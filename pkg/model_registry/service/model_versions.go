package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mlflow/mlflow-go/pkg/contract"
	"github.com/mlflow/mlflow-go/pkg/model_registry/store/sql/models"
	"github.com/mlflow/mlflow-go/pkg/protos"
)

func (m *ModelRegistryService) GetLatestVersions(
	ctx context.Context, input *protos.GetLatestVersions,
) (*protos.GetLatestVersions_Response, *contract.Error) {
	latestVersions, err := m.store.GetLatestVersions(ctx, input.GetName(), input.GetStages())
	if err != nil {
		return nil, err
	}

	return &protos.GetLatestVersions_Response{
		ModelVersions: latestVersions,
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
	// by some strange reason GetModelVersion.Version has a string type so we can't apply our validation,
	// that's why such a custom validation exists to satisfy Python tests.
	version := input.GetVersion()
	if _, err := strconv.Atoi(version); err != nil {
		return nil, contract.NewErrorWith(
			protos.ErrorCode_INVALID_PARAMETER_VALUE, "Model version must be an integer", err,
		)
	}

	modelVersion, err := m.store.GetModelVersion(ctx, input.GetName(), version, true)
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
