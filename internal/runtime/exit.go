package runtime

import "errors"

const (
	ExitOK     = 0
	ExitError  = 1
	ExitPanic  = 2
	ExitForced = 3
)

// ErrForcedShutdown indicates graceful shutdown timed out.
var ErrForcedShutdown = errors.New("forced shutdown")

// NormalizeExit maps runtime errors to process exit codes.
func NormalizeExit(err error) int {
	if err == nil {
		return ExitOK
	}

	var panicErr *RecoveredPanicError
	if errors.As(err, &panicErr) || errors.Is(err, ErrRecoveredPanic) {
		return ExitPanic
	}
	if errors.Is(err, ErrForcedShutdown) {
		return ExitForced
	}
	return ExitError
}
