// Package app implements the main application logic for Godo.
package app

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/systray"
	"github.com/jonesrussell/godo/internal/logger"
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
	store      storage.Store
	fyneApp    fyne.App
	config     *common.HTTPConfig
}

// New creates a new application instance
func New(
	name common.AppName,
	version common.AppVersion,
	id common.AppID,
	mainWindow gui.MainWindow,
	quickNote gui.QuickNote,
	hotkey hotkey.Manager,
	logger logger.Logger,
	store storage.Store,
	fyneApp fyne.App,
	config *common.HTTPConfig,
) *App {
	return &App{
		name:       name,
		version:    version,
		id:         id,
		mainWindow: mainWindow,
		quickNote:  quickNote,
		hotkey:     hotkey,
		logger:     logger,
		store:      store,
		fyneApp:    fyneApp,
		config:     config,
	}
}

// SetupUI initializes the user interface
func (a *App) SetupUI() {
	// Set up main window
	a.mainWindow.Show()

	// Set up systray
	systray.SetupSystray(a.fyneApp, a.mainWindow.GetWindow())

	// Register hotkey
	if err := a.hotkey.Register(); err != nil {
		a.logger.Error("failed to register hotkey", "error", err)
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
func (a *App) Store() storage.Store {
	return a.store
}
