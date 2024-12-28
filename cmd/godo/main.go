// Package main is the entry point for the Godo application
package main

import (
	"fmt"
	"os"

	"github.com/jonesrussell/godo/internal/container"
)

func main() {
	// Initialize the application with all dependencies
	app, cleanup, err := container.InitializeApp()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	// Setup UI components
	app.SetupUI()

	// Run the application
	app.Run()
}
