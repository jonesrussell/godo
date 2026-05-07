package runtime

import (
	"context"
	"time"
)

// ShutdownApp is the minimal contract for coordinated shutdown. Adapters will
// wrap application types (e.g. core.App) without embedding UI/hotkey logic in
// the runtime package.
type ShutdownApp interface {
	Shutdown(ctx context.Context) error
	ForceKillTimeout() time.Duration
}

// DefaultForceKillTimeout is used when ForceKillTimeout() returns zero.
const DefaultForceKillTimeout = 3 * time.Second
