package logger

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockTestLogger combines both test logger and mock functionality
type MockTestLogger struct {
	Logger
	mock.Mock
	t *testing.T
}

// NewMockTestLogger creates a new logger that combines test logging and mock expectations
func NewMockTestLogger(t *testing.T) *MockTestLogger {
	l := &MockTestLogger{
		Logger: NewTestLogger(t),
		t:      t,
	}
	// Set up default expectations
	l.On("Debug", mock.Anything, mock.Anything).Return()
	l.On("Info", mock.Anything, mock.Anything).Return()
	l.On("Error", mock.Anything, mock.Anything).Return()
	l.On("Warn", mock.Anything, mock.Anything).Return()
	l.On("WithError", mock.Anything).Return(l)
	return l
}

// Debug implements Logger with mock expectations
func (m *MockTestLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Debug(msg, args...)
}

// Info implements Logger with mock expectations
func (m *MockTestLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Info(msg, args...)
}

// Warn implements Logger with mock expectations
func (m *MockTestLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Warn(msg, args...)
}

// Error implements Logger with mock expectations
func (m *MockTestLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Error(msg, args...)
}

// WithError implements Logger with mock expectations
func (m *MockTestLogger) WithError(err error) Logger {
	m.Called(err)
	return m.Logger.WithError(err)
}
