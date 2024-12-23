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

// Implementation of Logger interface for TestLogger
func (l *TestLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("DEBUG: %s %v", msg, keysAndValues)
}

func (l *TestLogger) Info(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("INFO: %s %v", msg, keysAndValues)
}

func (l *TestLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("WARN: %s %v", msg, keysAndValues)
}

func (l *TestLogger) Error(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("ERROR: %s %v", msg, keysAndValues)
}

func (l *TestLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.T.Helper()
	l.T.Logf("FATAL: %s %v", msg, keysAndValues)
}

func (l *TestLogger) WithError(_ error) Logger {
	return &TestLogger{T: l.T}
}

func (l *TestLogger) WithField(_ string, _ interface{}) Logger {
	return &TestLogger{T: l.T}
}

func (l *TestLogger) WithFields(_ map[string]interface{}) Logger {
	return &TestLogger{T: l.T}
}
