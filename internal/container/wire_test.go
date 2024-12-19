package container

import (
	"os"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeApp(t *testing.T) {
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
		name string

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

			// Cleanup
			if cleanup != nil {
				cleanup()
			}
		})
	}
}
