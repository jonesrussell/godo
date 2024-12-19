package main

import (
	"fmt"
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
	"golang.design/x/hotkey"
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
		systrayIcon := assets.GetSystrayIconResource()
		appIcon := assets.GetAppIconResource()
		myApp.SetIcon(appIcon)

		// Create menu items
		quickNote := fyne.NewMenuItem("Quick Note (Ctrl+Alt+G)", nil)
		quickNote.Icon = systrayIcon
		quickNote.Shortcut = &desktop.CustomShortcut{
			KeyName:  fyne.KeyG,
			Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt,
		}

		preferences := fyne.NewMenuItem("Preferences", func() {
			logger.Debug("Opening preferences")
			// TODO: Add preferences dialog
		})

		separator := fyne.NewMenuItemSeparator()

		quit := fyne.NewMenuItem("Quit", func() {
			logger.Info("Application shutdown requested")
			myApp.Quit()
		})

		menu := fyne.NewMenu("Godo",
			quickNote,
			separator,
			preferences,
			separator,
			quit,
		)

		quickNote.Action = func() {
			logger.Debug("Opening quick note from tray")
			myApp.SendNotification(fyne.NewNotification("Quick Note", "Opening quick note..."))
			showQuickNote()
		}

		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(systrayIcon)
		logger.Info("System tray initialized")

		// Register quick note shortcut
		mainWindow := myApp.Driver().AllWindows()[0]
		mainWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
			KeyName:  fyne.KeyG,
			Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt,
		}, func(shortcut fyne.Shortcut) {
			logger.Debug("Quick Note shortcut triggered")
			showQuickNote()
		})
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
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		logger.Error("Failed to create app directory", "error", err)
		return "godo.db"
	}

	return filepath.Join(appDir, "godo.db")
}

func setupGlobalHotkey(callback func()) error {
	// Create a new hotkey combination: Ctrl+Alt+G
	hk := hotkey.New([]hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.ModAlt,
	}, hotkey.KeyG)

	// Register the hotkey
	if err := hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	// Handle hotkey events in a goroutine
	go func() {
		for {
			select {
			case <-hk.Keydown():
				logger.Debug("Global hotkey triggered")
				callback()
			}
		}
	}()

	return nil
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

	// Set up global hotkey
	if err := setupGlobalHotkey(qn.Show); err != nil {
		logger.Error("Failed to setup global hotkey", "error", err)
	}

	// Set up system tray with the quick note Show method
	setupSystemTray(myApp, qn.Show)

	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to the simplified Fyne app!"),
		widget.NewButton("Open Quick Note", qn.Show),
	))

	// Hide main window by default
	mainWindow.SetCloseIntercept(func() {
		logger.Debug("Window close intercepted, hiding instead")
		mainWindow.Hide()
	})

	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.CenterOnScreen()

	// Ensure window starts hidden
	mainWindow.Hide()

	// Run the application
	myApp.Run()
}
