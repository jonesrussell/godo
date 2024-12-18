package logger

import (
	"os"

	"github.com/jonesrussell/godo/internal/common"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals // logger needs to be globally accessible for application-wide logging
var log *zap.Logger

// Initialize sets up the logger with default configuration
func Initialize() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	var err error
	log, err = config.Build()
	if err != nil {
		return nil, err
	}
	return log, nil
}

// InitializeWithConfig sets up the logger with custom configuration
func InitializeWithConfig(cfg common.LogConfig) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	// Configure outputs
	if len(cfg.Output) > 0 {
		config.OutputPaths = cfg.Output
	}
	if len(cfg.ErrorOutput) > 0 {
		config.ErrorOutputPaths = cfg.ErrorOutput
	}

	// Set log level
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	config.Level = level

	log, err = config.Build()
	if err != nil {
		return nil, err
	}
	return log, nil
}

// Debug logs a debug message with structured fields
func Debug(msg string, keysAndValues ...interface{}) {
	log.Sugar().Debugw(msg, keysAndValues...)
}

// Info logs an info message with structured fields
func Info(msg string, keysAndValues ...interface{}) {
	log.Sugar().Infow(msg, keysAndValues...)
}

// Warn logs a warning message with structured fields
func Warn(msg string, keysAndValues ...interface{}) {
	log.Sugar().Warnw(msg, keysAndValues...)
}

// Error logs an error message with structured fields
func Error(msg string, keysAndValues ...interface{}) {
	log.Sugar().Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message with structured fields and exits
func Fatal(msg string, keysAndValues ...interface{}) {
	log.Sugar().Fatalw(msg, keysAndValues...)
	os.Exit(1)
}
