package api

import (
	"context"
	"net/http"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Runner manages the HTTP server lifecycle
type Runner struct {
	server   *Server
	logger   logger.Logger
	config   *common.HTTPConfig
	ready    chan struct{}
	shutdown chan struct{}
}

// NewRunner creates a new HTTP server runner
func NewRunner(store storage.TaskStore, l logger.Logger, config *common.HTTPConfig) *Runner {
	return &Runner{
		server:   NewServer(store, l),
		logger:   l,
		config:   config,
		ready:    make(chan struct{}),
		shutdown: make(chan struct{}),
	}
}

// Start starts the HTTP server in a goroutine with proper synchronization
func (r *Runner) Start(port int) {
	go func() {
		defer close(r.shutdown)

		// Signal that we're attempting to start
		r.logger.Info("Starting HTTP server", "port", port)

		if err := r.server.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			r.logger.Error("HTTP server error", "error", err)
			return
		}

		r.logger.Info("HTTP server stopped")
	}()

	// Wait a brief moment for server to initialize
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(r.ready)
	}()
}

// WaitForReady waits for the server to be ready to accept connections
func (r *Runner) WaitForReady(timeout time.Duration) bool {
	select {
	case <-r.ready:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Shutdown gracefully shuts down the HTTP server
func (r *Runner) Shutdown(ctx context.Context) error {
	err := r.server.Shutdown(ctx)

	// Wait for server goroutine to complete
	select {
	case <-r.shutdown:
		r.logger.Info("Server shutdown completed")
	case <-ctx.Done():
		r.logger.Warn("Server shutdown timeout")
	}

	return err
}
