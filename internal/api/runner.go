package api

import (
	"context"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Runner manages the HTTP server lifecycle
type Runner struct {
	server *Server
	logger logger.Logger
}

// NewRunner creates a new HTTP server runner
func NewRunner(store storage.Store, l logger.Logger) *Runner {
	return &Runner{
		server: NewServer(store, l),
		logger: l,
	}
}

// Start starts the HTTP server in a goroutine
func (r *Runner) Start(port int) {
	go func() {
		if err := r.server.Start(port); err != nil {
			r.logger.Error("HTTP server error", "error", err)
		}
	}()
}

// Shutdown gracefully shuts down the HTTP server
func (r *Runner) Shutdown(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}
