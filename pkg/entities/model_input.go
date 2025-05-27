package entities

import (
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
	"github.com/mlflow/mlflow-go-backend/pkg/utils"
)

type ModelInput struct {
	ModelID string
}

func (mi ModelInput) ToProto() *protos.ModelInput {
	return &protos.ModelInput{
		ModelId: utils.PtrTo(mi.ModelID),
	}
}

func NewModelInputFromProto(proto *protos.ModelInput) *ModelInput {
	return &ModelInput{
		ModelID: *proto.ModelId,
	}
}
