// Package main is the entry point for the Godo application
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonesrussell/godo/internal/application"
	"github.com/jonesrussell/godo/internal/application/container"
)

func main() {
	// Initialize the application with all dependencies
	myapp, cleanup, err := container.InitializeApp()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run signal handler in a separate goroutine
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Cleaning up...")

		// Get the concrete App type
		if godoApp, ok := myapp.(*application.App); ok {
			// Quit the application
			godoApp.Quit()
		}

		// Force kill after a delay if we're still running
		go func() {
			forceKillTimeout := myapp.ForceKillTimeout()
			if forceKillTimeout == 0 {
				forceKillTimeout = 2 * time.Second
			}
			time.Sleep(forceKillTimeout)
			fmt.Println("Forcing process termination...")
			os.Exit(1)
		}()
	}()

	// Run the application
	myapp.Run()

	// Run cleanup on normal exit
	cleanup()
}
