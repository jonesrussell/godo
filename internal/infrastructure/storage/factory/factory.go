// Package factory provides factory functions for creating storage implementations
package factory

import (
	"fmt"

	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/api"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/sqlite"
)

// NewUnifiedStorage creates a new storage implementation based on configuration
func NewUnifiedStorage(config *domainstorage.StorageConfig, log logger.Logger) (domainstorage.UnifiedNoteStorage, error) {
	if config == nil {
		return nil, fmt.Errorf("storage configuration is required")
	}

	if log == nil {
		return nil, fmt.Errorf("logger is required")
	}

	switch config.Type {
	case domainstorage.StorageTypeSQLite:
		return NewSQLiteStorage(&config.SQLite, log)
	case domainstorage.StorageTypeAPI:
		return NewAPIStorage(&config.API, log)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}

// NewSQLiteStorage creates a new SQLite storage implementation
func NewSQLiteStorage(config *domainstorage.SQLiteConfig, log logger.Logger) (domainstorage.UnifiedNoteStorage, error) {
	if config == nil {
		return nil, fmt.Errorf("SQLite configuration is required")
	}

	if config.FilePath == "" {
		return nil, fmt.Errorf("SQLite file path is required")
	}

	log.Debug("Creating SQLite storage", "file_path", config.FilePath)

	store, err := sqlite.New(config.FilePath, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite store: %w", err)
	}

	adapter := sqlite.NewUnifiedAdapter(store)
	log.Info("SQLite storage created successfully", "file_path", config.FilePath)

	return adapter, nil
}

// NewAPIStorage creates a new API storage implementation
func NewAPIStorage(config *domainstorage.APIConfig, log logger.Logger) (domainstorage.UnifiedNoteStorage, error) {
	if config == nil {
		return nil, fmt.Errorf("API configuration is required")
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("API base URL is required")
	}

	log.Debug("Creating API storage", "base_url", config.BaseURL, "timeout", config.Timeout)

	store, err := api.New(*config, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create API store: %w", err)
	}

	log.Info("API storage created successfully", "base_url", config.BaseURL)

	return store, nil
}
