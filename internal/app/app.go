// Package app implements the main application logic for Godo.
package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// App represents the main application
type App struct {
	config    *config.Config
	logger    logger.Logger
	store     storage.Store
	mainWin   MainWindow
	quickNote QuickNote
	hotkeys   HotkeyManager
	isDocker  bool // Replaces global buildTagDocker
}

// MainWindow defines the interface for the main application window
type MainWindow interface {
	Show()
	Setup()
}

// QuickNote defines the interface for the quick note window
type QuickNote interface {
	Show()
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, log logger.Logger, store storage.Store) *App {
	app := &App{
		config:   cfg,
		logger:   log,
		store:    store,
		isDocker: false, // Will be set to true by build_docker.go init()
	}

	// Initialize hotkey manager based on build tags
	if app.isDocker {
		app.hotkeys = NewNoopHotkeyManager(app)
	} else {
		app.hotkeys = NewDefaultHotkeyManager(app)
	}

	return app
}

// Run starts the application
func (a *App) Run() error {
	// Setup hotkeys
	if err := a.hotkeys.Setup(); err != nil {
		return err
	}

	// Show main window
	a.mainWin.Show()

	return nil
}

// SetupUI initializes the application UI
func (a *App) SetupUI() {
	a.mainWin.Setup()
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return a.config.App.Version
}

// SetIsDocker sets the Docker environment flag
func (a *App) SetIsDocker(isDocker bool) {
	a.isDocker = isDocker
}
