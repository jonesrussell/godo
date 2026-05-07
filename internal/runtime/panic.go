package runtime

import (
	"errors"
	"fmt"
)

// PanicLogger is the minimal surface needed to report recovered panics without
// printing to stdout/stderr.
type PanicLogger interface {
	Error(msg string, keysAndValues ...any)
}

type noopPanicLogger struct{}

func (noopPanicLogger) Error(string, ...any) {}

// RecoveredPanicError wraps a value recovered from panic.
type RecoveredPanicError struct {
	Value any
}

func (e *RecoveredPanicError) Error() string {
	return fmt.Sprintf("panic recovered: %v", e.Value)
}

// ErrRecoveredPanic is a stable sentinel for panic-derived failures.
var ErrRecoveredPanic = errors.New("recovered panic")

func (e *RecoveredPanicError) Unwrap() error {
	return ErrRecoveredPanic
}

// RecoverFromPanic turns a non-nil recover() value into an error and logs it.
// If recovered is nil, returns nil. If log is nil, a no-op logger is used.
func RecoverFromPanic(log PanicLogger, recovered any) error {
	if recovered == nil {
		return nil
	}
	if log == nil {
		log = noopPanicLogger{}
	}
	log.Error("panic recovered", "panic", recovered)
	return &RecoveredPanicError{Value: recovered}
}

// RecoverAndReport is equivalent to RecoverFromPanic.
func RecoverAndReport(log PanicLogger, recovered any) error {
	return RecoverFromPanic(log, recovered)
}

// WithPanicRecovery runs fn and converts any panic into a RecoveredPanicError.
// The returned error is nil if fn completes without panicking.
func WithPanicRecovery(log PanicLogger, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = RecoverFromPanic(log, r)
		}
	}()
	return fn()
}
