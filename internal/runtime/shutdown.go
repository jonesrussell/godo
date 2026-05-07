package runtime

import (
	"context"
	"errors"
	"sync"
	"time"
)

// CoordinatedShutdown runs a single ordered shutdown round against app.
//
// Order (skeleton; future WPs may extend app.Shutdown to sequence UI, hotkeys,
// storage, audit, harness without adding orchestration logic here):
//  1. Apply force-kill deadline from app.ForceKillTimeout() (or
//     DefaultForceKillTimeout when zero).
//  2. Invoke app.Shutdown with a context cancelled at that deadline.
//  3. Map deadline exceeded to ErrForcedShutdown for NormalizeExit (ExitForced).
//
// cleanup runs exactly once on return (success, error, or panic from app;
// panics are not recovered here). No stdout writes and no os.Exit.
func CoordinatedShutdown(parent context.Context, app ShutdownApp, cleanup func()) error {
	if app == nil {
		return errors.New("runtime: nil ShutdownApp")
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

	d := app.ForceKillTimeout()
	if d <= 0 {
		d = effectiveForceKillTimeout()
	}

	ctx, cancel := context.WithTimeout(parent, d)
	defer cancel()

	err := app.Shutdown(ctx)

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrForcedShutdown
	}
	if err != nil {
		return err
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return ErrForcedShutdown
	}

	return nil
}

// defaultForceKillTimeout is the active default when ForceKillTimeout() is zero.
// Tests may override it to avoid long sleeps; production stays at
// DefaultForceKillTimeout.
var defaultForceKillTimeout = DefaultForceKillTimeout

func effectiveForceKillTimeout() time.Duration {
	return defaultForceKillTimeout
}
