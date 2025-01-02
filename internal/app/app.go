// Package app implements the main application logic for Godo.
package app

import (
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

// SetupUI initializes the user interface
func (a *App) SetupUI() {
	// Set up systray first
	systray.SetupSystray(a.fyneApp, a.mainWindow.GetWindow(), a.quickNote)

	// Only show main window if not configured to start hidden
	if !a.config.UI.MainWindow.StartHidden {
		a.mainWindow.Show()
	}

	// Register hotkey
	if err := a.hotkey.Register(); err != nil {
		a.logger.Error("failed to register hotkey", "error", err)
	}

	// Start hotkey listener
	if err := a.hotkey.Start(); err != nil {
		a.logger.Error("failed to start hotkey listener", "error", err)
	}

	// Set up quick note handler
	if manager, ok := a.hotkey.(*hotkey.DefaultManager); ok {
		go func() {
			for range manager.GetHotkey().Keydown() {
				a.quickNote.Show()
			}
		}()
	}
}

// Run starts the application
func (a *App) Run() {
	a.logger.Info("starting application",
		"name", a.name,
		"version", a.version,
		"id", a.id,
	)

	// Run the application
	a.fyneApp.Run()
}

// Cleanup performs cleanup before application exit
func (a *App) Cleanup() {
	a.logger.Info("cleaning up application")

	if err := a.hotkey.Unregister(); err != nil {
		a.logger.Error("failed to unregister hotkey", "error", err)
	}

	if err := a.store.Close(); err != nil {
		a.logger.Error("failed to close store", "error", err)
	}
}

// Logger returns the application logger
func (a *App) Logger() logger.Logger {
	return a.logger
}

// Store returns the application store
func (a *App) Store() storage.TaskStore {
	return a.store
}
