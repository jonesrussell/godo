// Package app implements the main application logic for Godo.
package app

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// App implements the Application interface
type App struct {
	logger     logger.Logger
	fyneApp    fyne.App
	store      storage.Store
	mainWindow gui.MainWindow
	quickNote  gui.QuickNote
	hotkey     hotkey.Manager
	httpConfig *common.HTTPConfig
	name       common.AppName
	version    common.AppVersion
	id         common.AppID
}

// New creates a new App instance
func New(
	logger logger.Logger,
	fyneApp fyne.App,
	store storage.Store,
	mainWindow gui.MainWindow,
	quickNote gui.QuickNote,
	hotkey hotkey.Manager,
	httpConfig *common.HTTPConfig,
	name common.AppName,
	version common.AppVersion,
	id common.AppID,
) *App {
	return &App{
		logger:     logger,
		fyneApp:    fyneApp,
		store:      store,
		mainWindow: mainWindow,
		quickNote:  quickNote,
		hotkey:     hotkey,
		httpConfig: httpConfig,
		name:       name,
		version:    version,
		id:         id,
	}
}

// SetupUI initializes the UI components
func (a *App) SetupUI() {
	a.mainWindow.Show()
}

// Run starts the application
func (a *App) Run() {
	a.fyneApp.Run()
}

// Cleanup performs cleanup before application exit
func (a *App) Cleanup() {
	if err := a.store.Close(); err != nil {
		a.logger.Error("Failed to close store", "error", err)
	}
}
