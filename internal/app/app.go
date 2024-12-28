// Package app implements the main application logic for Godo.
package app

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/gui"
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
	Store     storage.Store
	mainWin   gui.MainWindow
	quickNote quicknote.Interface
	Hotkeys   hotkey.Manager
	version   string
}

// New creates a new application instance
func New(
	l logger.Logger,
	fyneApp fyne.App,
	store storage.Store,
	mainWin gui.MainWindow,
	quickNote quicknote.Interface,
	hotkeys hotkey.Manager,
) *App {
	return &App{
		Logger:    l,
		fyneApp:   fyneApp,
		Store:     store,
		mainWin:   mainWin,
		quickNote: quickNote,
		Hotkeys:   hotkeys,
		version:   "0.1.0",
	}
}

// Run starts the application
func (a *App) Run() error {
	a.Logger.Info("Starting application", "version", a.version)

	// Setup UI components first
	a.Logger.Debug("Setting up UI components")
	a.SetupUI()

	// Setup windows
	a.Logger.Debug("Setting up main window")
	if err := a.mainWin.Setup(); err != nil {
		a.Logger.Error("Failed to setup main window", "error", err)
		return err
	}
	a.Logger.Debug("Main window setup complete")

	a.Logger.Debug("Setting up quick note window")
	if err := a.quickNote.Setup(); err != nil {
		a.Logger.Error("Failed to setup quick note window", "error", err)
		return err
	}
	a.Logger.Debug("Quick note window setup complete")

	// Register global hotkeys last, after all windows are set up
	a.Logger.Debug("Registering global hotkeys")
	if err := a.Hotkeys.Register(); err != nil {
		a.Logger.Error("Failed to register hotkeys", "error", err)
		return err
	}
	a.Logger.Debug("Global hotkeys registered successfully")

	defer func() {
		a.Logger.Debug("Unregistering global hotkeys")
		if err := a.Hotkeys.Unregister(); err != nil {
			a.Logger.Error("Failed to unregister hotkeys", "error", err)
		}
	}()

	a.Logger.Info("Application initialized successfully, starting main loop")

	// Run the application
	a.fyneApp.Run()

	a.Logger.Info("Application main loop completed")
	return nil
}

// SetupUI initializes the application UI
func (a *App) SetupUI() {
	// Setup systray with menu
	a.Logger.Debug("Setting up system tray")
	menu := fyne.NewMenu("Godo",
		fyne.NewMenuItem("Quick Note", func() {
			a.Logger.Debug("Quick Note menu item clicked")
			a.quickNote.Show()
		}),
		fyne.NewMenuItem("Show", func() {
			a.Logger.Debug("Show menu item clicked")
			if win := a.mainWin.GetWindow(); win != nil {
				win.Show()
				a.Logger.Debug("Main window shown")
			} else {
				a.Logger.Error("Main window is nil")
			}
		}),
		fyne.NewMenuItem("Quit", func() {
			a.Logger.Debug("Quit menu item clicked")
			a.fyneApp.Quit()
		}),
	)

	tray := systray.New(a.fyneApp, a.Logger)
	tray.Setup(menu)
	tray.SetIcon(theme.GetSystrayIconResource())
	a.Logger.Debug("System tray setup complete")
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return a.version
}
