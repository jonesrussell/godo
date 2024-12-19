package main

import (
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/container"
	"github.com/jonesrussell/godo/internal/logger"
)

func main() {
	// Initialize logger with default config
	defaultConfig := &common.LogConfig{
		Level:       "info",
		Output:      []string{"stdout", "godo.log"},
		ErrorOutput: []string{"stderr", "godo.error.log"},
	}

	if _, err := logger.Initialize(defaultConfig); err != nil {
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
