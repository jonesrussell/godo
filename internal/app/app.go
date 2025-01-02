// Package app implements the main application logic for Godo.
package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/systray"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/options"
	"github.com/jonesrussell/godo/internal/storage"
)

// Application defines the interface for the main application
type Application interface {
	SetupUI() error
	Run()
	Cleanup()
	Logger() logger.Logger
	Store() storage.Store
}

// App implements the Application interface
type App struct {
	logger     logger.Logger
	config     *config.Config
	mainWindow gui.MainWindowManager
	quickNote  gui.QuickNoteManager
	hotkey     hotkey.Manager
	store      storage.Store
	fyneApp    fyne.App
}

// Params holds the parameters for creating a new App instance
type Params struct {
	Options   *options.AppOptions
	Hotkey    hotkey.Manager
	Version   string
	Commit    string
	BuildTime string
}

// New creates a new application instance using the options pattern
func New(params *Params) (*App, error) {
	app := &App{
		logger:     params.Options.Core.Logger,
		config:     params.Options.Core.Config,
		mainWindow: params.Options.GUI.MainWindow,
		quickNote:  params.Options.GUI.QuickNote,
		hotkey:     params.Hotkey,
		store:      params.Options.Core.Store,
		fyneApp:    params.Options.GUI.App,
	}

	app.logger.Info("Application initialized",
		"version", params.Version,
		"commit", params.Commit,
		"build_time", params.BuildTime,
		"config", app.config.Hotkeys.QuickNote)

	return app, nil
}

// SetupUI initializes the user interface components in the correct order
func (a *App) SetupUI() error {
	a.logger.Debug("Setting up UI components")

	// Set up system tray
	systray.SetupSystray(a.fyneApp, a.mainWindow, a.quickNote)

	// Show main window if not configured to start hidden
	if !a.config.UI.MainWindow.StartHidden {
		a.mainWindow.Show()
	}

	return nil
}

// Run starts the application
func (a *App) Run() {
	a.logger.Info("Starting application",
		"hotkey_active", a.hotkey != nil,
		"store_type", fmt.Sprintf("%T", a.store))

	// Set up UI components first
	if err := a.SetupUI(); err != nil {
		a.logger.Error("Failed to setup UI", "error", err)
		return
	}

	// Start hotkey manager if available
	if a.hotkey != nil {
		if err := a.hotkey.Start(); err != nil {
			a.logger.Error("Failed to start hotkey manager", "error", err)
			return
		}
		defer func() {
			if err := a.hotkey.Stop(); err != nil {
				a.logger.Error("Failed to stop hotkey manager", "error", err)
			}
		}()
	}

	// Run the application main loop
	a.fyneApp.Run()
}

// Cleanup performs cleanup before application exit
func (a *App) Cleanup() {
	a.logger.Info("Cleaning up application")

	// Stop hotkey manager if running
	if a.hotkey != nil {
		if err := a.hotkey.Stop(); err != nil {
			a.logger.Error("Failed to stop hotkey manager", "error", err)
		}
	}

	// Finally close the store
	if err := a.store.Close(); err != nil {
		a.logger.Error("Failed to close store", "error", err)
	}
	a.logger.Info("Store closed successfully")
}

// Logger returns the application logger
func (a *App) Logger() logger.Logger {
	return a.logger
}

// Store returns the task store
func (a *App) Store() storage.Store {
	return a.store
}
