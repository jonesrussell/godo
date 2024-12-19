package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		logger.Info("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		logger.Info("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		logger.Info("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		logger.Info("Lifecycle: Exited Foreground")
	})
}

func setupSystemTray(myApp fyne.App, showQuickNote func()) {
	if desk, ok := myApp.(desktop.App); ok {
		// Use different icons for systray and app
		systrayIcon := assets.GetSystrayIconResource()
		appIcon := assets.GetAppIconResource()

		// Set the application icon
		myApp.SetIcon(appIcon)

		quickNote := fyne.NewMenuItem("Quick Note", nil)
		quickNote.Icon = systrayIcon
		menu := fyne.NewMenu("Godo", quickNote)

		quickNote.Action = func() {
			logger.Debug("Opening quick note from tray")
			myApp.SendNotification(fyne.NewNotification("Quick Note", "Opening quick note..."))
			showQuickNote()
		}

		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(systrayIcon)
		logger.Info("System tray initialized")
	} else {
		logger.Warn("System tray not supported on this platform")
	}
}

func getDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Failed to get home directory", "error", err)
		return "godo.db"
	}

	appDir := filepath.Join(homeDir, ".godo")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		logger.Error("Failed to create app directory", "error", err)
		return "godo.db"
	}

	return filepath.Join(appDir, "godo.db")
}

func main() {
	// Initialize logger
	_, err := logger.Initialize()
	if err != nil {
		panic(err)
	}

	// Initialize SQLite storage
	dbPath := getDBPath()
	store, err := sqlite.New(dbPath)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		panic(err)
	}
	defer store.Close()

	// Create app with unique ID
	myApp := app.NewWithID("io.github.jonesrussell.godo")
	logLifecycle(myApp)

	mainWindow := myApp.NewWindow("Godo")

	// Create quick note instance with store
	qn := quicknote.New(mainWindow, store)

	// Set up system tray with the quick note Show method
	setupSystemTray(myApp, qn.Show)

	// Hide main window by default
	mainWindow.SetCloseIntercept(func() {
		logger.Debug("Window close intercepted, hiding instead")
		mainWindow.Hide()
	})

	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to the simplified Fyne app!"),
		widget.NewButton("Open Quick Note", qn.Show),
	))
	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.CenterOnScreen()

	// Start hidden
	mainWindow.Hide()
	myApp.Run()
}
