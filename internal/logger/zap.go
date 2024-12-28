package logger

import (
	"fmt"
	"strings"

	"github.com/jonesrussell/godo/internal/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	*zap.SugaredLogger
}

// New creates a new logger instance based on the provided configuration
func New(config *common.LogConfig) (Logger, error) {
	// Validate log level
	level, err := parseLogLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}

	// Create Zap logger configuration
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    getEncoderConfig(),
		OutputPaths:      config.Output,
		ErrorOutputPaths: config.ErrorOutput,
	}

	// Build the logger
	baseLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &zapLogger{baseLogger.Sugar()}, nil
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
func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) WithError(err error) Logger {
	return &zapLogger{l.SugaredLogger.With("error", err)}
}

func (l *zapLogger) WithField(key string, value interface{}) Logger {
	return &zapLogger{l.SugaredLogger.With(key, value)}
}

func (l *zapLogger) WithFields(fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return l
	}
	args := make([]interface{}, 0, len(fields)*fieldMultiplier)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &zapLogger{l.SugaredLogger.With(args...)}
}
