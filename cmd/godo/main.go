// Package main is the entry point for the Godo application
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonesrussell/godo/internal/container"
)

func main() {
	// Initialize the application with all dependencies
	app, cleanup, err := container.InitializeApp()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run cleanup on exit
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Cleaning up...")
		cleanup()
		os.Exit(0)
	}()

	// Setup UI components
	app.SetupUI()

	// Run the application
	app.Run()

	// If we get here normally (not through signal), still run cleanup
	cleanup()
}
