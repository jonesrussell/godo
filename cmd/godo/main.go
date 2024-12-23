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

	log, err := logger.New(defaultConfig)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Initialize app using dependency injection
	app, cleanup, err := container.InitializeApp()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer cleanup()

	// Setup and run the application
	app.SetupUI()
	app.Run()
}
