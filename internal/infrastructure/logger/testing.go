package logger

// TestLogger is a logger implementation specifically for testing
type TestLogger struct {
	T TestingT // Interface to support both *testing.T and *testing.B
}

// TestingT is an interface wrapper around *testing.T and *testing.B
type TestingT interface {
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

// NewTestLogger creates a new logger for testing that writes to the test log
func NewTestLogger(t TestingT) Logger {
	return &TestLogger{T: t}
}

// Debug logs a debug message with optional key-value pairs
func (l *TestLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("DEBUG: %s %v", msg, keysAndValues)
}

// Info logs an info message with optional key-value pairs
func (l *TestLogger) Info(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("INFO: %s %v", msg, keysAndValues)
}

// Warn logs a warning message with optional key-value pairs
func (l *TestLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("WARN: %s %v", msg, keysAndValues)
}

// Error logs an error message with optional key-value pairs
func (l *TestLogger) Error(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("ERROR: %s %v", msg, keysAndValues)
}

// Fatal logs a fatal message with optional key-value pairs and exits
func (l *TestLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("FATAL: %s %v", msg, keysAndValues)
}

// WithError returns a new logger with the error field set
func (l *TestLogger) WithError(_ error) Logger {
	return &TestLogger{T: l.T}
}

// WithField returns a new logger with the given field set
func (l *TestLogger) WithField(_ string, _ interface{}) Logger {
	return &TestLogger{T: l.T}
}

// WithFields returns a new logger with the given fields set
func (l *TestLogger) WithFields(_ map[string]interface{}) Logger {
	return &TestLogger{T: l.T}
}
