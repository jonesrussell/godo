package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	t.Run("Development", func(t *testing.T) {
		config := &common.LogConfig{
			Level:       "debug",
			Output:      []string{"stdout"},
			ErrorOutput: []string{"stderr"},
		}
		logger, err := New(config)
		require.NoError(t, err)
		require.NotNil(t, logger)

		// Verify logger is configured for development
		zapLogger := logger.(*zapLogger).SugaredLogger.Desugar()
		assert.True(t, zapLogger.Core().Enabled(zapcore.DebugLevel))
	})

	t.Run("Production", func(t *testing.T) {
		config := &common.LogConfig{
			Level:       "info",
			Output:      []string{"stdout"},
			ErrorOutput: []string{"stderr"},
		}
		logger, err := New(config)
		require.NoError(t, err)
		require.NotNil(t, logger)

		// Verify logger is configured for production
		zapLogger := logger.(*zapLogger).SugaredLogger.Desugar()
		assert.False(t, zapLogger.Core().Enabled(zapcore.DebugLevel))
	})
}

func TestLogLevels(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.DebugLevel)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	tests := []struct {
		name    string
		logFunc func(msg string, keysAndValues ...interface{})
		level   string
		message string
	}{
		{
			name:    "Debug",
			logFunc: logger.Debug,
			level:   "debug",
			message: "debug message",
		},
		{
			name:    "Info",
			logFunc: logger.Info,
			level:   "info",
			message: "info message",
		},
		{
			name:    "Warn",
			logFunc: logger.Warn,
			level:   "warn",
			message: "warning message",
		},
		{
			name:    "Error",
			logFunc: logger.Error,
			level:   "error",
			message: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message)

			var log map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &log)
			require.NoError(t, err)

			assert.Equal(t, tt.level, log["level"])
			assert.Equal(t, tt.message, log["msg"])
		})
	}
}

func TestNewTestLogger(t *testing.T) {
	testLogger := NewTestLogger(t)
	require.NotNil(t, testLogger)

	// Test all log levels
	testLogger.Debug("debug message")
	testLogger.Info("info message")
	testLogger.Warn("warning message")
	testLogger.Error("error message")
}

func TestWithFields(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.DebugLevel)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	// Create a logger with fields
	fields := map[string]interface{}{
		"string_field": "value",
		"int_field":    42,
		"bool_field":   true,
	}
	loggerWithFields := logger.WithFields(fields)

	// Log a message
	loggerWithFields.Info("test message")

	// Parse the log output
	var log map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &log)
	require.NoError(t, err)

	// Verify fields are present
	assert.Equal(t, "value", log["string_field"])
	assert.Equal(t, float64(42), log["int_field"]) // JSON numbers are float64
	assert.Equal(t, true, log["bool_field"])
	assert.Equal(t, "test message", log["msg"])
	assert.Equal(t, "info", log["level"])
}

func TestLoggerConcurrency(t *testing.T) {
	config := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}
	logger, err := New(config)
	require.NoError(t, err)

	// Log concurrently from multiple goroutines
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.WithFields(map[string]interface{}{
				"goroutine_id": id,
			}).Info("concurrent log message")
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
