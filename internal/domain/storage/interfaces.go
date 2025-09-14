// Package storage defines storage interfaces for the domain layer
package storage

import (
	"context"

	"github.com/jonesrussell/godo/internal/domain/model"
)

//go:generate mockgen -destination=../../test/mocks/mock_unified_storage.go -package=mocks github.com/jonesrussell/godo/internal/domain/storage UnifiedNoteStorage

// UnifiedNoteStorage defines the unified interface for note storage operations
// This interface supports both local SQLite and remote API storage backends
type UnifiedNoteStorage interface {
	// Basic CRUD operations
	CreateNote(ctx context.Context, content string) (*model.Note, error)
	GetNote(ctx context.Context, id string) (*model.Note, error)
	GetAllNotes(ctx context.Context) ([]*model.Note, error)
	UpdateNote(ctx context.Context, id string, content string, done bool) (*model.Note, error)
	DeleteNote(ctx context.Context, id string) error

	// Convenience operations
	ToggleDone(ctx context.Context, id string) (*model.Note, error)
	MarkDone(ctx context.Context, id string) (*model.Note, error)
	MarkUndone(ctx context.Context, id string) (*model.Note, error)

	// Resource management
	Close() error
}

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeSQLite StorageType = "sqlite"
	StorageTypeAPI    StorageType = "api"
)

// StorageConfig holds configuration for storage backends
type StorageConfig struct {
	Type   StorageType  `mapstructure:"type" json:"type"`
	SQLite SQLiteConfig `mapstructure:"sqlite" json:"sqlite"`
	API    APIConfig    `mapstructure:"api" json:"api"`
}

// SQLiteConfig holds SQLite-specific configuration
type SQLiteConfig struct {
	FilePath string `mapstructure:"file_path" json:"file_path"`
}

// APIConfig holds API-specific configuration
type APIConfig struct {
	BaseURL            string `mapstructure:"base_url" json:"base_url"`
	Timeout            int    `mapstructure:"timeout_seconds" json:"timeout_seconds"`
	RetryCount         int    `mapstructure:"retry_count" json:"retry_count"`
	RetryDelay         int    `mapstructure:"retry_delay_ms" json:"retry_delay_ms"`
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify" json:"insecure_skip_verify"`
}

// DefaultStorageConfig returns a default storage configuration
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		Type: StorageTypeSQLite,
		SQLite: SQLiteConfig{
			FilePath: "godo.db",
		},
		API: APIConfig{
			BaseURL:            "http://localhost:8000/api",
			Timeout:            30,
			RetryCount:         3,
			RetryDelay:         1000,
			InsecureSkipVerify: false,
		},
	}
}
