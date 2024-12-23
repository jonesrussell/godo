package logger

import (
	"fmt"
	"strings"

	"github.com/jonesrussell/godo/internal/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the minimum logging interface needed by your application
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// zapLogger implements Logger
type zapLogger struct {
	*zap.SugaredLogger
}

// New creates a new logger instance based on the provided configuration
func New(config *common.LogConfig) (Logger, error) {
	// Create Zap logger configuration
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(getZapLevel(config.Level)),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      config.Output,
		ErrorOutputPaths: config.ErrorOutput,
	}

	// Build the logger
	baseLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	// Return our wrapped logger
	return &zapLogger{
		SugaredLogger: baseLogger.Sugar(),
	}, nil
}

// getZapLevel converts string level to zapcore.Level
func getZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
