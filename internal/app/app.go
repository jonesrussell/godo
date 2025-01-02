// Package app implements the main application logic for Godo.
package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/systray"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/options"
	"github.com/jonesrussell/godo/internal/storage"
)

// App implements the Application interface
type App struct {
	name       common.AppName
	version    common.AppVersion
	id         common.AppID
	mainWindow gui.MainWindow
	quickNote  gui.QuickNote
	hotkey     hotkey.Manager
	logger     logger.Logger
	store      storage.TaskStore
	fyneApp    fyne.App
	config     *config.Config
}

// Params holds the parameters for creating a new App instance
type Params struct {
	Options *options.AppOptions
	Hotkey  hotkey.Manager
}

// New creates a new application instance using the options pattern
func New(params *Params) *App {
	return &App{
		name:       params.Options.Name,
		version:    params.Options.Version,
		id:         params.Options.ID,
		mainWindow: params.Options.GUI.MainWindow,
		quickNote:  params.Options.GUI.QuickNote,
		hotkey:     params.Hotkey,
		logger:     params.Options.Core.Logger,
		store:      params.Options.Core.Store,
		fyneApp:    params.Options.GUI.App,
		config:     params.Options.Core.Config,
	}
}

// setupHotkey initializes and starts the global hotkey system
func (a *App) setupHotkey() error {
	a.logger.Debug("Setting up global hotkey system",
		"config", a.config.Hotkeys.QuickNote)

	if err := a.hotkey.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	if err := a.hotkey.Start(); err != nil {
		return fmt.Errorf("failed to start hotkey listener: %w", err)
	}

	a.logger.Info("Hotkey system initialized successfully")
	return nil
}

// SetupUI initializes the user interface components in the correct order
func (a *App) SetupUI() error {
	a.logger.Debug("Setting up UI components")

	// 1. Set up systray first as it's the most visible component
	systray.SetupSystray(a.fyneApp, a.mainWindow.GetWindow(), a.quickNote)

	// 2. Show main window if not configured to start hidden
	if !a.config.UI.MainWindow.StartHidden {
		a.mainWindow.Show()
	}

	return nil
}

// Run starts the application
func (a *App) Run() {
	a.logger.Info("Starting application",
		"name", a.name,
		"version", a.version,
		"id", a.id,
	)

	// Set up UI components first
	if err := a.SetupUI(); err != nil {
		a.logger.Error("Failed to setup UI", "error", err)
		return
	}

	// Set up hotkey system
	if err := a.setupHotkey(); err != nil {
		a.logger.Error("Failed to setup hotkey system", "error", err)
		// Continue running even if hotkey fails
	}

	// Run the application main loop
	a.fyneApp.Run()
}

// Cleanup performs cleanup before application exit
func (a *App) Cleanup() {
	a.logger.Info("Cleaning up application")

	// First stop the hotkey manager
	if err := a.hotkey.Stop(); err != nil {
		a.logger.Error("Failed to stop hotkey manager", "error", err)
	} else {
		a.logger.Info("Hotkey manager stopped successfully")
	}

	// Then close the store
	if err := a.store.Close(); err != nil {
		a.logger.Error("Failed to close store", "error", err)
	} else {
		a.logger.Info("Store closed successfully")
	}

	a.logger.Info("Cleanup completed")
}

// Logger returns the application logger
func (a *App) Logger() logger.Logger {
	return a.logger
}

// Store returns the application store
func (a *App) Store() storage.TaskStore {
	return a.store
}
