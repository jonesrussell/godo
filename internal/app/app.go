// Package app implements the main application logic for Godo.
package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// App represents the main application
type App struct {
	config  *config.Config
	logger  logger.Logger
	store   storage.Store
	mainWin *mainwindow.Window
	hotkeys HotkeyManager
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, log logger.Logger, store storage.Store, mainWin *mainwindow.Window) *App {
	app := &App{
		config:  cfg,
		logger:  log,
		store:   store,
		mainWin: mainWin,
	}

	// Initialize hotkey manager based on build tags
	app.hotkeys = initHotkeyManager(app)

	return app
}

// Run starts the application
func (a *App) Run() error {
	// Setup hotkeys
	if err := a.hotkeys.Setup(); err != nil {
		return err
	}

	// Show main window
	if a.mainWin != nil {
		a.mainWin.Show()
	}

	return nil
}

// SetupUI initializes the application UI
func (a *App) SetupUI() {
	if a.mainWin != nil {
		a.mainWin.Setup()
	}
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return a.config.App.Version
}
