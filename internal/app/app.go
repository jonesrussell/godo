// Package app implements the main application logic for Godo.
package app

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/systray"
	"github.com/jonesrussell/godo/internal/logger"
)

// App implements the main application
type App struct {
	name       string
	version    string
	mainWindow gui.MainWindow
	quickNote  gui.QuickNote
	hotkey     hotkey.Manager
	logger     logger.Logger
	fyneApp    fyne.App
}

// New creates a new application instance
func New(
	name string,
	version string,
	mainWindow gui.MainWindow,
	quickNote gui.QuickNote,
	hotkey hotkey.Manager,
	logger logger.Logger,
	fyneApp fyne.App,
) *App {
	return &App{
		name:       name,
		version:    version,
		mainWindow: mainWindow,
		quickNote:  quickNote,
		hotkey:     hotkey,
		logger:     logger,
		fyneApp:    fyneApp,
	}
}

// Start starts the application
func (a *App) Start() error {
	a.logger.Info("starting application",
		"name", a.name,
		"version", a.version,
	)

	// Set up hotkey
	if err := a.hotkey.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	// Set up UI
	a.SetupUI()

	// Run the application
	a.fyneApp.Run()
	return nil
}

// SetupUI initializes the user interface
func (a *App) SetupUI() {
	// Set up main window
	a.mainWindow.Show()

	// Set up systray
	systray.SetupSystray(a.fyneApp, a.mainWindow.GetWindow())
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("stopping application")

	if err := a.hotkey.Unregister(); err != nil {
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}

	a.fyneApp.Quit()
	return nil
}
