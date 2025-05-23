package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/mlflow/mlflow-go-backend/pkg/contract"
	"github.com/mlflow/mlflow-go-backend/pkg/model_registry/store/sql/models"
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
	"github.com/mlflow/mlflow-go-backend/pkg/utils"
)

func (m *ModelRegistryService) GetLatestVersions(
	ctx context.Context, input *protos.GetLatestVersions,
) (*protos.GetLatestVersions_Response, *contract.Error) {
	latestVersions, err := m.store.GetLatestVersions(ctx, input.GetName(), input.GetStages())
	if err != nil {
		return nil, err
	}

	modelVersions := make([]*protos.ModelVersion, 0, len(latestVersions))
	for _, lastVersion := range latestVersions {
		modelVersions = append(modelVersions, lastVersion.ToProto())
	}

	return &protos.GetLatestVersions_Response{
		ModelVersions: modelVersions,
	}, nil
}

func (m *ModelRegistryService) DeleteModelVersion(
	ctx context.Context, input *protos.DeleteModelVersion,
) (*protos.DeleteModelVersion_Response, *contract.Error) {
	if err := m.store.DeleteModelVersion(ctx, input.GetName(), input.GetVersion()); err != nil {
		return nil, err
	}

	return &protos.DeleteModelVersion_Response{}, nil
}

func (m *ModelRegistryService) GetModelVersion(
	ctx context.Context, input *protos.GetModelVersion,
) (*protos.GetModelVersion_Response, *contract.Error) {
	modelVersion, err := m.store.GetModelVersion(ctx, input.GetName(), input.GetVersion(), true)
	if err != nil {
		return nil, err
	}

	return &protos.GetModelVersion_Response{
		ModelVersion: modelVersion.ToProto(),
	}, nil
}

func (m *ModelRegistryService) UpdateModelVersion(
	ctx context.Context, input *protos.UpdateModelVersion,
) (*protos.UpdateModelVersion_Response, *contract.Error) {
	modelVersion, err := m.store.UpdateModelVersion(ctx, input.GetName(), input.GetVersion(), input.GetDescription())
	if err != nil {
		return nil, err
	}

	return &protos.UpdateModelVersion_Response{
		ModelVersion: modelVersion.ToProto(),
	}, nil
}

func (m *ModelRegistryService) TransitionModelVersionStage(
	ctx context.Context, input *protos.TransitionModelVersionStage,
) (*protos.TransitionModelVersionStage_Response, *contract.Error) {
	stage, ok := models.CanonicalMapping[strings.ToLower(input.GetStage())]
	if !ok {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			fmt.Sprintf(
				"Invalid Model Version stage: unknown. Value must be one of %s, %s, %s, %s.",
				models.ModelVersionStageNone,
				models.ModelVersionStageStaging,
				models.ModelVersionStageProduction,
				models.ModelVersionStageArchived,
			),
		)
	}

	modelVersion, err := m.store.TransitionModelVersionStage(
		ctx,
		input.GetName(),
		input.GetVersion(),
		stage,
		input.GetArchiveExistingVersions(),
	)
	if err != nil {
		return nil, err
	}

	return &protos.TransitionModelVersionStage_Response{
		ModelVersion: modelVersion.ToProto(),
	}, nil
}

func (m *ModelRegistryService) DeleteModelVersionTag(
	ctx context.Context, input *protos.DeleteModelVersionTag,
) (*protos.DeleteModelVersionTag_Response, *contract.Error) {
	if err := m.store.DeleteModelVersionTag(
		ctx, input.GetName(), input.GetVersion(), input.GetKey(),
	); err != nil {
		return nil, err
	}

	return &protos.DeleteModelVersionTag_Response{}, nil
}

func (m *ModelRegistryService) GetModelVersionByAlias(
	ctx context.Context, input *protos.GetModelVersionByAlias,
) (*protos.GetModelVersionByAlias_Response, *contract.Error) {
	modelVersion, err := m.store.GetModelVersionByAlias(ctx, input.GetName(), input.GetAlias())
	if err != nil {
		return nil, err
	}

	return &protos.GetModelVersionByAlias_Response{
		ModelVersion: modelVersion.ToProto(),
	}, nil
}

func (m *ModelRegistryService) SetModelVersionTag(
	ctx context.Context, input *protos.SetModelVersionTag,
) (*protos.SetModelVersionTag_Response, *contract.Error) {
	if err := m.store.SetModelVersionTag(
		ctx,
		input.GetName(),
		input.GetVersion(),
		input.GetKey(),
		input.GetValue(),
	); err != nil {
		return nil, err
	}

	return &protos.SetModelVersionTag_Response{}, nil
}

//nolint:revive,stylecheck
func (m *ModelRegistryService) GetModelVersionDownloadUri(
	ctx context.Context, input *protos.GetModelVersionDownloadUri,
) (*protos.GetModelVersionDownloadUri_Response, *contract.Error) {
	artifactURI, err := m.store.GetModelVersionDownloadURI(ctx, input.GetName(), input.GetVersion())
	if err != nil {
		return nil, err
	}

	return &protos.GetModelVersionDownloadUri_Response{
		ArtifactUri: utils.PtrTo(artifactURI),
	}, nil
}
