package config_test

import (
	"os"
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
			"test",
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
				os.Setenv("GODO_LOGGING_LEVEL", "invalid")
			},
			cleanupEnv: func() {
				os.Unsetenv("GODO_LOGGING_LEVEL")
			},
			expectError: true,
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
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}
