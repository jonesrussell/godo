// Package app implements the main application logic for Godo.
package app

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/storage"
	"go.uber.org/zap"
)

// App represents the main application
type App struct {
	Logger    *zap.Logger
	fyneApp   fyne.App
	store     storage.Store
	mainWin   *mainwindow.Window
	quickNote quicknote.Interface
	Hotkeys   HotkeyManager
	version   string
}

// New creates a new application instance
func New(logger *zap.Logger, fyneApp fyne.App, store storage.Store, hotkeys HotkeyManager) *App {
	return &App{
		Logger:  logger,
		fyneApp: fyneApp,
		store:   store,
		Hotkeys: hotkeys,
		version: "0.1.0",
	}
}

// Run starts the application
func (a *App) Run() error {
	// Setup windows
	if err := a.mainWin.Setup(); err != nil {
		a.Logger.Error("Failed to setup main window", zap.Error(err))
		return err
	}

	if err := a.quickNote.Setup(); err != nil {
		a.Logger.Error("Failed to setup quick note window", zap.Error(err))
		return err
	}

	// Register global hotkeys
	if err := a.Hotkeys.Register(); err != nil {
		a.Logger.Error("Failed to register hotkeys", zap.Error(err))
		return err
	}
	defer func() {
		if err := a.Hotkeys.Unregister(); err != nil {
			a.Logger.Error("Failed to unregister hotkeys", zap.Error(err))
		}
	}()

	// Run the application
	a.fyneApp.Run()

	return nil
}

// SetupUI initializes the application UI
func (a *App) SetupUI() {
	// Create main window
	a.mainWin = mainwindow.New(a.store, a.Logger)

	// Create quick note window
	a.quickNote = quicknote.New(a.store, a.Logger)
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return a.version
}
