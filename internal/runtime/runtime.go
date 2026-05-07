package runtime

import (
	"context"
	"errors"
)

// Run is the single runtime entrypoint for process lifecycle orchestration.
func Run(rootCtx context.Context, app App, cleanup func(), log PanicLogger) int {
	if app == nil {
		return NormalizeExit(errors.New("runtime: nil app"))
	}
	if rootCtx == nil {
		rootCtx = context.Background()
	}

	runDone := make(chan error, 1)
	go func() {
		runDone <- WithPanicRecovery(log, func() error {
			app.Run()
			return nil
		})
	}()

	select {
	case runErr := <-runDone:
		shutdownErr := CoordinatedShutdown(context.Background(), app, cleanup)
		if runErr != nil {
			return NormalizeExit(runErr)
		}
		return NormalizeExit(shutdownErr)
	case <-rootCtx.Done():
		shutdownErr := CoordinatedShutdown(context.Background(), app, cleanup)
		select {
		case runErr := <-runDone:
			if runErr != nil {
				return NormalizeExit(runErr)
			}
		default:
		}
		return NormalizeExit(shutdownErr)
	}
}
