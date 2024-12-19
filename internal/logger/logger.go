package logger

import (
	"sync"

	"github.com/jonesrussell/godo/internal/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

type zapLogger struct {
	log *zap.Logger
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Fatalw(msg, keysAndValues...)
}

type loggerManager struct {
	sync.RWMutex
	logger Logger
}

var (
	manager *loggerManager
	once    sync.Once
)

func getManager() *loggerManager {
	once.Do(func() {
		manager = &loggerManager{}
	})
	return manager
}

// Initialize creates and sets up the logger
func Initialize(config *common.LogConfig) (Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      config.Output,
		ErrorOutputPaths: config.ErrorOutput,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	zl := &zapLogger{log: logger}
	m := getManager()
	m.Lock()
	m.logger = zl
	m.Unlock()
	return zl, nil
}

// Package-level functions use thread-safe manager
func Debug(msg string, keysAndValues ...interface{}) {
	m := getManager()
	m.RLock()
	defer m.RUnlock()
	if m.logger != nil {
		m.logger.Debug(msg, keysAndValues...)
	}
}

func Info(msg string, keysAndValues ...interface{}) {
	m := getManager()
	m.RLock()
	defer m.RUnlock()
	if m.logger != nil {
		m.logger.Info(msg, keysAndValues...)
	}
}

func Warn(msg string, keysAndValues ...interface{}) {
	m := getManager()
	m.RLock()
	defer m.RUnlock()
	if m.logger != nil {
		m.logger.Warn(msg, keysAndValues...)
	}
}

func Error(msg string, keysAndValues ...interface{}) {
	m := getManager()
	m.RLock()
	defer m.RUnlock()
	if m.logger != nil {
		m.logger.Error(msg, keysAndValues...)
	}
}

func Fatal(msg string, keysAndValues ...interface{}) {
	m := getManager()
	m.RLock()
	defer m.RUnlock()
	if m.logger != nil {
		m.logger.Fatal(msg, keysAndValues...)
	}
}
