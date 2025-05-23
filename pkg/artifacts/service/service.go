package service

import (
	"context"

	"github.com/mlflow/mlflow-go-backend/pkg/config"
)

type ArtifactsService struct {
	config *config.Config
}

func NewArtifactsService(_ context.Context, config *config.Config) (*ArtifactsService, error) {
	return &ArtifactsService{
		config: config,
	}, nil
}

func (as ArtifactsService) Destroy() error {
	return nil
}
