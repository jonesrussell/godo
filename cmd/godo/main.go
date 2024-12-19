package main

import (
	"github.com/jonesrussell/godo/internal/container"
	"github.com/jonesrussell/godo/internal/logger"
)

func main() {
	// Initialize logger
	if _, err := logger.Initialize(); err != nil {
		panic(err)
	}

	// Initialize app using dependency injection
	app, cleanup, err := container.InitializeApp()
	if err != nil {
		logger.Error("Failed to initialize application", "error", err)
		panic(err)
	}
	defer cleanup()

	// Setup and run the application
	app.SetupUI()
	app.Run()
}
