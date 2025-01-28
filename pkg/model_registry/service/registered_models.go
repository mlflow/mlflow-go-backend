package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mlflow/mlflow-go-backend/pkg/contract"
	"github.com/mlflow/mlflow-go-backend/pkg/entities"
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
)

var (
	RegisteredModelAliasRegex        = regexp.MustCompile(`^[\w\-]*$`)
	RegisteredModelAliasVersionRegex = regexp.MustCompile(`^[vV]\d+$`)
)

func (m *ModelRegistryService) UpdateRegisteredModel(
	ctx context.Context, input *protos.UpdateRegisteredModel,
) (*protos.UpdateRegisteredModel_Response, *contract.Error) {
	registeredModel, err := m.store.UpdateRegisteredModel(ctx, input.GetName(), input.GetDescription())
	if err != nil {
		return nil, err
	}

	return &protos.UpdateRegisteredModel_Response{
		RegisteredModel: registeredModel.ToProto(),
	}, nil
}

func (m *ModelRegistryService) RenameRegisteredModel(
	ctx context.Context, input *protos.RenameRegisteredModel,
) (*protos.RenameRegisteredModel_Response, *contract.Error) {
	registeredModel, err := m.store.RenameRegisteredModel(ctx, input.GetName(), input.GetNewName())
	if err != nil {
		return nil, err
	}

	return &protos.RenameRegisteredModel_Response{
		RegisteredModel: registeredModel.ToProto(),
	}, nil
}

func (m *ModelRegistryService) DeleteRegisteredModel(
	ctx context.Context, input *protos.DeleteRegisteredModel,
) (*protos.DeleteRegisteredModel_Response, *contract.Error) {
	if err := m.store.DeleteRegisteredModel(ctx, input.GetName()); err != nil {
		return nil, err
	}

	return &protos.DeleteRegisteredModel_Response{}, nil
}

func (m *ModelRegistryService) GetRegisteredModel(
	ctx context.Context, input *protos.GetRegisteredModel,
) (*protos.GetRegisteredModel_Response, *contract.Error) {
	registeredModel, err := m.store.GetRegisteredModel(ctx, input.GetName())
	if err != nil {
		return nil, err
	}

	return &protos.GetRegisteredModel_Response{
		RegisteredModel: registeredModel.ToProto(),
	}, nil
}

func (m *ModelRegistryService) SetRegisteredModelTag(
	ctx context.Context, input *protos.SetRegisteredModelTag,
) (*protos.SetRegisteredModelTag_Response, *contract.Error) {
	if err := m.store.SetRegisteredModelTag(ctx, input.GetName(), input.GetKey(), input.GetValue()); err != nil {
		return nil, err
	}

	return &protos.SetRegisteredModelTag_Response{}, nil
}

func (m *ModelRegistryService) CreateRegisteredModel(
	ctx context.Context, input *protos.CreateRegisteredModel,
) (*protos.CreateRegisteredModel_Response, *contract.Error) {
	tags := make([]*entities.RegisteredModelTag, 0, len(input.GetTags()))
	for _, tag := range input.GetTags() {
		tags = append(tags, entities.NewRegisteredModelTagFromProto(tag))
	}

	registeredModel, err := m.store.CreateRegisteredModel(ctx, input.GetName(), input.GetDescription(), tags)
	if err != nil {
		return nil, err
	}

	return &protos.CreateRegisteredModel_Response{
		RegisteredModel: registeredModel.ToProto(),
	}, nil
}

func (m *ModelRegistryService) DeleteRegisteredModelTag(
	ctx context.Context, input *protos.DeleteRegisteredModelTag,
) (*protos.DeleteRegisteredModelTag_Response, *contract.Error) {
	if err := m.store.DeleteRegisteredModelTag(ctx, input.GetName(), input.GetKey()); err != nil {
		return nil, err
	}

	return &protos.DeleteRegisteredModelTag_Response{}, nil
}

func (m *ModelRegistryService) SetRegisteredModelAlias(
	ctx context.Context, input *protos.SetRegisteredModelAlias,
) (*protos.SetRegisteredModelAlias_Response, *contract.Error) {
	alias := input.GetAlias()
	if !RegisteredModelAliasRegex.MatchString(alias) {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			fmt.Sprintf(
				"Invalid alias name: %s. Names may only contain alphanumerics, underscores (_), and dashes (-).",
				alias,
			),
		)
	}

	if alias == "latest" {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			"'latest' alias name (case insensitive) is reserved.",
		)
	}

	if RegisteredModelAliasVersionRegex.MatchString(alias) {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			fmt.Sprintf("Version alias name '%s' is reserved.", alias),
		)
	}

	if err := m.store.SetRegisteredModelAlias(
		ctx,
		input.GetName(),
		alias,
		input.GetVersion(),
	); err != nil {
		return nil, err
	}

	return &protos.SetRegisteredModelAlias_Response{}, nil
}

func (m *ModelRegistryService) DeleteRegisteredModelAlias(
	ctx context.Context, input *protos.DeleteRegisteredModelAlias,
) (*protos.DeleteRegisteredModelAlias_Response, *contract.Error) {
	alias := input.GetAlias()
	if !RegisteredModelAliasRegex.MatchString(alias) {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			fmt.Sprintf(
				"Invalid alias name: %s. Names may only contain alphanumerics, underscores (_), and dashes (-).",
				alias,
			),
		)
	}

	if alias == "latest" {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			"'latest' alias name (case insensitive) is reserved.",
		)
	}

	if RegisteredModelAliasVersionRegex.MatchString(alias) {
		return nil, contract.NewError(
			protos.ErrorCode_INVALID_PARAMETER_VALUE,
			fmt.Sprintf("Version alias name '%s' is reserved.", alias),
		)
	}

	if err := m.store.DeleteRegisteredModelAlias(
		ctx,
		input.GetName(),
		alias,
	); err != nil {
		return nil, err
	}

	return &protos.DeleteRegisteredModelAlias_Response{}, nil
}
