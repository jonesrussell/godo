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

func TestLoggerOperations(t *testing.T) {
	tests := []struct {
		name       string
		config     *common.LogConfig
		logFunc    func(Logger)
		wantLevel  string
		wantFields map[string]interface{}
		wantErr    bool
	}{
		{
			name: "debug level logging",
			config: &common.LogConfig{
				Level:   "debug",
				Console: true,
			},
			logFunc: func(l Logger) {
				l.Debug("debug message", "key", "value")
			},
			wantLevel: "debug",
			wantFields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "info level logging",
			config: &common.LogConfig{
				Level:   "info",
				Console: true,
			},
			logFunc: func(l Logger) {
				l.Info("info message", "count", 42)
			},
			wantLevel: "info",
			wantFields: map[string]interface{}{
				"count": float64(42), // JSON numbers are float64
			},
		},
		{
			name: "warn level logging",
			config: &common.LogConfig{
				Level:   "warn",
				Console: true,
			},
			logFunc: func(l Logger) {
				l.Warn("warn message", "active", true)
			},
			wantLevel: "warn",
			wantFields: map[string]interface{}{
				"active": true,
			},
		},
		{
			name: "error level logging",
			config: &common.LogConfig{
				Level:   "error",
				Console: true,
			},
			logFunc: func(l Logger) {
				l.Error("error message", "error_code", "E123")
			},
			wantLevel: "error",
			wantFields: map[string]interface{}{
				"error_code": "E123",
			},
		},
		{
			name: "invalid log level",
			config: &common.LogConfig{
				Level:   "invalid",
				Console: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				_, err := New(tt.config)
				assert.Error(t, err)
				return
			}

			// Create a buffer to capture log output
			var buf bytes.Buffer
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(&buf),
				zap.NewAtomicLevelAt(zapcore.DebugLevel),
			)
			testLogger := &zapLogger{
				SugaredLogger: zap.New(core).Sugar(),
			}

			// Execute the log function
			tt.logFunc(testLogger)

			// Parse the log output
			var logEntry map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &logEntry)
			require.NoError(t, err)

			// Verify log level
			assert.Equal(t, tt.wantLevel, logEntry["level"])

			// Verify log fields
			for key, want := range tt.wantFields {
				got, exists := logEntry[key]
				assert.True(t, exists, "field %s should exist", key)
				assert.Equal(t, want, got, "field %s should have correct value", key)
			}
		})
	}
}

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger(t)
	assert.NotNil(t, logger)

	// Test all log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestLoggerWithOptions(t *testing.T) {
	tests := []struct {
		name    string
		config  *common.LogConfig
		wantErr bool
	}{
		{
			name: "with file output",
			config: &common.LogConfig{
				Level:   "debug",
				Console: false,
				Output:  []string{"test.log"},
			},
			wantErr: false,
		},
		{
			name: "with console and file output",
			config: &common.LogConfig{
				Level:   "info",
				Console: true,
				Output:  []string{"test.log"},
			},
			wantErr: false,
		},
		{
			name: "with invalid file path",
			config: &common.LogConfig{
				Level:   "info",
				Console: false,
				Output:  []string{"\x00invalid\x00path\x00test.log"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, logger)

			// Test logging
			logger.Info("test message")
		})
	}
}
