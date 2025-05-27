package models

import (
	"math"

	"github.com/mlflow/mlflow-go-backend/pkg/entities"
)

type LoggedModelMetric struct {
	ModelID           string  `gorm:"column:model_id"`
	MetricName        string  `gorm:"column:metric_name"`
	MetricTimestampMs int64   `gorm:"column:metric_timestamp_ms"`
	MetricStep        int64   `gorm:"column:metric_step"`
	MetricValue       float64 `gorm:"column:metric_value"`
	ExperimentID      int     `gorm:"column:experiment_id"`
	RunID             string  `gorm:"column:run_id"`
	DatasetUUID       string  `gorm:"column:dataset_uuid"`
	DatasetName       string  `gorm:"column:dataset_name"`
	DatasetDigest     string  `gorm:"column:dataset_digest"`
}

func NewLoggedMetricFromEntity(runID string, metric *entities.Metric) *LoggedModelMetric {
	model := LoggedModelMetric{
		RunID:             runID,
		ModelID:           metric.ModelID,
		MetricName:        metric.Key,
		MetricTimestampMs: metric.Timestamp,
		DatasetName:       metric.DatasetName,
		DatasetDigest:     metric.DatasetDigest,
	}

	if metric.Step != 0 {
		model.MetricStep = metric.Step
	}

	switch {
	case math.IsNaN(metric.Value):
		model.MetricValue = 0
	case math.IsInf(metric.Value, 0):
		// NB: SQL cannot represent Infs => We replace +/- Inf with max/min 64b float value
		if metric.Value > 0 {
			model.MetricValue = math.MaxFloat64
		} else {
			model.MetricValue = -math.MaxFloat64
		}
	default:
		model.MetricValue = metric.Value
	}

	return &model
}
