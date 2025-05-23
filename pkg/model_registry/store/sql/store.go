package sql

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/mlflow/mlflow-go-backend/pkg/config"
	"github.com/mlflow/mlflow-go-backend/pkg/sql"
)

type ModelRegistrySQLStore struct {
	config *config.Config
	db     *gorm.DB
}

func NewModelRegistrySQLStore(ctx context.Context, config *config.Config) (*ModelRegistrySQLStore, error) {
	database, err := sql.NewDatabase(ctx, config.ModelRegistryStoreURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %q: %w", config.ModelRegistryStoreURI, err)
	}

	return &ModelRegistrySQLStore{
		config: config,
		db:     database,
	}, nil
}

func (m *ModelRegistrySQLStore) Destroy() error {
	if err := sql.CloseDatabase(m.db); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
