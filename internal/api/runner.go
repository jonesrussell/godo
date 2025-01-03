package api

import (
	"context"
	"fmt"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Runner manages the HTTP server lifecycle
type Runner struct {
	server *Server
	logger logger.Logger
	config *common.HTTPConfig
}

// NewRunner creates a new HTTP server runner
func NewRunner(store storage.Store, l logger.Logger, config *common.HTTPConfig) *Runner {
	// Convert logger.Logger to *zap.Logger
	zapLogger, ok := l.(*logger.ZapLogger)
	if !ok {
		panic("logger must be a ZapLogger")
	}

	return &Runner{
		server: NewServer(store, zapLogger.SugaredLogger.Desugar()),
		logger: l,
		config: config,
	}
}

// Start starts the HTTP server in a goroutine
func (r *Runner) Start(port int) {
	go func() {
		addr := fmt.Sprintf(":%d", port)
		if err := r.server.Start(addr); err != nil {
			r.logger.Error("HTTP server error", "error", err)
		}
	}()
}

// Shutdown gracefully shuts down the HTTP server
func (r *Runner) Shutdown(ctx context.Context) error {
	return r.server.Stop(ctx)
}
