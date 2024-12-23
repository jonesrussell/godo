package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestConfig(t *testing.T) {
	// Create a new observer core
	observedZapCore, logs := observer.New(zap.DebugLevel)

	// Create a test logger that writes to the observer
	zapLogger := zap.New(observedZapCore)
	defer zapLogger.Sync()

	// Create a new Provider with our test logger
	provider := config.NewProvider(
		[]string{"testdata"},
		"config",
		"yaml",
	)

	// Set test mode to prevent path resolution
	os.Setenv("GODO_TEST_MODE", "true")
	defer os.Unsetenv("GODO_TEST_MODE")

	t.Run("Load default config", func(t *testing.T) {
		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "Godo", cfg.App.Name)
		assert.Equal(t, "0.1.0", cfg.App.Version)

		// Verify logs
		logEntries := logs.All()
		assert.True(t, len(logEntries) > 0, "Expected some log entries")

		// Check for specific log messages
		hasStartMessage := false
		for _, entry := range logEntries {
			if entry.Message == "starting config load" {
				hasStartMessage = true
				break
			}
		}
		assert.True(t, hasStartMessage, "Expected 'starting config load' message")
	})

	t.Run("Environment variables override config", func(t *testing.T) {
		os.Setenv("GODO_APP_NAME", "TestApp")
		os.Setenv("GODO_DATABASE_PATH", "test.db")
		defer func() {
			os.Unsetenv("GODO_APP_NAME")
			os.Unsetenv("GODO_DATABASE_PATH")
		}()

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "TestApp", cfg.App.Name)
		assert.Equal(t, "test.db", cfg.Database.Path)
	})

	t.Run("Invalid config validation", func(t *testing.T) {
		os.Setenv("GODO_APP_NAME", "")
		defer os.Unsetenv("GODO_APP_NAME")

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
		)

		_, err := provider.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app name is required")
	})

	t.Run("Invalid log level", func(t *testing.T) {
		os.Setenv("GODO_LOGGER_LEVEL", "invalid")
		defer os.Unsetenv("GODO_LOGGER_LEVEL")

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
		)

		_, err := provider.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})

	t.Run("Path resolution in production mode", func(t *testing.T) {
		os.Unsetenv("GODO_TEST_MODE")
		defer os.Setenv("GODO_TEST_MODE", "true")

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.True(t, filepath.IsAbs(cfg.Database.Path))
		assert.Contains(t, cfg.Database.Path, "godo.db")
	})
}

func TestConfigFileErrors(t *testing.T) {
	t.Run("Invalid YAML syntax", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "invalid.yaml")
		err := os.WriteFile(tmpFile, []byte("invalid: yaml: content:"), 0o600)
		require.NoError(t, err)

		provider := config.NewProvider(
			[]string{tmpDir},
			"invalid",
			"yaml",
		)

		_, err = provider.Load()
		assert.Error(t, err)
	})
}

func TestNewDefaultConfig(t *testing.T) {
	cfg := config.NewDefaultConfig()

	assert.Equal(t, "Godo", cfg.App.Name)
	assert.Equal(t, "0.1.0", cfg.App.Version)
	assert.Equal(t, "io.github.jonesrussell.godo", cfg.App.ID)
	assert.Equal(t, "info", cfg.Logger.Level)
	assert.True(t, cfg.Logger.Console)
	assert.Equal(t, "godo.db", cfg.Database.Path)
}
