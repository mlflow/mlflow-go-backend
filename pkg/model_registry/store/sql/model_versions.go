package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/mlflow/mlflow-go-backend/pkg/contract"
	"github.com/mlflow/mlflow-go-backend/pkg/entities"
	"github.com/mlflow/mlflow-go-backend/pkg/model_registry/store/sql/models"
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
)

func (m *ModelRegistrySQLStore) GetLatestVersions(
	ctx context.Context, name string, stages []string,
) ([]*protos.ModelVersion, *contract.Error) {
	if _, err := m.GetRegisteredModel(ctx, name); err != nil {
		return nil, err
	}

	var modelVersions []*models.ModelVersion

	subQuery := m.db.
		WithContext(ctx).
		Model(&models.ModelVersion{}).
		Select("name, MAX(version) AS max_version").
		Where("name = ?", name).
		Where("current_stage <> ?", models.StageDeletedInternal).
		Group("name, current_stage")

	if len(stages) > 0 {
		for idx, stage := range stages {
			stages[idx] = strings.ToLower(stage)
			if canonicalStage, ok := models.CanonicalMapping[stages[idx]]; ok {
				stages[idx] = canonicalStage.String()

				continue
			}

			return nil, contract.NewError(
				protos.ErrorCode_BAD_REQUEST,
				fmt.Sprintf(
					"Invalid Model Version stage: %s. Value must be one of %s.",
					stage,
					models.AllModelVersionStages(),
				),
			)
		}

		subQuery = subQuery.Where("current_stage IN (?)", stages)
	}

	err := m.db.
		WithContext(ctx).
		Model(&models.ModelVersion{}).
		Joins("JOIN (?) AS sub ON model_versions.name = sub.name AND model_versions.version = sub.max_version", subQuery).
		Find(&modelVersions).Error
	if err != nil {
		return nil, contract.NewErrorWith(
			protos.ErrorCode_INTERNAL_ERROR,
			fmt.Sprintf("failed to query latest model version for %q", name),
			err,
		)
	}

	results := make([]*protos.ModelVersion, 0, len(modelVersions))
	for _, modelVersion := range modelVersions {
		results = append(results, modelVersion.ToProto())
	}

	return results, nil
}

func (m *ModelRegistrySQLStore) GetModelVersion(
	ctx context.Context, name, version string, eager bool,
) (*entities.ModelVersion, *contract.Error) {
	query := m.db.WithContext(
		ctx,
	).Where(
		"name = ?", name,
	).Where(
		"version = ?", version,
	).Where(
		"current_stage != ?", models.StageDeletedInternal,
	)

	// preload Tags only by demand.
	if eager {
		query = query.Preload("Tags")
	}

	var modelVersion models.ModelVersion
	if err := query.First(
		&modelVersion,
	).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contract.NewError(
				protos.ErrorCode_RESOURCE_DOES_NOT_EXIST,
				fmt.Sprintf("Model Version (name=%s, version=%s) not found", name, version),
			)
		}

		return nil, contract.NewErrorWith(
			protos.ErrorCode_INTERNAL_ERROR,
			fmt.Sprintf("failed to get Model Version by name %s and version %s", name, version),
			err,
		)
	}

	var registeredModelAliases []models.RegisteredModelAlias
	if err := m.db.WithContext(ctx).Where(
		"name = ?", modelVersion.Name,
	).Where(
		"version = ?", modelVersion.Version,
	).Find(
		&registeredModelAliases,
	).Error; err != nil {
		return nil, contract.NewErrorWith(
			protos.ErrorCode_INTERNAL_ERROR,
			fmt.Sprintf("failed to get Registered Model Aliases by name %s and version %s", name, version),
			err,
		)
	}

	modelVersion.Aliases = append(modelVersion.Aliases, registeredModelAliases...)

	return modelVersion.ToEntity(), nil
}

func (m *ModelRegistrySQLStore) DeleteModelVersion(ctx context.Context, name, version string) *contract.Error {
	registeredModel, err := m.GetRegisteredModel(ctx, name)
	if err != nil {
		return err
	}

	modelVersion, err := m.GetModelVersion(ctx, name, version, false)
	if err != nil {
		return err
	}

	if err := m.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		if err := transaction.Model(
			&models.RegisteredModel{},
		).Where(
			"name = ?", registeredModel.Name,
		).Updates(&models.RegisteredModel{
			LastUpdatedTime: time.Now().UnixMilli(),
		}).Error; err != nil {
			return err
		}

		if err := transaction.Where(
			"name = ?", registeredModel.Name,
		).Where(
			"version = ?", version,
		).Delete(
			&models.RegisteredModelAlias{},
		).Error; err != nil {
			return err
		}

		if err := transaction.Model(
			&models.ModelVersion{},
		).Where(
			"name = ?", modelVersion.Name,
		).Where(
			"version = ?", modelVersion.Version,
		).Updates(&models.ModelVersion{
			RunID:           "REDACTED-RUN-ID",
			UserID:          sql.NullString{Valid: true},
			Source:          "REDACTED-SOURCE-PATH",
			RunLink:         "REDACTED-RUN-LINK",
			CurrentStage:    models.StageDeletedInternal,
			Description:     sql.NullString{Valid: true},
			StatusMessage:   sql.NullString{Valid: true},
			LastUpdatedTime: time.Now().UnixMilli(),
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return contract.NewErrorWith(protos.ErrorCode_INTERNAL_ERROR, "error deleting model version", err)
	}

	return nil
}

func (m *ModelRegistrySQLStore) UpdateModelVersion(
	ctx context.Context, name, version, description string,
) (*entities.ModelVersion, *contract.Error) {
	modelVersion, err := m.GetModelVersion(ctx, name, version, false)
	if err != nil {
		return nil, err
	}

	if err := m.db.WithContext(ctx).Model(
		&models.ModelVersion{},
	).Where(
		"name = ?", modelVersion.Name,
	).Where(
		"version = ?", modelVersion.Version,
	).Updates(&models.ModelVersion{
		Description:     sql.NullString{String: description, Valid: description != ""},
		LastUpdatedTime: time.Now().UnixMilli(),
	}).Error; err != nil {
		return nil, contract.NewErrorWith(protos.ErrorCode_INTERNAL_ERROR, "error updating model version", err)
	}

	return modelVersion, nil
}

//nolint:funlen,cyclop
func (m *ModelRegistrySQLStore) TransitionModelVersionStage(
	ctx context.Context, name, version string, stage models.ModelVersionStage, archiveExistingVersions bool,
) (*entities.ModelVersion, *contract.Error) {
	isActiveStage := false
	if _, ok := models.DefaultStagesForGetLatestVersions[strings.ToLower(stage.String())]; ok {
		isActiveStage = true
	}

	if archiveExistingVersions && !isActiveStage {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			fmt.Sprintf(
				`Model version transition cannot archive existing model versions because '%s' is not an Active stage. 
Valid stages are %s`,
				stage, models.AllModelVersionStages(),
			),
		)
	}

	modelVersion, err := m.GetModelVersion(ctx, name, version, false)
	if err != nil {
		return nil, err
	}

	registeredModel, err := m.GetRegisteredModel(ctx, name)
	if err != nil {
		return nil, err
	}

	if err := m.db.Transaction(func(transaction *gorm.DB) error {
		lastUpdatedTime := time.Now().UnixMilli()
		if err := transaction.Model(
			&models.RegisteredModel{},
		).Where(
			"name = ?", registeredModel.Name,
		).Updates(&models.RegisteredModel{
			LastUpdatedTime: lastUpdatedTime,
		}).Error; err != nil {
			return err
		}

		if err := transaction.Model(
			&models.ModelVersion{},
		).Where(
			"name = ?", modelVersion.Name,
		).Where(
			"version = ?", modelVersion.Version,
		).Updates(&models.ModelVersion{
			CurrentStage:    stage,
			LastUpdatedTime: lastUpdatedTime,
		}).Error; err != nil {
			return err
		}

		if archiveExistingVersions {
			if err := transaction.Where(
				"name = ?", name,
			).Where(
				"version != ?", version,
			).Where(
				"current_stage = ?", stage,
			).Updates(&models.ModelVersion{
				CurrentStage:    models.ModelVersionStageArchived,
				LastUpdatedTime: lastUpdatedTime,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, contract.NewErrorWith(
			protos.ErrorCode_INTERNAL_ERROR, "error transitioning model version stage", err,
		)
	}

	return modelVersion, nil
}


func (m *ModelRegistrySQLStore) DeleteModelVersionTag(
	ctx context.Context, name, version, key string,
) *contract.Error {
	modelVersion, err := m.GetModelVersion(ctx, name, version, false)
	if err != nil {
		return err
	}

	if err := m.db.WithContext(ctx).Where(
		"name = ?", modelVersion.Name,
	).Where(
		"key = ?", key,
	).Where(
		"version = ?", version,
	).Delete(
		&models.ModelVersionTag{},
	).Error; err != nil {
		return contract.NewErrorWith(
			protos.ErrorCode_INTERNAL_ERROR, "error deleting model version tag", err,
		)
	}

	return nil
}

func (m *ModelRegistrySQLStore) GetModelVersionByAlias(
	ctx context.Context, name, alias string,
) (*entities.ModelVersion, *contract.Error) {
	var registeredModelAlias models.RegisteredModelAlias
	if err := m.db.WithContext(
		ctx,
	).Where(
		"name = ?", name,
	).Where(
		"alias = ?", alias,
	).First(
		&registeredModelAlias,
	).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contract.NewError(
				protos.ErrorCode_INVALID_PARAMETER_VALUE,
				fmt.Sprintf("Registered model alias %s not found.", alias),
			)
		}

		return nil, contract.NewErrorWith(
			protos.ErrorCode_INTERNAL_ERROR, "error getting registered model alias", err,
		)
	}

	modelVersion, err := m.GetModelVersion(
		ctx,
		name,
		strconv.Itoa(int(registeredModelAlias.Version)),
		false,
	)
	if err != nil {
		return nil, err
	}

	return modelVersion, nil
}
