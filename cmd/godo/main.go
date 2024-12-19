package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
	"golang.design/x/hotkey"
)

type App struct {
	fyneApp    fyne.App
	mainWindow fyne.Window
	quickNote  *quicknote.QuickNote
	store      storage.Store
}

func NewApp() (*App, error) {
	// Initialize SQLite storage
	store, err := initializeStorage()
	if err != nil {
		return nil, err
	}

	fyneApp := fyneapp.NewWithID("io.github.jonesrussell.godo")
	mainWindow := fyneApp.NewWindow("Godo")

	app := &App{
		fyneApp:    fyneApp,
		mainWindow: mainWindow,
		store:      store,
	}

	app.quickNote = quicknote.New(mainWindow, store)

	return app, nil
}

func (a *App) setupUI() {
	a.setupLifecycleLogging()
	a.setupSystemTray()
	a.setupMainWindow()
	if err := a.setupGlobalHotkey(); err != nil {
		logger.Error("Failed to setup global hotkey", "error", err)
	}
}

func (a *App) setupLifecycleLogging() {
	a.fyneApp.Lifecycle().SetOnStarted(func() {
		logger.Info("Lifecycle: Started")
	})
	a.fyneApp.Lifecycle().SetOnStopped(func() {
		logger.Info("Lifecycle: Stopped")
	})
}

func (a *App) setupSystemTray() {
	if desk, ok := a.fyneApp.(desktop.App); ok {
		logger.Debug("Loading system tray icon")
		systrayIcon := assets.GetSystrayIconResource()
		appIcon := assets.GetAppIconResource()
		a.fyneApp.SetIcon(appIcon)

		menu := fyne.NewMenu("Godo",
			fyne.NewMenuItem("Quick Note", a.quickNote.Show),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				a.fyneApp.Quit()
			}),
		)

		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(systrayIcon)
		logger.Info("System tray initialized")
	} else {
		logger.Warn("System tray not supported on this platform")
	}
}

func (a *App) setupMainWindow() {
	a.mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to Godo!"),
		widget.NewButton("Open Quick Note", a.quickNote.Show),
	))

	a.mainWindow.SetCloseIntercept(func() {
		logger.Debug("Window close intercepted, hiding instead")
		a.mainWindow.Hide()
	})

	a.mainWindow.Resize(fyne.NewSize(800, 600))
	a.mainWindow.CenterOnScreen()
	a.mainWindow.Hide()
}

func (a *App) setupGlobalHotkey() error {
	return setupGlobalHotkey(a.quickNote.Show)
}

func (a *App) Run() {
	a.fyneApp.Run()
}

func (a *App) Cleanup() {
	if db, ok := a.store.(*sqlite.Store); ok {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database", "error", err)
		}
	}
}

func main() {
	// Initialize logger
	if _, err := logger.Initialize(); err != nil {
		panic(err)
	}

	app, err := NewApp()
	if err != nil {
		logger.Error("Failed to create application", "error", err)
		panic(err)
	}
	defer app.Cleanup()

	app.setupUI()
	app.Run()
}

func initializeStorage() (storage.Store, error) {
	dbPath := getDBPath()
	return sqlite.New(dbPath)
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
	hk := hotkey.New([]hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.ModAlt,
	}, hotkey.KeyG)

	if err := hk.Register(); err != nil {
		return err
	}

	go func() {
		for range hk.Keydown() {
			logger.Debug("Global hotkey triggered")
			callback()
		}
	}()

	return nil
}
