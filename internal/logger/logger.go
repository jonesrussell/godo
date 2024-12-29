// Package logger provides logging functionality for the application
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// fieldMultiplier is used to calculate the capacity for key-value pairs
const fieldMultiplier = 2

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	WithError(err error) Logger
}

// ZapLogger implements the Logger interface using zap
type ZapLogger struct {
	*zap.SugaredLogger
}

// Config holds logger configuration
type Config struct {
	Level       string
	Development bool
	Encoding    string
}

// NewLogger creates a new logger instance
func NewLogger(cfg *Config) (*ZapLogger, func(), error) {
	var zapCfg zap.Config

	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapCfg.Encoding = cfg.Encoding
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	if err := zapCfg.Level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, nil, fmt.Errorf("could not parse log level: %v", err)
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("could not build logger: %v", err)
	}

	return &ZapLogger{logger.Sugar()}, func() {
		_ = logger.Sync()
	}, nil
}

// Debug logs a debug message
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.Debugw(msg, keysAndValues...)
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.Infow(msg, keysAndValues...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.Warnw(msg, keysAndValues...)
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.Errorw(msg, keysAndValues...)
}

// WithError returns a logger with an error field
func (l *ZapLogger) WithError(err error) Logger {
	return &ZapLogger{l.With("error", err)}
}
