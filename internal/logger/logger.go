package logger

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log        *zap.Logger
	sugar      *zap.SugaredLogger
	isUIActive atomic.Bool
)

func init() {
	// Ensure logs directory exists in project root
	logsDir := filepath.Join("logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic("failed to create logs directory: " + err.Error())
	}

	// Configure logging to file
	logFile := filepath.Join(logsDir, "godo.log")
	writer, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	// Custom encoder config
	config := zapcore.EncoderConfig{
		TimeKey:    "T",
		LevelKey:   "L",
		MessageKey: "M",
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(l.CapitalString())
		},
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05-07:00"))
		},
	}

	fileEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(writer), zapcore.DebugLevel),
	)

	log = zap.New(core)
	sugar = log.Sugar()
}

// SetUIActive sets the UI active state
func SetUIActive(active bool) {
	isUIActive.Store(active)
}

// IsUIActive returns whether the UI is currently active
func IsUIActive() bool {
	return isUIActive.Load()
}

func Info(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func Error(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

func Debug(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}

func Warn(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

// WithFields adds structured context to the logger
func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	return sugar.With(fields)
}

// Sync flushes any buffered log entries
func Sync() error {
	return sugar.Sync()
}
