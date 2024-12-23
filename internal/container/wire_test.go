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

func TestProviders(t *testing.T) {
	t.Run("provideLogger creates basic logger", func(t *testing.T) {
		log, err := provideLogger()
		require.NoError(t, err)
		assert.NotNil(t, log)
	})

	t.Run("provideSQLite creates store with correct path", func(t *testing.T) {
		// Create test config
		cfg := &config.Config{
			Database: config.DatabaseConfig{
				Path: filepath.Join(t.TempDir(), "test.db"),
			},
		}

		log, err := logger.NewZapLogger(&logger.Config{
			Level:   "debug",
			Console: true,
		})
		require.NoError(t, err)

		store, cleanup, err := provideSQLite(cfg, log)
		require.NoError(t, err)
		assert.NotNil(t, store)
		assert.NotNil(t, cleanup)

		cleanup()
	})

	t.Run("provideConfig loads configuration", func(t *testing.T) {
		// Create temporary test directory
		tmpDir := t.TempDir()
		configsDir := filepath.Join(tmpDir, "configs")
		require.NoError(t, os.MkdirAll(configsDir, 0o755))

		// Create test config file
		defaultConfig := `
app:
  name: "Godo"
  version: "0.1.0"
  id: "io.github.jonesrussell.godo"
database:
  path: "test.db"
logger:
  level: "info"
  console: true
hotkeys:
  quick_note: "Ctrl+Alt+G"
`
		err := os.WriteFile(filepath.Join(configsDir, "default.yaml"), []byte(defaultConfig), 0o600)
		require.NoError(t, err)

		// Set working directory to temp dir for config loading
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Chdir(originalWd))
		}()
		require.NoError(t, os.Chdir(tmpDir))

		cfg, err := provideConfig()
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "Godo", cfg.App.Name)
		assert.Equal(t, "0.1.0", cfg.App.Version)
	})
}

func TestInitializeApp(t *testing.T) {
	tests := []struct {
		name    string
		envVar  string
		setup   func(t *testing.T) (string, func())
		wantErr bool
	}{
		{
			name:   "initializes with default config",
			envVar: "",
			setup: func(t *testing.T) (string, func()) {
				tmpDir := t.TempDir()
				configsDir := filepath.Join(tmpDir, "configs")
				require.NoError(t, os.MkdirAll(configsDir, 0o755))

				defaultConfig := `
app:
  name: "Godo"
  version: "0.1.0"
  id: "io.github.jonesrussell.godo"
database:
  path: "test.db"
logger:
  level: "info"
  console: true
`
				err := os.WriteFile(filepath.Join(configsDir, "default.yaml"), []byte(defaultConfig), 0o600)
				require.NoError(t, err)

				originalWd, err := os.Getwd()
				require.NoError(t, err)

				require.NoError(t, os.Chdir(tmpDir))

				return originalWd, func() {
					require.NoError(t, os.Chdir(originalWd))
				}
			},
			wantErr: false,
		},
		{
			name:   "handles invalid config",
			envVar: "invalid",
			setup: func(t *testing.T) (string, func()) {
				tmpDir := t.TempDir()
				configsDir := filepath.Join(tmpDir, "configs")
				require.NoError(t, os.MkdirAll(configsDir, 0o755))

				invalidConfig := `invalid: yaml: content`
				err := os.WriteFile(filepath.Join(configsDir, "default.yaml"), []byte(invalidConfig), 0o600)
				require.NoError(t, err)

				originalWd, err := os.Getwd()
				require.NoError(t, err)

				require.NoError(t, os.Chdir(tmpDir))

				return originalWd, func() {
					require.NoError(t, os.Chdir(originalWd))
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.envVar != "" {
				os.Setenv("GODO_ENV", tt.envVar)
				defer os.Unsetenv("GODO_ENV")
			}

			originalWd, cleanup := tt.setup(t)
			defer cleanup()

			// Test
			app, cleanup2, err := InitializeApp()

			// Ensure cleanup runs
			t.Cleanup(func() {
				if cleanup2 != nil {
					cleanup2()
				}
				require.NoError(t, os.Chdir(originalWd))
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, app)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
			}
		})
	}
}
