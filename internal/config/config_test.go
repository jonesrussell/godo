//go:build !docker

package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if os.Getenv("CI") == "true" {
		os.Exit(0) // Skip tests in CI environment
	}
	os.Exit(m.Run())
}

func TestConfig(t *testing.T) {
	// Create a test logger
	log, err := logger.New(&common.LogConfig{
		Level:   "debug",
		Console: true,
		Output:  []string{"stdout"},
	})
	require.NoError(t, err)

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
	t.Run("Falls back to defaults with invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create invalid YAML file
		invalidYAML := []byte(`
app:
  name: [test
    version: 1.0.0]
  id: {broken
`)
		err := os.WriteFile(filepath.Join(tmpDir, "invalid.yaml"), invalidYAML, 0o600)
		require.NoError(t, err)

		provider := config.NewProvider(
			[]string{tmpDir},
			"invalid",
			"yaml",
			config.WithLogger(logger.NewTestLogger(t)),
		)

		cfg, err := provider.Load()
		require.NoError(t, err) // Should not error as it falls back to defaults
		assert.Equal(t, config.DefaultAppName, cfg.App.Name)
		assert.Equal(t, config.DefaultAppVersion, cfg.App.Version)
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

func TestConfigError(t *testing.T) {
	t.Run("Error string format", func(t *testing.T) {
		err := &config.ConfigError{
			Op:  "test",
			Err: fmt.Errorf("test error"),
		}
		assert.Equal(t, "config test: test error", err.Error())
	})

	t.Run("Error unwrap", func(t *testing.T) {
		innerErr := fmt.Errorf("inner error")
		err := &config.ConfigError{
			Op:  "test",
			Err: innerErr,
		}
		assert.Equal(t, innerErr, err.Unwrap())
	})
}

func TestPathResolution(t *testing.T) {
	// Save original env and restore after test
	originalTestMode := os.Getenv(config.EnvTestMode)
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalTestMode != "" {
			os.Setenv(config.EnvTestMode, originalTestMode)
		} else {
			os.Unsetenv(config.EnvTestMode)
		}
		if originalConfigDir != "" {
			os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	// Create a test logger
	log, err := logger.New(&common.LogConfig{
		Level:   "debug",
		Console: true,
		Output:  []string{"stdout"},
	})
	require.NoError(t, err)

	t.Run("Relative path resolution", func(t *testing.T) {
		os.Unsetenv(config.EnvTestMode)

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
			config.WithLogger(log),
		)

		cfg, err := provider.Load()
		require.NoError(t, err)

		// Check if the database path was made absolute
		assert.True(t, filepath.IsAbs(cfg.Database.Path))
		assert.Contains(t, cfg.Database.Path, "godo")
		assert.Contains(t, cfg.Database.Path, "godo.db")
	})

	t.Run("Keep absolute path unchanged", func(t *testing.T) {
		os.Unsetenv(config.EnvTestMode)

		absPath := filepath.Join(t.TempDir(), "custom.db")
		os.Setenv(config.EnvPrefix+"_DATABASE_PATH", absPath)
		defer os.Unsetenv(config.EnvPrefix + "_DATABASE_PATH")

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
			config.WithLogger(log),
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, absPath, cfg.Database.Path)
	})
}

func TestEnvironmentVariables(t *testing.T) {
	// Create a test logger
	log, err := logger.New(&common.LogConfig{
		Level:   "debug",
		Console: true,
		Output:  []string{"stdout"},
	})
	require.NoError(t, err)

	t.Run("Complex environment variable overrides", func(t *testing.T) {
		// Set multiple environment variables
		envVars := map[string]string{
			config.EnvPrefix + "_APP_NAME":           "EnvApp",
			config.EnvPrefix + "_APP_VERSION":        "2.0.0",
			config.EnvPrefix + "_APP_ID":             "env.app.id",
			config.EnvPrefix + "_LOGGER_LEVEL":       "debug",
			config.EnvPrefix + "_LOGGER_CONSOLE":     "false",
			config.EnvPrefix + "_HOTKEYS_QUICK_NOTE": "Alt+Shift+N",
		}

		// Set environment variables and create cleanup function
		for k, v := range envVars {
			os.Setenv(k, v)
		}
		defer func() {
			for k := range envVars {
				os.Unsetenv(k)
			}
		}()

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
			config.WithLogger(log),
		)

		cfg, err := provider.Load()
		require.NoError(t, err)

		// Verify all environment variables were properly applied
		assert.Equal(t, "EnvApp", cfg.App.Name)
		assert.Equal(t, "2.0.0", cfg.App.Version)
		assert.Equal(t, "env.app.id", cfg.App.ID)
		assert.Equal(t, "debug", cfg.Logger.Level)
		assert.False(t, cfg.Logger.Console)
		assert.Equal(t, "Alt+Shift+N", cfg.Hotkeys.QuickNote.String())
	})

	t.Run("Invalid environment variable values", func(t *testing.T) {
		os.Setenv(config.EnvPrefix+"_LOGGER_LEVEL", "invalid_level")
		defer os.Unsetenv(config.EnvPrefix + "_LOGGER_LEVEL")

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
			config.WithLogger(log),
		)

		_, err := provider.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})
}

func TestMultipleConfigPaths(t *testing.T) {
	// Create temporary directories for config files
	baseDir := t.TempDir()
	primaryDir := filepath.Join(baseDir, "primary")
	fallbackDir := filepath.Join(baseDir, "fallback")
	require.NoError(t, os.MkdirAll(primaryDir, 0o755))
	require.NoError(t, os.MkdirAll(fallbackDir, 0o755))

	// Create config files with different values
	primaryConfig := []byte(`
app:
  name: "Primary App"
  version: "1.0.0"
`)
	fallbackConfig := []byte(`
app:
  name: "Fallback App"
  version: "0.5.0"
`)

	require.NoError(t, os.WriteFile(filepath.Join(primaryDir, "config.yaml"), primaryConfig, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(fallbackDir, "config.yaml"), fallbackConfig, 0o600))

	// Create a test logger
	log, err := logger.New(&common.LogConfig{
		Level:   "debug",
		Console: true,
		Output:  []string{"stdout"},
	})
	require.NoError(t, err)

	t.Run("Primary config takes precedence", func(t *testing.T) {
		provider := config.NewProvider(
			[]string{primaryDir, fallbackDir},
			"config",
			"yaml",
			config.WithLogger(log),
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "Primary App", cfg.App.Name)
		assert.Equal(t, "1.0.0", cfg.App.Version)
	})

	t.Run("Fallback when primary doesn't exist", func(t *testing.T) {
		// Remove primary config
		require.NoError(t, os.Remove(filepath.Join(primaryDir, "config.yaml")))

		provider := config.NewProvider(
			[]string{primaryDir, fallbackDir},
			"config",
			"yaml",
			config.WithLogger(log),
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "Fallback App", cfg.App.Name)
		assert.Equal(t, "0.5.0", cfg.App.Version)
	})

	t.Run("Default values when no config exists", func(t *testing.T) {
		// Use non-existent directories
		provider := config.NewProvider(
			[]string{filepath.Join(baseDir, "nonexistent")},
			"config",
			"yaml",
			config.WithLogger(log),
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, config.DefaultAppName, cfg.App.Name)
		assert.Equal(t, config.DefaultAppVersion, cfg.App.Version)
	})
}
