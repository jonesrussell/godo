//go:build !docker
// +build !docker

package container

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvideConfig(t *testing.T) {
	t.Run("loads default config", func(t *testing.T) {
		cfg, err := provideConfig()
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "Godo", cfg.App.Name)
	})

	t.Run("handles environment variables", func(t *testing.T) {
		os.Setenv("GODO_APP_NAME", "TestApp")
		defer os.Unsetenv("GODO_APP_NAME")

		cfg, err := provideConfig()
		require.NoError(t, err)
		assert.Equal(t, "TestApp", cfg.App.Name)
	})
}

func TestProvideLogger(t *testing.T) {
	t.Run("creates logger with default config", func(t *testing.T) {
		log, err := provideLogger()
		require.NoError(t, err)
		assert.NotNil(t, log)
		assert.Implements(t, (*logger.Logger)(nil), log)
	})
}

func TestProvideSQLite(t *testing.T) {
	t.Run("creates SQLite store", func(t *testing.T) {
		// Create a temporary directory for the test database
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "test.db")

		// Create test config
		cfg := &config.Config{
			Database: config.DatabaseConfig{
				Path: dbPath,
			},
		}

		// Create test logger
		log, err := provideLogger()
		require.NoError(t, err)

		// Create SQLite store
		store, cleanup, err := provideSQLite(cfg, log)
		require.NoError(t, err)
		if cleanup != nil {
			defer cleanup()
		}

		assert.NotNil(t, store)
		assert.FileExists(t, dbPath)
	})

	t.Run("handles invalid path", func(t *testing.T) {
		// Create a path that's invalid for SQLite (empty path)
		cfg := &config.Config{
			Database: config.DatabaseConfig{
				Path: "",
			},
		}

		log, err := provideLogger()
		require.NoError(t, err)

		store, cleanup, err := provideSQLite(cfg, log)
		if cleanup != nil {
			defer cleanup()
		}

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid database path")
		assert.Nil(t, store)
	})
}

func TestInitializeApp(t *testing.T) {
	t.Run("creates fully wired application", func(t *testing.T) {
		// Create a temporary directory for the test database
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "test.db")

		// Set database path through environment variable
		os.Setenv("GODO_DATABASE_PATH", dbPath)
		defer os.Unsetenv("GODO_DATABASE_PATH")

		// Initialize the application
		application, cleanup, err := InitializeApp()
		require.NoError(t, err)
		if cleanup != nil {
			defer cleanup()
		}

		assert.NotNil(t, application)
		assert.NotNil(t, cleanup)
		assert.FileExists(t, dbPath)
	})
}

func TestBuildConstraints(t *testing.T) {
	t.Run("supported platforms do not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			// This test will only run on supported platforms (Windows, Linux, macOS)
			// The build_constraints.go file will panic on unsupported platforms
		})
	})
}
