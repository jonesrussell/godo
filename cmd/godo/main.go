// Package main is the entry point for the Godo application
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonesrussell/godo/internal/api"
	"github.com/jonesrussell/godo/internal/common"
	godocontainer "github.com/jonesrussell/godo/internal/container"
	"github.com/jonesrussell/godo/internal/logger"
)

const (
	// Default HTTP server port
	defaultHTTPPort = 8080
)

func run() error {
	fmt.Println("=== Starting Godo... ===")

	// Initialize logger first
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}
	log, err := logger.New(logConfig)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return err
	}

	// Create application container
	container, err := godocontainer.Initialize(log)
	if err != nil {
		fmt.Printf("Failed to initialize container: %v\n", err)
		log.Error("Failed to initialize container", "error", err)
		return err
	}

	// Create HTTP server runner
	httpRunner := api.NewRunner(container.Store, container.Logger)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		httpRunner.Start(defaultHTTPPort)
	}()

	// Handle shutdown in a separate goroutine
	go func() {
		sig := <-sigChan
		log.Info("Received signal", "signal", sig.String())
		// Shutdown HTTP server
		ctx := context.Background()
		if err := httpRunner.Shutdown(ctx); err != nil {
			log.Error("Failed to shutdown HTTP server", "error", err)
		}
		// Signal the app to quit
		container.App.Quit()
	}()

	// Run the Fyne app in the main goroutine
	return container.App.Run()
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
