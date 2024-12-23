package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestLogger(t *testing.T) logger.Logger {
	t.Helper()
	log, err := logger.NewZapLogger(&logger.Config{
		Level:    "debug",
		Console:  true,
		File:     false,
		FilePath: "",
	})
	require.NoError(t, err)
	return log
}

func TestLoad(t *testing.T) {
	// Use the setupTestLogger function instead of inline creation
	log := setupTestLogger(t)

	// Test config loading
	cfg, err := Load(log)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Create test config files
	defaultConfig := `
app:
  name: "Godo"
  version: "0.1.0"
database:
  path: "test.db"
  max_open_conns: 1
  max_idle_conns: 1
logging:
  level: "info"
`
	err = os.WriteFile(filepath.Join(tmpDir, "default.yaml"), []byte(defaultConfig), 0o600)
	require.NoError(t, err)

	// Test environment config
	testConfig := `
database:
  path: "test_env.db"
  max_open_conns: 1
  max_idle_conns: 1
`
	err = os.WriteFile(filepath.Join(tmpDir, "test.yaml"), []byte(testConfig), 0o600)
	require.NoError(t, err)

	// Temporarily change working directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()

	// Create configs directory in temp directory
	configsDir := filepath.Join(tmpDir, "configs")
	err = os.Mkdir(configsDir, 0o755)
	require.NoError(t, err)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Copy test configs to configs directory
	err = os.WriteFile(filepath.Join(configsDir, "default.yaml"), []byte(defaultConfig), 0o600)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(configsDir, "test.yaml"), []byte(testConfig), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		env     string
		envVars map[string]string

		want    *Config
		wantErr bool
	}{
		{
			name: "loads default config",
			env:  "development",
			want: &Config{
				App: AppConfig{
					Name:    "Godo",
					Version: "0.1.0",
				},
				Database: DatabaseConfig{
					Path:         "test.db",
					MaxOpenConns: 1,
					MaxIdleConns: 1,
				},
			},
		},
		{
			name: "loads environment config",
			env:  "test",
			envVars: map[string]string{
				"GODO_ENV": "test",
			},
			want: &Config{
				App: AppConfig{
					Name:    "Godo",
					Version: "0.1.0",
				},
				Database: DatabaseConfig{
					Path:         "test_env.db",
					MaxOpenConns: 1,
					MaxIdleConns: 1,
				},
			},
		},
		{
			name: "applies environment variables",
			env:  "development",
			envVars: map[string]string{
				"GODO_DATABASE_PATH": "env.db",
				"GODO_LOGGING_LEVEL": "debug",
			},
			want: &Config{
				App: AppConfig{
					Name:    "Godo",
					Version: "0.1.0",
				},
				Database: DatabaseConfig{
					Path:         "env.db",
					MaxOpenConns: 1,
					MaxIdleConns: 1,
				},
				Logging: common.LogConfig{
					Level: "debug",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			t.Cleanup(func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			})

			// Test with logger instance
			got, err := Load(log)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				if tt.want != nil {
					assert.Equal(t, tt.want.App, got.App)
					assert.Equal(t, tt.want.Database, got.Database)
					if tt.envVars != nil && tt.envVars["GODO_LOGGING_LEVEL"] != "" {
						assert.Equal(t, tt.want.Logging.Level, got.Logging.Level)
					}
				}
			}
		})
	}
}
