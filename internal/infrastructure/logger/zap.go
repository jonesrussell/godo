package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate mockgen -destination=../../test/mocks/mock_logger.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/logger Logger

const fieldMultiplier = 2

type ZapLogger struct {
	*zap.SugaredLogger
	logFile    io.Closer
	mu         sync.RWMutex
	closed     bool
	hasConsole bool
}

// New creates a new logger instance based on the provided configuration
func New(config *LogConfig) (Logger, error) {
	if config == nil {
		return nil, fmt.Errorf("logger config cannot be nil")
	}

	level, err := parseLogLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}

	writers, errorWriters, logFile, hasConsole, err := setupWriters(config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup writers: %w", err)
	}

	core := createCore(writers, errorWriters, level)
	baseLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		SugaredLogger: baseLogger.Sugar(),
		logFile:       logFile,
		hasConsole:    hasConsole,
	}, nil
}

func setupWriters(config *LogConfig) ([]zapcore.WriteSyncer, []zapcore.WriteSyncer, io.Closer, bool, error) {
	var writers []zapcore.WriteSyncer
	var errorWriters []zapcore.WriteSyncer
	var logFile io.Closer
	hasConsole := false

	// Add console output if enabled
	if config.Console {
		// Create safe console writers that won't fail on Windows during shutdown
		stdoutWriter := newSafeConsoleWriter(os.Stdout)
		stderrWriter := newSafeConsoleWriter(os.Stderr)

		writers = append(writers, stdoutWriter)
		errorWriters = append(errorWriters, stderrWriter)
		hasConsole = true
	}

	// Handle file logging
	if config.File && config.FilePath != "" {
		file, err := setupFileLogging(config.FilePath)
		if err != nil {
			return nil, nil, nil, false, err
		}

		logFile = file
		fileSync := zapcore.AddSync(file)
		writers = append(writers, fileSync)
		errorWriters = append(errorWriters, fileSync)
	}

	// Default to stdout if no writers specified
	if len(writers) == 0 {
		stdoutWriter := newSafeConsoleWriter(os.Stdout)
		stderrWriter := newSafeConsoleWriter(os.Stderr)
		writers = append(writers, stdoutWriter)
		errorWriters = append(errorWriters, stderrWriter)
		hasConsole = true
	}

	return writers, errorWriters, logFile, hasConsole, nil
}

// safeConsoleWriter wraps console writers to handle Windows sync failures gracefully
type safeConsoleWriter struct {
	writer   io.Writer
	mu       sync.RWMutex
	disabled bool
}

func newSafeConsoleWriter(w io.Writer) zapcore.WriteSyncer {
	return &safeConsoleWriter{writer: w}
}

func (w *safeConsoleWriter) Write(p []byte) (n int, err error) {
	w.mu.RLock()
	disabled := w.disabled
	w.mu.RUnlock()

	if disabled {
		return len(p), nil // Pretend write succeeded
	}

	n, err = w.writer.Write(p)
	if err != nil && isConsoleInvalidError(err) {
		// Disable future writes to prevent spam
		w.mu.Lock()
		w.disabled = true
		w.mu.Unlock()
		return len(p), nil // Pretend write succeeded
	}
	return n, err
}

func (w *safeConsoleWriter) Sync() error {
	w.mu.RLock()
	disabled := w.disabled
	w.mu.RUnlock()

	if disabled {
		return nil
	}

	// On Windows, syncing stdout/stderr can fail during shutdown
	if syncer, ok := w.writer.(interface{ Sync() error }); ok {
		err := syncer.Sync()
		if err != nil && isConsoleInvalidError(err) {
			// Disable future operations
			w.mu.Lock()
			w.disabled = true
			w.mu.Unlock()
			return nil // Don't propagate the error
		}
		return err
	}
	return nil
}

// isConsoleInvalidError checks if the error is related to invalid console handles
func isConsoleInvalidError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "invalid") ||
		strings.Contains(errStr, "handle is invalid") ||
		strings.Contains(errStr, "/dev/stdout") ||
		strings.Contains(errStr, "/dev/stderr")
}

func setupFileLogging(filePath string) (*os.File, error) {
	logDir := filepath.Dir(filePath)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, err)
	}

	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}

	return logFile, nil
}

func createCore(writers, errorWriters []zapcore.WriteSyncer, level zapcore.Level) zapcore.Core {
	encoderConfig := getEncoderConfig()

	// Use JSON encoder for file logging, console encoder for console
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if runtime.GOOS == "windows" {
		// On Windows, use a more conservative encoder config
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // No colors
	}

	// Main core for all levels
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		zap.NewAtomicLevelAt(level),
	)

	// Error core for error-level messages
	errorCore := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(errorWriters...),
		zap.NewAtomicLevelAt(zapcore.ErrorLevel),
	)

	return zapcore.NewTee(core, errorCore)
}

func parseLogLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unsupported log level: %s", level)
	}
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Use colored output only on non-Windows or when explicitly supported
	if runtime.GOOS != "windows" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	return encoderConfig
}

// Interface implementations with closed state checking
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	if l.isClosed() {
		return
	}
	l.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	if l.isClosed() {
		return
	}
	l.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	if l.isClosed() {
		return
	}
	l.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	if l.isClosed() {
		return
	}
	l.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	if l.isClosed() {
		return
	}
	l.Fatalw(msg, keysAndValues...)
}

func (l *ZapLogger) WithError(err error) Logger {
	if l.isClosed() {
		return l
	}
	return &ZapLogger{
		SugaredLogger: l.With("error", err),
		logFile:       l.logFile,
		hasConsole:    l.hasConsole,
	}
}

func (l *ZapLogger) WithField(key string, value interface{}) Logger {
	if l.isClosed() {
		return l
	}
	return &ZapLogger{
		SugaredLogger: l.With(key, value),
		logFile:       l.logFile,
		hasConsole:    l.hasConsole,
	}
}

func (l *ZapLogger) WithFields(fields map[string]interface{}) Logger {
	if l.isClosed() || len(fields) == 0 {
		return l
	}

	args := make([]interface{}, 0, len(fields)*fieldMultiplier)
	for k, v := range fields {
		args = append(args, k, v)
	}

	return &ZapLogger{
		SugaredLogger: l.With(args...),
		logFile:       l.logFile,
		hasConsole:    l.hasConsole,
	}
}

// Sync flushes any buffered log entries with Windows-safe error handling
func (l *ZapLogger) Sync() error {
	if l.isClosed() {
		return nil
	}

	if logger := l.SugaredLogger.Desugar(); logger != nil {
		err := logger.Sync()
		// On Windows during shutdown, sync errors are common and expected
		if err != nil && isConsoleInvalidError(err) {
			// Log to stderr if possible, but don't fail
			if runtime.GOOS == "windows" {
				return nil // Ignore console sync errors on Windows
			}
		}
		return err
	}
	return nil
}

// Close properly closes the logger and its resources
func (l *ZapLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}

	// Sync file outputs only (skip console on Windows to avoid errors)
	if l.logFile != nil {
		if logger := l.SugaredLogger.Desugar(); logger != nil {
			_ = logger.Sync() // Best effort sync, ignore errors during shutdown
		}
	}

	// Close log file if it exists
	if l.logFile != nil {
		if err := l.logFile.Close(); err != nil {
			l.closed = true
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}

	l.closed = true
	return nil
}

func (l *ZapLogger) isClosed() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.closed
}
