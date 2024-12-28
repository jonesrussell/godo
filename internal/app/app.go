// Package app implements the main application logic for Godo.
package app

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/mainwindow/systray"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/gui/theme"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// App represents the main application
type App struct {
	Logger    logger.Logger
	fyneApp   fyne.App
	store     storage.Store
	mainWin   *mainwindow.Window
	quickNote quicknote.Interface
	Hotkeys   HotkeyManager
	version   string
}

// New creates a new application instance
func New(log logger.Logger, fyneApp fyne.App, store storage.Store, hotkeys HotkeyManager) *App {
	return &App{
		Logger:  log,
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
		a.Logger.Error("Failed to setup main window", "error", err)
		return err
	}

	if err := a.quickNote.Setup(); err != nil {
		a.Logger.Error("Failed to setup quick note window", "error", err)
		return err
	}

	// Register global hotkeys
	if err := a.Hotkeys.Register(); err != nil {
		a.Logger.Error("Failed to register hotkeys", "error", err)
		return err
	}
	defer func() {
		if err := a.Hotkeys.Unregister(); err != nil {
			a.Logger.Error("Failed to unregister hotkeys", "error", err)
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

	// Setup systray with menu
	menu := fyne.NewMenu("Godo",
		fyne.NewMenuItem("Show", func() {
			if win := a.mainWin.GetWindow(); win != nil {
				win.Show()
			}
		}),
		fyne.NewMenuItem("Quit", func() {
			a.fyneApp.Quit()
		}),
	)

	systray := systray.New(a.fyneApp, a.Logger)
	systray.Setup(menu)
	systray.SetIcon(theme.AppIcon())

	// Create quick note window
	a.quickNote = quicknote.New(a.store, a.Logger)
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return a.version
}
