package entities

import (
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
	"github.com/mlflow/mlflow-go-backend/pkg/utils"
)

type ModelOutput struct {
	Step    int64
	ModelID string
}

func (mo ModelOutput) ToProto() *protos.ModelOutput {
	return &protos.ModelOutput{
		Step:    utils.PtrTo(mo.Step),
		ModelId: utils.PtrTo(mo.ModelID),
	}
}
