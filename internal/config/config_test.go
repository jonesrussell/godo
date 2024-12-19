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
	// Initialize test logger
	log := setupTestLogger(t)

	// Create test config files
	tmpDir := t.TempDir()
	defaultConfig := `
app:
  name: "Godo"
  version: "0.1.0"
database:
  path: "test.db"
  max_open_conns: 1
  max_idle_conns: 1
`
	err := os.WriteFile(filepath.Join(tmpDir, "default.yaml"), []byte(defaultConfig), 0o600)
	require.NoError(t, err)

	// Test environment config
	testConfig := `
database:
  path: "test_env.db"
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

	err = os.Mkdir("configs", 0o755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join("configs", "default.yaml"), []byte(defaultConfig), 0o600)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join("configs", "test.yaml"), []byte(testConfig), 0o600)
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
				"GODO_LOG_LEVEL":     "debug",
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
			// Setup environment variables in a separate slice to handle cleanup
			var envToClean []string
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				envToClean = append(envToClean, k)
			}
			// Clean up environment variables after the test
			t.Cleanup(func() {
				for _, k := range envToClean {
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
					if tt.envVars != nil && tt.envVars["GODO_LOG_LEVEL"] != "" {
						assert.Equal(t, tt.want.Logging.Level, got.Logging.Level)
					}
				}
			}
		})
	}
}
