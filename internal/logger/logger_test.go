package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *common.LogConfig
		wantErr bool
	}{
		{
			name: "valid debug level",
			config: &common.LogConfig{
				Level:       "debug",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "valid info level",
			config: &common.LogConfig{
				Level:       "info",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "valid warn level",
			config: &common.LogConfig{
				Level:       "warn",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "valid error level",
			config: &common.LogConfig{
				Level:       "error",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "invalid level",
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
			logger, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, logger)

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
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

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

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger(t)
	require.NotNil(t, logger)

	// Test all log levels
	logger.Debug("debug message", "key", "value")
	logger.Info("info message", "key", "value")
	logger.Warn("warn message", "key", "value")
	logger.Error("error message", "key", "value", "error", "test error")
}

func TestNoopLogger(t *testing.T) {
	logger := NewNoopLogger()
	require.NotNil(t, logger)

	// Test that none of these operations panic
	assert.NotPanics(t, func() {
		logger.Debug("debug message", "key", "value")
		logger.Info("info message", "key", "value")
		logger.Warn("warn message", "key", "value")
		logger.Error("error message", "key", "value")
		logger.WithError(assert.AnError)
		logger.WithField("key", "value")
		logger.WithFields(map[string]interface{}{"key": "value"})
	})
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		level   string
		want    zapcore.Level
		wantErr bool
	}{
		{"debug", zapcore.DebugLevel, false},
		{"info", zapcore.InfoLevel, false},
		{"warn", zapcore.WarnLevel, false},
		{"error", zapcore.ErrorLevel, false},
		{"invalid", zapcore.InfoLevel, true},
		{"", zapcore.InfoLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got, err := parseLogLevel(tt.level)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetEncoderConfig(t *testing.T) {
	config := getEncoderConfig()

	// Verify encoder config settings
	assert.NotNil(t, config.EncodeTime)
	assert.NotNil(t, config.EncodeLevel)
	assert.NotNil(t, config.EncodeDuration)
	assert.NotNil(t, config.EncodeCaller)
}

func TestLoggerWithMethods(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	tests := []struct {
		name       string
		logFunc    func()
		wantFields map[string]interface{}
	}{
		{
			name: "WithError",
			logFunc: func() {
				logger.WithError(assert.AnError).Error("error message")
			},
			wantFields: map[string]interface{}{
				"error": assert.AnError.Error(),
			},
		},
		{
			name: "WithField",
			logFunc: func() {
				logger.WithField("key", "value").Info("info message")
			},
			wantFields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "WithFields",
			logFunc: func() {
				fields := map[string]interface{}{
					"string_key": "value",
					"int_key":    42,
					"bool_key":   true,
				}
				logger.WithFields(fields).Info("info message")
			},
			wantFields: map[string]interface{}{
				"string_key": "value",
				"int_key":    float64(42), // JSON numbers are float64
				"bool_key":   true,
			},
		},
		{
			name: "WithFields empty map",
			logFunc: func() {
				logger.WithFields(nil).Info("info message")
			},
			wantFields: map[string]interface{}{},
		},
		{
			name: "WithField nil value",
			logFunc: func() {
				logger.WithField("key", nil).Info("info message")
			},
			wantFields: map[string]interface{}{
				"key": nil,
			},
		},
		{
			name: "WithError nil error",
			logFunc: func() {
				logger.WithError(nil).Error("error message")
			},
			wantFields: map[string]interface{}{},
		},
		{
			name: "chained WithField calls",
			logFunc: func() {
				logger.WithField("key1", "value1").
					WithField("key2", "value2").
					Info("info message")
			},
			wantFields: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "WithField after WithFields",
			logFunc: func() {
				logger.WithFields(map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				}).WithField("key3", "value3").Info("info message")
			},
			wantFields: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "WithError after WithField",
			logFunc: func() {
				logger.WithField("key", "value").
					WithError(assert.AnError).
					Error("error message")
			},
			wantFields: map[string]interface{}{
				"key":   "value",
				"error": assert.AnError.Error(),
			},
		},
		{
			name: "WithFields after WithError",
			logFunc: func() {
				logger.WithError(assert.AnError).
					WithFields(map[string]interface{}{
						"key": "value",
					}).Error("error message")
			},
			wantFields: map[string]interface{}{
				"error": assert.AnError.Error(),
				"key":   "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			// Split the buffer into lines and process each line
			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			require.NotEmpty(t, lines)

			// Parse the last line (most recent log entry)
			var log map[string]interface{}
			err := json.Unmarshal([]byte(lines[len(lines)-1]), &log)
			require.NoError(t, err)

			// Verify fields
			for key, want := range tt.wantFields {
				got, exists := log[key]
				assert.True(t, exists, "field %s should exist", key)
				assert.Equal(t, want, got, "field %s should have correct value", key)
			}
		})
	}
}

func TestLoggerEdgeCases(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	tests := []struct {
		name    string
		logFunc func()
	}{
		{
			name: "empty message",
			logFunc: func() {
				logger.Info("")
			},
		},
		{
			name: "nil fields",
			logFunc: func() {
				logger.Info("message", nil)
			},
		},
		{
			name: "odd number of fields",
			logFunc: func() {
				logger.Info("message", "key1", "value1", "key2")
			},
		},
		{
			name: "non-string key",
			logFunc: func() {
				logger.Info("message", 42, "value")
			},
		},
		{
			name: "complex value",
			logFunc: func() {
				logger.Info("message", "key", struct{ Value string }{"test"})
			},
		},
		{
			name: "multiple WithField calls",
			logFunc: func() {
				logger.WithField("key1", "value1").
					WithField("key2", "value2").
					Info("message")
			},
		},
		{
			name: "WithField after WithFields",
			logFunc: func() {
				logger.WithFields(map[string]interface{}{"key1": "value1"}).
					WithField("key2", "value2").
					Info("message")
			},
		},
		{
			name: "WithError after WithField",
			logFunc: func() {
				logger.WithField("key", "value").
					WithError(assert.AnError).
					Error("message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			assert.NotPanics(t, tt.logFunc)

			// Verify that something was logged
			assert.NotEmpty(t, buf.String())
		})
	}
}

func TestLoggerConcurrent(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	const numGoroutines = 10
	const logsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Use a mutex to protect the buffer
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				mu.Lock()
				logger.WithFields(map[string]interface{}{
					"goroutine": id,
					"iteration": j,
				}).Info("concurrent log")
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Count the number of log entries by counting newlines
	logContent := buf.String()
	logLines := strings.Split(logContent, "\n")
	var validLines int
	for _, line := range logLines {
		if strings.TrimSpace(line) != "" {
			validLines++
		}
	}
	assert.Equal(t, numGoroutines*logsPerGoroutine, validLines)
}

func TestLoggerConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *common.LogConfig
		wantErr bool
	}{
		{
			name: "valid config with console output",
			config: &common.LogConfig{
				Level:       "debug",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "valid config with file output",
			config: &common.LogConfig{
				Level:       "info",
				Output:      []string{"test.log"},
				ErrorOutput: []string{"test.error.log"},
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: &common.LogConfig{
				Level:       "invalid",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantErr: true,
		},
		{
			name: "empty output paths",
			config: &common.LogConfig{
				Level:       "info",
				Output:      []string{},
				ErrorOutput: []string{},
			},
			wantErr: false,
		},
		{
			name: "nil config",
			config: &common.LogConfig{
				Level: "info",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, logger)

			// Test that the logger can be used
			logger.Info("test message")
		})
	}
}

func TestLoggerEncoderConfig(t *testing.T) {
	config := getEncoderConfig()

	// Test encoder config settings
	assert.Equal(t, zapcore.ISO8601TimeEncoder, config.EncodeTime)
	assert.Equal(t, zapcore.CapitalColorLevelEncoder, config.EncodeLevel)
	assert.Equal(t, zapcore.StringDurationEncoder, config.EncodeDuration)
	assert.Equal(t, zapcore.ShortCallerEncoder, config.EncodeCaller)

	// Test encoder config with sample data
	encoder := zapcore.NewJSONEncoder(config)
	buf, err := encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.InfoLevel,
		Time:       time.Now(),
		LoggerName: "test",
		Message:    "test message",
		Caller:     zapcore.EntryCaller{Defined: true, File: "test.go", Line: 42},
	}, []zapcore.Field{
		zap.String("key", "value"),
	})
	require.NoError(t, err)

	var log map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &log)
	require.NoError(t, err)

	assert.Contains(t, log, "ts")
	assert.Contains(t, log, "level")
	assert.Contains(t, log, "logger")
	assert.Contains(t, log, "caller")
	assert.Contains(t, log, "msg")
	assert.Contains(t, log, "key")
}

func TestLoggerLevelParsing(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		want    zapcore.Level
		wantErr bool
	}{
		{
			name:    "debug level",
			level:   "debug",
			want:    zapcore.DebugLevel,
			wantErr: false,
		},
		{
			name:    "info level",
			level:   "info",
			want:    zapcore.InfoLevel,
			wantErr: false,
		},
		{
			name:    "warn level",
			level:   "warn",
			want:    zapcore.WarnLevel,
			wantErr: false,
		},
		{
			name:    "error level",
			level:   "error",
			want:    zapcore.ErrorLevel,
			wantErr: false,
		},
		{
			name:    "invalid level",
			level:   "invalid",
			want:    zapcore.InfoLevel,
			wantErr: true,
		},
		{
			name:    "empty level",
			level:   "",
			want:    zapcore.InfoLevel,
			wantErr: true,
		},
		{
			name:    "mixed case level",
			level:   "InFo",
			want:    zapcore.InfoLevel,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogLevel(tt.level)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestLoggerWithFieldsNilMap(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	// Test with nil map
	logger.WithFields(nil).Info("test message")

	var log map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &log)
	require.NoError(t, err)
	assert.Equal(t, "test message", log["msg"])
}

func TestLoggerWithFieldsEmptyMap(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	// Test with empty map
	logger.WithFields(map[string]interface{}{}).Info("test message")

	var log map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &log)
	require.NoError(t, err)
	assert.Equal(t, "test message", log["msg"])
}

func TestLoggerWithFieldsNestedMap(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	// Test with nested map
	logger.WithFields(map[string]interface{}{
		"nested": map[string]interface{}{
			"key": "value",
		},
	}).Info("test message")

	var log map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &log)
	require.NoError(t, err)
	assert.Equal(t, "test message", log["msg"])
	nested, ok := log["nested"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", nested["key"])
}

func TestLoggerWithFieldsComplexTypes(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &zapLogger{SugaredLogger: zap.New(core).Sugar()}

	type testStruct struct {
		Field string
	}

	// Test with various types
	logger.WithFields(map[string]interface{}{
		"int":      42,
		"float":    3.14,
		"bool":     true,
		"string":   "value",
		"slice":    []string{"a", "b", "c"},
		"struct":   testStruct{Field: "value"},
		"map":      map[string]string{"key": "value"},
		"nil":      nil,
		"duration": time.Second,
		"time":     time.Now(),
	}).Info("test message")

	var log map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &log)
	require.NoError(t, err)
	assert.Equal(t, "test message", log["msg"])
	assert.Equal(t, float64(42), log["int"])
	assert.Equal(t, 3.14, log["float"])
	assert.Equal(t, true, log["bool"])
	assert.Equal(t, "value", log["string"])
	assert.NotNil(t, log["slice"])
	assert.NotNil(t, log["struct"])
	assert.NotNil(t, log["map"])
	assert.Nil(t, log["nil"])
	assert.NotNil(t, log["duration"])
	assert.NotNil(t, log["time"])
}

func TestLoggerCleanup(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := os.CreateTemp("", "test-*.log")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a logger with the temporary file
	config := &common.LogConfig{
		Level:       "info",
		Output:      []string{tmpFile.Name()},
		ErrorOutput: []string{tmpFile.Name()},
	}

	logger, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Write some logs
	logger.Info("test message")
	logger.Error("test error")

	// Verify that the file contains the logs
	content, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	assert.Contains(t, string(content), "test message")
	assert.Contains(t, string(content), "test error")
}

func TestLoggerWithInvalidFile(t *testing.T) {
	// Try to create a logger with an invalid file path
	config := &common.LogConfig{
		Level:       "info",
		Output:      []string{"/invalid/path/test.log"},
		ErrorOutput: []string{"/invalid/path/test.error.log"},
	}

	logger, err := New(config)
	assert.Error(t, err)
	assert.Nil(t, logger)
}
