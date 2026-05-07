package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// NewRootContext creates the root runtime context bound to process interrupt
// signals. The returned cancel function must always be called by the caller to
// stop signal delivery and release internal resources.
func NewRootContext() (context.Context, context.CancelFunc) {
	return WithSignals(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

// WithSignals creates a context cancelled by either parent cancellation or any
// of the provided OS signals. This WP03 helper only wires signal/context
// behavior; shutdown orchestration is intentionally deferred.
func WithSignals(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return signal.NotifyContext(parent, signals...)
}
