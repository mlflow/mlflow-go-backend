package models

import "github.com/mlflow/mlflow-go-backend/pkg/entities"

type Output struct {
	ID              string `gorm:"column:input_uuid;not null"`
	Step            int64  `gorm:"column:step"`
	SourceType      string `gorm:"column:source_type;primaryKey"`
	SourceID        string `gorm:"column:source_id;primaryKey"`
	DestinationType string `gorm:"column:destination_type;primaryKey"`
	DestinationID   string `gorm:"column:destination_id;primaryKey"`
}

func (o *Output) TableName() string {
	return "inputs"
}

func (o *Output) ToEntity() *entities.ModelOutput {
	return &entities.ModelOutput{
		Step:    o.Step,
		ModelID: o.DestinationID,
	}
}
