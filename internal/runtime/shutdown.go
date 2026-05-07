package runtime

import (
	"context"
	"errors"
	"sync"
)

// CoordinatedShutdown runs a single deterministic shutdown pass.
func CoordinatedShutdown(parent context.Context, app ShutdownApp, cleanup func()) error {
	if app == nil {
		return errors.New("runtime: nil shutdown app")
	}
	if parent == nil {
		parent = context.Background()
	}

	var once sync.Once
	runCleanup := func() {
		once.Do(func() {
			if cleanup != nil {
				cleanup()
			}
		})
	}
	defer runCleanup()

	timeout := app.ForceKillTimeout()
	if timeout <= 0 {
		timeout = defaultForceKillTimeout
	}

	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	err := app.Shutdown(ctx)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return ErrForcedShutdown
	}
	return err
}
