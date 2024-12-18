package main

import (
	"errors"
	"os"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/marcsauter/single"
)

func main() {
	if err := run(); err != nil {
		logger.Error("Application error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize logger
	if _, err := logger.Initialize(); err != nil {
		return err
	}

	// Single instance check
	if err := checkSingleInstance(); err != nil {
		return err
	}

	// Initialize app
	application, err := initializeApp()
	if err != nil {
		logger.Error("Failed to initialize app", "error", err)
		return err
	}
	defer cleanup(application)

	// Create and run GUI
	uiApp := gui.New(application)
	return uiApp.Run()
}

func checkSingleInstance() error {
	// Create single instance lock
	s := single.New("godo")
	if err := s.CheckLock(); err != nil {
		if err == single.ErrAlreadyRunning {
			logger.Info("Godo is already running. Look for the icon in your system tray.")
			return nil
		}
		return errors.New("failed to check application lock: " + err.Error())
	}
	defer func() {
		if err := s.TryUnlock(); err != nil {
			logger.Error("Failed to unlock single instance", "error", err)
		}
	}()
	return nil
}

func initializeApp() (*app.App, error) {
	cfg, err := initializeConfig()
	if err != nil {
		return nil, errors.New("failed to initialize config: " + err.Error())
	}

	application, err := app.InitializeAppWithConfig(cfg)
	if err != nil {
		return nil, errors.New("failed to initialize app: " + err.Error())
	}
	return application, nil
}

func cleanup(application *app.App) {
	logger.Info("Cleaning up application...")
	if err := application.Cleanup(); err != nil {
		logger.Error("Failed to cleanup", "error", err)
	}
}

func initializeConfig() (*config.Config, error) {
	env := os.Getenv("GODO_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		return nil, err
	}

	logConfig := common.LogConfig{
		Level:       cfg.Logging.Level,
		Output:      cfg.Logging.Output,
		ErrorOutput: cfg.Logging.ErrorOutput,
	}

	if _, err := logger.InitializeWithConfig(logConfig); err != nil {
		return nil, err
	}

	return cfg, nil
}
