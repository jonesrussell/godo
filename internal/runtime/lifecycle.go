package runtime

import (
	"context"
	"time"
)

// Lifecycle is the minimal contract runtime orchestration expects from the app.
//
// TODO: evolve this interface if additional lifecycle hooks are required.
type Lifecycle interface {
	// Start begins application execution and blocks until completion or shutdown.
	Start(ctx context.Context) error

	// Shutdown requests graceful teardown for all app-managed resources.
	Shutdown(ctx context.Context) error

	// ForceKillTimeout defines maximum wait before forced process termination.
	ForceKillTimeout() time.Duration
}
