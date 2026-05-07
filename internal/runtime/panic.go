package runtime

import (
	"errors"
	"fmt"
)

// PanicLogger is the minimal logging interface used during panic recovery.
type PanicLogger interface {
	Error(msg string, keysAndValues ...any)
}

type noopPanicLogger struct{}

func (noopPanicLogger) Error(string, ...any) {}

// RecoveredPanicError wraps panic values converted to errors.
type RecoveredPanicError struct {
	Value any
}

func (e *RecoveredPanicError) Error() string {
	return fmt.Sprintf("panic recovered: %v", e.Value)
}

// ErrRecoveredPanic marks errors sourced from panic recovery.
var ErrRecoveredPanic = errors.New("recovered panic")

func (e *RecoveredPanicError) Unwrap() error {
	return ErrRecoveredPanic
}

// RecoverFromPanic converts recover() values into structured errors.
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

// RecoverAndReport is an alias maintained for readability in call sites.
func RecoverAndReport(log PanicLogger, recovered any) error {
	return RecoverFromPanic(log, recovered)
}

// WithPanicRecovery executes fn and converts panics into returned errors.
func WithPanicRecovery(log PanicLogger, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = RecoverFromPanic(log, r)
		}
	}()
	return fn()
}
