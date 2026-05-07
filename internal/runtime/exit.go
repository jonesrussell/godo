package runtime

import (
	"errors"
)

// Process exit codes for NormalizeExit. Callers in main should map these to
// os.Exit; this package must not call os.Exit.
const (
	ExitOK     = 0
	ExitError  = 1
	ExitPanic  = 2
	ExitForced = 3
)

// ErrForcedShutdown indicates the process should exit with ExitForced.
// Reserved for future force-kill / timeout paths.
var ErrForcedShutdown = errors.New("forced shutdown")

// NormalizeExit maps runtime errors to process exit codes.
func NormalizeExit(err error) int {
	if err == nil {
		return ExitOK
	}
	var pe *RecoveredPanicError
	if errors.As(err, &pe) {
		return ExitPanic
	}
	if errors.Is(err, ErrRecoveredPanic) {
		return ExitPanic
	}
	if errors.Is(err, ErrForcedShutdown) {
		return ExitForced
	}
	return ExitError
}
