// Package app implements the main application logic for Godo.
package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// App represents the main application
type App struct {
	config    *config.Config
	logger    logger.Logger
	store     storage.Store
	mainWin   *mainwindow.Window
	quickNote quicknote.Interface
	hotkeys   HotkeyManager
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, log logger.Logger, store storage.Store, mainWin *mainwindow.Window) *App {
	app := &App{
		config:  cfg,
		logger:  log,
		store:   store,
		mainWin: mainWin,
	}

	// Create quick note window
	app.quickNote = quicknote.New(store)
	app.quickNote.Initialize(mainWin.GetApp(), log)

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

	// Run the main window loop
	if a.mainWin != nil {
		a.mainWin.Run()
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
