package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

// Runner manages the HTTP server lifecycle
type Runner struct {
	server   *Server
	logger   logger.Logger
	config   *config.HTTPConfig
	ready    chan struct{}
	shutdown chan struct{}
}

// NewRunner creates a new HTTP server runner
func NewRunner(
	taskService service.TaskService,
	l logger.Logger,
	config *config.HTTPConfig,
) *Runner {
	return &Runner{
		server:   NewServer(taskService, l),
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

		err := r.server.Start(port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			r.logger.Error("HTTP server failed to start",
				"port", port,
				"error", err.Error(),
				"hint", "Check if another process is using this port or configure a different port")
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
