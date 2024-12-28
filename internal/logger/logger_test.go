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
	"go.uber.org/zap/zaptest"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name      string
		config    *common.LogConfig
		wantLevel zapcore.Level
		wantErr   bool
	}{
		{
			name: "debug level",
			config: &common.LogConfig{
				Level: "debug",
			},
			wantLevel: zapcore.DebugLevel,
		},
		{
			name: "info level",
			config: &common.LogConfig{
				Level: "info",
			},
			wantLevel: zapcore.InfoLevel,
		},
		{
			name: "warn level",
			config: &common.LogConfig{
				Level: "warn",
			},
			wantLevel: zapcore.WarnLevel,
		},
		{
			name: "error level",
			config: &common.LogConfig{
				Level: "error",
			},
			wantLevel: zapcore.ErrorLevel,
		},
		{
			name: "invalid level defaults to info",
			config: &common.LogConfig{
				Level: "invalid",
			},
			wantLevel: zapcore.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewZapLogger(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logger)

			// Test that the logger can be used without panicking
			logger.Debug("test debug message")
			logger.Info("test info message")
			logger.Warn("test warn message")
			logger.Error("test error message")
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	var buf bytes.Buffer
	core := zaptest.NewLogger(t, zaptest.WrapOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(&buf),
			zapcore.DebugLevel,
		)
	}))).Core()
	logger := &ZapLogger{logger: zap.New(core).Sugar()}

	tests := []struct {
		name       string
		logFunc    func()
		wantMsg    string
		wantFields map[string]interface{}
	}{
		{
			name: "debug with fields",
			logFunc: func() {
				logger.Debug("debug message", "key", "value")
			},
			wantMsg: "debug message",
			wantFields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "info with fields",
			logFunc: func() {
				logger.Info("info message", "count", 42)
			},
			wantMsg: "info message",
			wantFields: map[string]interface{}{
				"count": float64(42), // JSON numbers are float64
			},
		},
		{
			name: "warn with fields",
			logFunc: func() {
				logger.Warn("warn message", "active", true)
			},
			wantMsg: "warn message",
			wantFields: map[string]interface{}{
				"active": true,
			},
		},
		{
			name: "error with fields",
			logFunc: func() {
				logger.Error("error message", "code", "E123")
			},
			wantMsg: "error message",
			wantFields: map[string]interface{}{
				"code": "E123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			var log map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &log)
			require.NoError(t, err)

			assert.Equal(t, tt.wantMsg, log["msg"])
			for key, want := range tt.wantFields {
				got, exists := log[key]
				assert.True(t, exists, "field %s should exist", key)
				assert.Equal(t, want, got, "field %s should have correct value", key)
			}
		})
	}
}

func TestLoggerWithFields(t *testing.T) {
	config := &common.LogConfig{Level: "debug"}
	logger, err := NewZapLogger(config)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test WithField
	fieldLogger := logger.WithField("key", "value")
	require.NotNil(t, fieldLogger)
	fieldLogger.Info("test message")

	// Test WithFields
	fields := map[string]interface{}{
		"string_field": "value",
		"int_field":    123,
		"bool_field":   true,
	}
	fieldsLogger := logger.WithFields(fields)
	require.NotNil(t, fieldsLogger)
	fieldsLogger.Info("test message")

	// Test WithError
	err = assert.AnError
	errorLogger := logger.WithError(err)
	require.NotNil(t, errorLogger)
	errorLogger.Error("test error")
}

func TestGetZapLevel(t *testing.T) {
	tests := []struct {
		level string
		want  zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"invalid", zapcore.InfoLevel},
		{"", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got := getZapLevel(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}
