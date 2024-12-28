package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonesrussell/godo/internal/api"
	godocontainer "github.com/jonesrussell/godo/internal/container"
	"go.uber.org/zap"
)

func run() error {
	fmt.Println("=== Starting Godo... ===")

	// Initialize logger first
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return err
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Failed to sync logger: %v\n", err)
		}
	}()

	// Create application container
	container, err := godocontainer.Initialize(logger.Sugar())
	if err != nil {
		fmt.Printf("Failed to initialize container: %v\n", err)
		logger.Error("Failed to initialize container", zap.Error(err))
		return err
	}

	// Create HTTP server runner
	httpRunner := api.NewRunner(container.Store, container.Logger)
	httpRunner.Start(8080)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run app in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := container.App.Run(); err != nil {
			errChan <- err
		}
	}()

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		logger.Info("Received signal", zap.String("signal", sig.String()))
		// Shutdown HTTP server
		ctx := context.Background()
		if err := httpRunner.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown HTTP server", zap.Error(err))
		}
	case err := <-errChan:
		logger.Error("Application error", zap.Error(err))
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
