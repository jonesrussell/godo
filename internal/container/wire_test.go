package container

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeApp(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	configsDir := filepath.Join(tmpDir, "configs")
	require.NoError(t, os.MkdirAll(configsDir, 0o755))

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
	err := os.WriteFile(filepath.Join(configsDir, "default.yaml"), []byte(defaultConfig), 0o600)
	require.NoError(t, err)

	// Create production config
	productionConfig := `
database:
  path: "prod.db"
logging:
  level: "warn"
`
	err = os.WriteFile(filepath.Join(configsDir, "production.yaml"), []byte(productionConfig), 0o600)
	require.NoError(t, err)

	// Set working directory to temp dir for config loading
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(originalWd))
	}()
	require.NoError(t, os.Chdir(tmpDir))

	// Initialize logger with test config
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	log, err := logger.New(logConfig)
	require.NoError(t, err, "Failed to create logger")
	require.NotNil(t, log)

	tests := []struct {
		name    string
		envVar  string
		wantErr bool
	}{
		{
			name:    "initializes with development environment",
			envVar:  "",
			wantErr: false,
		},
		{
			name:    "initializes with production environment",
			envVar:  "production",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.envVar != "" {
				os.Setenv("GODO_ENV", tt.envVar)
				defer os.Unsetenv("GODO_ENV")
			} else {
				os.Unsetenv("GODO_ENV")
			}

			// Test
			app, cleanup, err := InitializeApp()

			// Ensure cleanup runs at the end of each test
			t.Cleanup(func() {
				if cleanup != nil {
					cleanup()
				}
			})

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, app)
				assert.Nil(t, cleanup)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, app)
				assert.NotNil(t, cleanup)
			}
		})
	}
}
