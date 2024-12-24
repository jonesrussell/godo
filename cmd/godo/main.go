package main

import (
	"fmt"
	"github.com/jonesrussell/godo/internal/app"
	"os"

	"github.com/jonesrussell/godo/internal/container"
	"go.uber.org/zap"
)

func run() error {
	logger, err := initializeLogger()
	if err != nil {
		return err
	}
	defer func() {
		_ = logger.Sync() // Ensure logger resources are flushed and released
	}()

	appInstance, cleanupApp, err := initializeAppWithLogger(logger)
	if err != nil {
		logger.Error("App initialization failed", zap.Error(err))
		return err
	}
	defer cleanupApp() // Ensure app resources are released on exit

	return handleAppRun(appInstance)
}

func initializeLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func initializeAppWithLogger(logger *zap.Logger) (*app.App, func(), error) {
	appInstance, cleanup, err := container.InitializeApp()
	if err != nil {
		logger.Error("Failed to initialize app", zap.Error(err))
		return nil, nil, err
	}
	appInstance.SetupUI()
	return appInstance, cleanup, nil
}

func handleAppRun(app *app.App) error {
	if err := app.Run(); err != nil {
		app.Logger.Error("Failed to run app", zap.Error(err))
		if unregisterErr := app.Hotkeys.Unregister(); unregisterErr != nil {
			app.Logger.Error("Failed to unregister hotkeys after app failure", zap.Error(unregisterErr))
		}
		return fmt.Errorf("app failed to run: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
