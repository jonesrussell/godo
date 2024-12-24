package main

import (
	"os"

	"github.com/jonesrussell/godo/internal/container"
	"go.uber.org/zap"
)

func run() error {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync() // Ignore sync errors
	}()

	// Create app
	app, cleanup, err := container.InitializeApp()
	if err != nil {
		logger.Error("Failed to initialize app", zap.Error(err))
		cleanup()
		return err
	}
	defer cleanup()

	// Setup UI
	app.SetupUI()

	// Run app
	if err := app.Run(); err != nil {
		logger.Error("Failed to run app", zap.Error(err))
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
