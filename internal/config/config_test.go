package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("Load default config", func(t *testing.T) {
		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "Godo", cfg.App.Name)
		assert.Equal(t, "0.1.0", cfg.App.Version)
	})

	t.Run("Load test environment config", func(t *testing.T) {
		// Set up test environment
		os.Setenv("GODO_DATABASE_PATH", "test.db")
		defer os.Unsetenv("GODO_DATABASE_PATH")

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "test.db", cfg.Database.Path)
	})

	t.Run("Config file not found uses defaults", func(t *testing.T) {
		provider := config.NewProvider(
			[]string{"nonexistent"},
			"nonexistent",
			"yaml",
		)

		cfg, err := provider.Load()
		require.NoError(t, err)
		assert.Equal(t, "Godo", cfg.App.Name)
		assert.Equal(t, "godo.db", cfg.Database.Path)
	})
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		expectError bool
	}{
		{
			name: "valid database path",
			setupEnv: func() {
				os.Setenv("GODO_DATABASE_PATH", "valid.db")
			},
			cleanupEnv: func() {
				os.Unsetenv("GODO_DATABASE_PATH")
			},
			expectError: false,
		},
		{
			name: "invalid log level",
			setupEnv: func() {
				os.Setenv("GODO_LOGGER_LEVEL", "invalid")
				os.Setenv("GODO_CONFIG", "testdata/config.yaml")
			},
			cleanupEnv: func() {
				os.Unsetenv("GODO_LOGGER_LEVEL")
				os.Unsetenv("GODO_CONFIG")
			},
			expectError: true,
		},
		{
			name: "missing required fields",
			setupEnv: func() {
				os.Setenv("GODO_APP_NAME", "")
			},
			cleanupEnv: func() {
				os.Unsetenv("GODO_APP_NAME")
			},
			expectError: true,
		},
		{
			name: "valid log level - debug",
			setupEnv: func() {
				os.Setenv("GODO_LOGGER_LEVEL", "debug")
			},
			cleanupEnv: func() {
				os.Unsetenv("GODO_LOGGER_LEVEL")
			},
			expectError: false,
		},
		{
			name: "valid log level - warn",
			setupEnv: func() {
				os.Setenv("GODO_LOGGER_LEVEL", "warn")
			},
			cleanupEnv: func() {
				os.Unsetenv("GODO_LOGGER_LEVEL")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				tt.setupEnv()
			}
			defer func() {
				if tt.cleanupEnv != nil {
					tt.cleanupEnv()
				}
			}()

			provider := config.NewProvider(
				[]string{"testdata"},
				"config",
				"yaml",
			)

			cfg, err := provider.Load()
			if tt.expectError {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), "invalid log level")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
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

func TestConfigFileErrors(t *testing.T) {
	t.Run("Invalid YAML syntax", func(t *testing.T) {
		// Create temporary file with invalid YAML
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

func TestConfigEnvBindingErrors(t *testing.T) {
	t.Run("Invalid environment variable binding", func(t *testing.T) {
		// Set an invalid environment variable
		os.Setenv("GODO_INVALID_VAR", "invalid")
		defer os.Unsetenv("GODO_INVALID_VAR")

		provider := config.NewProvider(
			[]string{"testdata"},
			"config",
			"yaml",
		)

		cfg, err := provider.Load()
		require.NoError(t, err) // Should not fail as invalid vars are ignored
		assert.NotNil(t, cfg)
	})
}

func TestConfigValidationRules(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &config.Config{
				Logger: config.LoggerConfig{
					Level: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "empty log level",
			cfg: &config.Config{
				Logger: config.LoggerConfig{
					Level: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.ValidateConfig(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
