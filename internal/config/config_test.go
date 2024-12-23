package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestConfig(t *testing.T) {
	// Create a zap observer for testing
	observedZapCore, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(observedZapCore)
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			t.Logf("failed to sync logger: %v", err)
		}
	}()

	// Use your standard logger implementation
	log := logger.NewZapLogger(zapLogger)

	provider := config.NewProvider(
		[]string{"testdata"},
		"config",
		"yaml",
		config.WithLogger(log),
	)

	// Set test mode to prevent path resolution
	os.Setenv(config.EnvTestMode, "true")
	defer os.Unsetenv(config.EnvTestMode)

	t.Run("Load default config", func(t *testing.T) {
		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, config.DefaultAppName, cfg.App.Name)
		assert.Equal(t, config.DefaultAppVersion, cfg.App.Version)

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
		os.Setenv(config.EnvPrefix+"_APP_NAME", "TestApp")
		os.Setenv(config.EnvPrefix+"_DATABASE_PATH", "test.db")
		defer func() {
			os.Unsetenv(config.EnvPrefix + "_APP_NAME")
			os.Unsetenv(config.EnvPrefix + "_DATABASE_PATH")
		}()

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "TestApp", cfg.App.Name)
		assert.Equal(t, "test.db", cfg.Database.Path)
	})

	t.Run("Invalid_config_validation", func(t *testing.T) {
		// Setup invalid config
		cfg := &config.Config{
			App: config.AppConfig{
				Name: "", // Invalid: empty name
			},
			Logger: common.LogConfig{
				Level: "invalid", // Invalid log level
			},
		}

		// Test validation
		err := config.ValidateConfig(cfg)
		assert.Error(t, err, "should fail validation with empty app name and invalid log level")

		// Optional: Check specific validation errors
		if configErr, ok := err.(*config.ConfigError); ok {
			assert.Contains(t, configErr.Error(), "app name is required")
			assert.Contains(t, configErr.Error(), "invalid log level")
		}
	})

	t.Run("Invalid log level", func(t *testing.T) {
		os.Setenv(config.EnvPrefix+"_LOGGER_LEVEL", "invalid")
		defer os.Unsetenv(config.EnvPrefix + "_LOGGER_LEVEL")

		_, err := provider.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})

	t.Run("Path resolution in production mode", func(t *testing.T) {
		os.Unsetenv(config.EnvTestMode)
		defer os.Setenv(config.EnvTestMode, "true")

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.True(t, filepath.IsAbs(cfg.Database.Path))
		assert.Contains(t, cfg.Database.Path, config.DefaultDBPath)
	})
}

func TestConfigFileErrors(t *testing.T) {
	t.Run("Invalid YAML syntax", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create invalid YAML with the correct filename (invalid.yaml)
		invalidYAML := []byte(`
app:
  name: [test
    version: 1.0.0]
  id: {broken:
`)
		err := os.WriteFile(filepath.Join(tmpDir, "invalid.yaml"), invalidYAML, 0o600)
		require.NoError(t, err)

		provider := config.NewProvider(
			[]string{tmpDir},
			"invalid", // This matches the filename we created
			"yaml",
			config.WithLogger(logger.NewTestLogger(t)), // Add logging for better debugging
		)

		cfg, err := provider.Load()
		if err == nil {
			t.Logf("Config loaded when it shouldn't: %+v", cfg)
			t.Fatal("Expected error for invalid YAML, got nil")
		}
		assert.Contains(t, err.Error(), "yaml", "Error should mention YAML parsing")
	})
}

func TestNewDefaultConfig(t *testing.T) {
	cfg := config.NewDefaultConfig()

	assert.Equal(t, config.DefaultAppName, cfg.App.Name)
	assert.Equal(t, config.DefaultAppVersion, cfg.App.Version)
	assert.Equal(t, "io.github.jonesrussell/godo", cfg.App.ID)
	assert.Equal(t, config.DefaultLogLevel, cfg.Logger.Level)
	assert.True(t, cfg.Logger.Console)
	assert.Equal(t, config.DefaultDBPath, cfg.Database.Path)
}
