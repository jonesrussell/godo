// Package main is the entry point for the Godo application
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonesrussell/godo/internal/application/container"
	"github.com/jonesrussell/godo/internal/application/core"
)

func main() {
	// Initialize the application with all dependencies
	myapp, cleanup, err := container.InitializeApp()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Ensure cleanup runs in all exit scenarios
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Application panicked: %v\n", r)
			cleanup()
			os.Exit(1)
		}
	}()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create done channel to coordinate shutdown
	done := make(chan bool, 1)

	// Run signal handler in a separate goroutine
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Cleaning up...")

		// Run cleanup BEFORE terminating the application
		cleanup()

		// Get the concrete App type and quit
		if godoApp, ok := myapp.(*core.App); ok {
			godoApp.Quit()
		}

		// Signal that we're done
		done <- true

		// Force kill after a delay if we're still running
		go func() {
			forceKillTimeout := myapp.ForceKillTimeout()
			if forceKillTimeout == 0 {
				forceKillTimeout = 3 * time.Second
			}
			time.Sleep(forceKillTimeout)
			fmt.Println("Forcing process termination...")
			os.Exit(1)
		}()
	}()

	// Run the application from the main goroutine (required by Fyne)
	// The signal handler will call myapp.Quit() when a signal is received
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Application error: application panicked: %v\n", r)
			cleanup()
			os.Exit(1)
		}
	}()

	myapp.Run()

	// Normal app completion
	cleanup()
}
