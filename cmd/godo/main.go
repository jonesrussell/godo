package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonesrussell/godo/internal/di"
	"github.com/jonesrussell/godo/internal/logger"
)

func main() {
	logger.Info("Starting Godo application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize dependency injection container
	logger.Info("Initializing dependency injection...")
	app, err := di.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}
	logger.Info("Dependency injection initialized successfully")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal: %v", sig)
		cancel()
	}()

	logger.Info("Starting application run...")
	if err := app.Run(ctx); err != nil {
		logger.Fatal("Application error: %v", err)
	}
}
