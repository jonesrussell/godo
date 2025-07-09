package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate mockgen -destination=../../test/mocks/mock_logger.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/logger Logger

// fieldMultiplier is used to calculate the capacity for key-value pairs
const fieldMultiplier = 2

type ZapLogger struct {
	*zap.SugaredLogger
}

// New creates a new logger instance based on the provided configuration
func New(config *LogConfig) (Logger, error) {
	// Validate log level
	level, err := parseLogLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}

	// Set up output paths
	outputPaths := config.Output
	if len(outputPaths) == 0 {
		outputPaths = []string{"stdout"}
	}

	// Set up error output paths
	errorOutputPaths := []string{"stderr"}

	// Handle file logging
	if config.File && config.FilePath != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(config.FilePath)
		if dirErr := os.MkdirAll(logDir, 0o755); dirErr != nil {
			return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, dirErr)
		}

		// Add file paths to outputs
		outputPaths = append(outputPaths, config.FilePath)
		errorOutputPaths = append(errorOutputPaths, config.FilePath)
	}

	// Create Zap logger configuration
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    getEncoderConfig(),
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
	}

	// Build the logger
	baseLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &ZapLogger{baseLogger.Sugar()}, nil
}

func parseLogLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unsupported log level: %s", level)
	}
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return encoderConfig
}

// Implement the interface methods
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.Fatalw(msg, keysAndValues...)
}

func (l *ZapLogger) WithError(err error) Logger {
	return &ZapLogger{l.With("error", err)}
}

func (l *ZapLogger) WithField(key string, value interface{}) Logger {
	return &ZapLogger{l.With(key, value)}
}

func (l *ZapLogger) WithFields(fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return l
	}
	args := make([]interface{}, 0, len(fields)*fieldMultiplier)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &ZapLogger{l.With(args...)}
}
