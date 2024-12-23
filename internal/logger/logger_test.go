package logger

import (
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *common.LogConfig
		wantErr bool
	}{
		{
			name: "creates with debug level",
			config: &common.LogConfig{
				Level:       "debug",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "creates with info level",
			config: &common.LogConfig{
				Level:       "info",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "fails with invalid level",
			config: &common.LogConfig{
				Level:       "invalid",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			logger, err := New(tt.config)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, logger)

				// Test logging
				logger.Info("test message")
			}
		})
	}
}

func TestLoggingFunctions(t *testing.T) {
	// Setup test logger
	config := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	logger, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test all logging methods
	logger.Debug("debug message", "key", "value")
	logger.Info("info message", "key", "value")
	logger.Warn("warn message", "key", "value")
	logger.Error("error message", "key", "value")
}
