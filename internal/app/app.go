// Package app implements the main application logic for Godo.
package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow/systray"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/gui/theme"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
	"golang.design/x/hotkey"
)

// App represents the main application instance.
type App struct {
	fyneApp       fyne.App
	mainWindow    fyne.Window
	quickNote     QuickNoteService
	systray       systray.Interface
	store         storage.Store
	config        *config.Config
	log           logger.Logger
	hotkeyFactory HotkeyFactory
}

// NewApp creates a new application instance.
func NewApp(cfg *config.Config, store storage.Store, log logger.Logger, hotkeyFactory HotkeyFactory) *App {
	if hotkeyFactory == nil {
		hotkeyFactory = NewHotkeyFactory()
	}

	log.Debug("Creating new Fyne app")
	fyneApp := fyneapp.NewWithID("io.github.jonesrussell.godo")
	log.Debug("Created Fyne app", "type", fmt.Sprintf("%T", fyneApp))

	mainWindow := fyneApp.NewWindow(cfg.App.Name)
	log.Debug("Created main window")

	app := &App{
		fyneApp:       fyneApp,
		mainWindow:    mainWindow,
		store:         store,
		config:        cfg,
		log:           log,
		hotkeyFactory: hotkeyFactory,
	}

	log.Debug("Creating system tray service")
	app.systray = systray.New(fyneApp, log)

	quickNoteConfig := quicknote.Config{
		App:        fyneApp,
		MainWindow: mainWindow,
		Store:      store,
		Logger:     log,
	}
	app.quickNote = quicknote.New(quickNoteConfig)

	return app
}

// SetupUI initializes the application UI and hotkeys.
func (a *App) SetupUI() {
	a.setupSystemTray()
	a.setupMainWindow()

	// Setup lifecycle events
	a.fyneApp.Lifecycle().SetOnStarted(func() {
		a.log.Info("Lifecycle: Started")
		if err := a.setupGlobalHotkey(); err != nil {
			a.log.Error("Failed to setup global hotkey", "error", err)
			msg := fmt.Sprintf("Failed to register global hotkey (%s).", string(a.config.Hotkeys.QuickNote))
			dialog.ShowError(fmt.Errorf("%s Quick note feature will not work", msg), a.mainWindow)
		} else {
			a.log.Info("Global hotkey setup complete", "hotkey", a.config.Hotkeys.QuickNote)
		}
	})

	a.fyneApp.Lifecycle().SetOnStopped(func() {
		a.log.Info("Lifecycle: Stopped")
	})
}

// setupGlobalHotkey configures the global hotkey for quick notes.
//
// Implementation Notes:
// 1. The hotkey is hardcoded to Ctrl+Alt+G for reliability
// 2. Uses golang.design/x/hotkey for cross-platform global hotkey support
// 3. Runs a dedicated goroutine to handle hotkey events
// 4. The hotkey registration must happen after the application starts
//
// Troubleshooting:
// - If hotkey doesn't work, check if another application has registered Ctrl+Alt+G
// - Verify the application has the necessary permissions for global hotkeys
// - On Windows, try running as administrator if hotkey registration fails
// - Check logs for "Global hotkey registered successfully" message
//
// Known Issues:
// - Configuration via config.Hotkeys.QuickNote is not currently used
// - The hotkey cannot be changed at runtime
// - Only supports the Ctrl+Alt+G combination
func (a *App) setupGlobalHotkey() error {
	a.log.Debug("Setting up hotkey")
	hk := a.hotkeyFactory.NewHotkey(getHotkeyModifiers(), hotkey.KeyG)

	if err := hk.Register(); err != nil {
		return err
	}

	go func() {
		for range hk.Keydown() {
			a.log.Debug("Global hotkey triggered")
			a.quickNote.Show()
		}
	}()

	return nil
}

func (a *App) Run() {
	a.fyneApp.Run()
}

func (a *App) setupSystemTray() {
	a.log.Debug("Starting system tray setup", "app_type", fmt.Sprintf("%T", a.fyneApp))

	// Load icons
	a.log.Debug("Loading system tray icons")
	systrayIcon := theme.GetSystrayIconResource()
	if systrayIcon == nil {
		a.log.Error("Failed to load system tray icon - resource is nil")
		return
	}
	a.log.Debug("Loaded system tray icon", "name", systrayIcon.Name(), "content_length", len(systrayIcon.Content()))

	appIcon := theme.GetAppIconResource()
	if appIcon == nil {
		a.log.Error("Failed to load app icon - resource is nil")
		return
	}
	a.log.Debug("Loaded app icon", "name", appIcon.Name(), "content_length", len(appIcon.Content()))

	// Set app icon
	a.fyneApp.SetIcon(appIcon)
	a.log.Debug("Set application icon")

	// Set systray icon first
	a.log.Debug("Setting system tray icon")
	a.systray.SetIcon(systrayIcon)

	// Create menu
	a.log.Debug("Creating system tray menu")
	menu := fyne.NewMenu("Godo",
		fyne.NewMenuItem("Show", func() {
			a.mainWindow.Show()
			a.mainWindow.CenterOnScreen()
		}),
		fyne.NewMenuItem("Quick Note", a.quickNote.Show),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			a.fyneApp.Quit()
		}),
	)
	a.log.Debug("Created system tray menu", "items", len(menu.Items))

	// Setup menu
	a.log.Debug("Setting up system tray with menu")
	a.systray.Setup(menu)

	if !a.systray.IsReady() {
		a.log.Error("System tray setup failed - not ready after setup")
	} else {
		a.log.Info("System tray setup completed successfully")
	}
}

func (a *App) setupMainWindow() {
	header := container.NewHBox(
		widget.NewLabel(a.config.App.Name),
		widget.NewSeparator(),
	)

	toolbar := container.NewHBox(
		widget.NewButton("New Todo", a.quickNote.Show),
		widget.NewSeparator(),
		widget.NewButton("Show All", func() {
			a.mainWindow.Show()
			a.mainWindow.CenterOnScreen()
		}),
	)

	content := container.NewVBox(
		widget.NewLabel("Your todos will appear here"),
	)

	hotkeyText := "Press " + string(a.config.Hotkeys.QuickNote) + " for quick notes"
	versionText := "v" + a.config.App.Version

	statusBar := container.NewHBox(
		widget.NewLabel(hotkeyText),
		widget.NewSeparator(),
		widget.NewLabel(versionText),
	)

	mainContent := container.NewBorder(
		container.NewVBox(header, toolbar),
		statusBar,
		nil,
		nil,
		content,
	)

	a.mainWindow.SetContent(mainContent)
	a.mainWindow.Resize(fyne.NewSize(800, 600))
	a.mainWindow.CenterOnScreen()
	a.mainWindow.Hide()
}

func (a *App) Cleanup() {
	if db, ok := a.store.(*sqlite.Store); ok {
		if err := db.Close(); err != nil {
			a.log.Error("Failed to close database", "error", err)
		}
	}
}

// SetQuickNoteService allows injection of a QuickNoteService for testing
func (a *App) SetQuickNoteService(service QuickNoteService) {
	a.quickNote = service
}

// ShowQuickNote exposes the quick note functionality for testing
func (a *App) ShowQuickNote() {
	a.quickNote.Show()
}

func (a *App) SaveNote(content string) error {
	return a.store.SaveNote(content)
}

func (a *App) GetNotes() ([]string, error) {
	return a.store.GetNotes()
}

func (a *App) GetVersion() string {
	return a.config.App.Version
}

// SetSystrayService allows injection of a SystemTrayService for testing
func (a *App) SetSystrayService(service systray.Interface) {
	a.systray = service
}

// GetMainWindow returns the main application window for testing purposes
func (a *App) GetMainWindow() fyne.Window {
	return a.mainWindow
}

// SetMainWindow allows setting the main window for testing purposes
func (a *App) SetMainWindow(window fyne.Window) {
	a.mainWindow = window
}
