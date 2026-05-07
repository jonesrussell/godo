package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// NewRootContext creates a root context bound to SIGINT and SIGTERM.
func NewRootContext() (context.Context, context.CancelFunc) {
	return WithSignals(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

// WithSignals creates a context canceled by the provided OS signals.
func WithSignals(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return signal.NotifyContext(parent, signals...)
}
