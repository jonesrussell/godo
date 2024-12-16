package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonesrussell/godo/internal/di"
	"github.com/jonesrussell/godo/internal/logger"
)

const shutdownTimeout = 30 * time.Second

// setupLogger configures and manages logger lifecycle
func setupLogger() func() {
	return func() {
		if err := logger.Sync(); err != nil {
			os.Stderr.WriteString("Failed to sync logger: " + err.Error() + "\n")
		}
	}
}

// setupSignalHandler creates a signal handler and returns a cancel function
func setupSignalHandler(ctx context.Context, cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal: %v", sig)
		cancel()
	}()
}

// initializeApplication handles app initialization
func initializeApplication() (*di.App, error) {
	logger.Info("Initializing dependency injection...")
	app, err := di.InitializeApp()
	if err != nil {
		return nil, err
	}
	logger.Info("Dependency injection initialized successfully")
	return app, nil
}

func main() {
	defer setupLogger()()

	logger.Info("Starting Godo application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := initializeApplication()
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}

	setupSignalHandler(ctx, cancel)

	logger.Info("Starting application run...")
	if err := app.Run(ctx); err != nil {
		logger.Fatal("Application error: %v", err)
	}

	logger.Info("Starting graceful shutdown...")
}
