package logger

import (
	"errors"
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

func TestLoggingWithContext(t *testing.T) {
	logger := NewTestLogger()

	t.Run("WithField", func(t *testing.T) {
		contextLogger := logger.WithField("requestID", "123")
		require.NotNil(t, contextLogger)
		contextLogger.Info("test message")
	})

	t.Run("WithFields", func(t *testing.T) {
		fields := map[string]interface{}{
			"requestID": "123",
			"userID":    "456",
		}
		contextLogger := logger.WithFields(fields)
		require.NotNil(t, contextLogger)
		contextLogger.Info("test message")
	})

	t.Run("WithError", func(t *testing.T) {
		err := errors.New("test error")
		contextLogger := logger.WithError(err)
		require.NotNil(t, contextLogger)
		contextLogger.Error("operation failed")
	})
}

func TestNoopLogger(t *testing.T) {
	logger := NewNoopLogger()

	// These should not panic
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	// Test context methods
	withField := logger.WithField("key", "value")
	assert.NotNil(t, withField)

	withFields := logger.WithFields(map[string]interface{}{"key": "value"})
	assert.NotNil(t, withFields)

	withError := logger.WithError(errors.New("test"))
	assert.NotNil(t, withError)
}
