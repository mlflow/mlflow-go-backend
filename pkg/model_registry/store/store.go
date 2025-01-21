package store

import (
	"context"

	"github.com/mlflow/mlflow-go-backend/pkg/contract"
	"github.com/mlflow/mlflow-go-backend/pkg/entities"
	"github.com/mlflow/mlflow-go-backend/pkg/model_registry/store/sql/models"
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
)

type ModelRegistryStore interface {
	contract.Destroyer
	ModelVersionStore
	RegisteredModelStore
}

type ModelVersionStore interface {
	GetLatestVersions(ctx context.Context, name string, stages []string) ([]*protos.ModelVersion, *contract.Error)
	GetModelVersion(ctx context.Context, name, version string, eager bool) (*entities.ModelVersion, *contract.Error)
	DeleteModelVersion(ctx context.Context, name, version string) *contract.Error
	UpdateModelVersion(ctx context.Context, name, version, description string) (*entities.ModelVersion, *contract.Error)
	TransitionModelVersionStage(
		ctx context.Context, name, version string, stage models.ModelVersionStage, archiveExistingVersions bool,
	) (*entities.ModelVersion, *contract.Error)
}

type RegisteredModelStore interface {
	CreateRegisteredModel(
		ctx context.Context, name, description string, tags []*entities.RegisteredModelTag,
	) (*entities.RegisteredModel, *contract.Error)
	GetRegisteredModel(ctx context.Context, name string) (*entities.RegisteredModel, *contract.Error)
	UpdateRegisteredModel(ctx context.Context, name, description string) (*entities.RegisteredModel, *contract.Error)
	RenameRegisteredModel(ctx context.Context, name, newName string) (*entities.RegisteredModel, *contract.Error)
	DeleteRegisteredModel(ctx context.Context, name string) *contract.Error
	SetRegisteredModelTag(ctx context.Context, name, key, value string) *contract.Error
	DeleteRegisteredModelTag(ctx context.Context, name, key string) *contract.Error
	SetRegisteredModelAlias(ctx context.Context, name, alias, version string) *contract.Error
	DeleteRegisteredModelAlias(ctx context.Context, name, alias string) *contract.Error
}
