package entities

import (
	"github.com/mlflow/mlflow-go-backend/pkg/protos"
	"github.com/mlflow/mlflow-go-backend/pkg/utils"
)

type ModelVersionTag struct {
	Key   string
	Value string
}

func (mvt ModelVersionTag) ToProto() *protos.ModelVersionTag {
	return &protos.ModelVersionTag{
		Key:   utils.PtrTo(mvt.Key),
		Value: utils.PtrTo(mvt.Value),
	}
}

func NewModelVersionTag(proto *protos.ModelVersionTag) *ModelVersionTag {
	return &ModelVersionTag{
		Key:   proto.GetKey(),
		Value: proto.GetValue(),
	}
}
