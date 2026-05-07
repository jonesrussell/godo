package runtime

import (
	"context"
	"time"
)

// App is the runtime entrypoint contract.
type App interface {
	Run()
	Shutdown(ctx context.Context) error
	ForceKillTimeout() time.Duration
}

// ShutdownApp is the minimal contract required by CoordinatedShutdown.
type ShutdownApp interface {
	Shutdown(ctx context.Context) error
	ForceKillTimeout() time.Duration
}

const defaultForceKillTimeout = 3 * time.Second
