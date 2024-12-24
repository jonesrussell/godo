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

type App struct {
	fyneApp    fyne.App
	mainWindow fyne.Window
	quickNote  QuickNoteService
	systray    systray.Interface
	store      storage.Store
	config     *config.Config
	log        logger.Logger
	hotkeyC    chan struct{}
}

func NewApp(cfg *config.Config, store storage.Store, log logger.Logger) *App {
	log.Debug("Creating new Fyne app")
	fyneApp := fyneapp.NewWithID("io.github.jonesrussell.godo")
	log.Debug("Created Fyne app", "type", fmt.Sprintf("%T", fyneApp))

	mainWindow := fyneApp.NewWindow(cfg.App.Name)
	log.Debug("Created main window")

	app := &App{
		fyneApp:    fyneApp,
		mainWindow: mainWindow,
		store:      store,
		config:     cfg,
		log:        log,
		hotkeyC:    make(chan struct{}, 1),
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

func (a *App) SetupUI() {
	a.setupSystemTray()
	a.setupMainWindow()

	// Setup lifecycle events
	a.fyneApp.Lifecycle().SetOnStarted(func() {
		a.log.Info("Lifecycle: Started")
		if err := a.setupGlobalHotkey(); err != nil {
			a.log.Error("Failed to setup global hotkey", "error", err)
			dialog.ShowError(fmt.Errorf("Failed to register global hotkey (%s). Quick note feature will not work.", a.config.Hotkeys.QuickNote), a.mainWindow)
		} else {
			a.log.Info("Global hotkey setup complete", "hotkey", a.config.Hotkeys.QuickNote)
		}

		// Start hotkey event handler
		go a.handleHotkeyEvents()
	})

	a.fyneApp.Lifecycle().SetOnStopped(func() {
		a.log.Info("Lifecycle: Stopped")
		close(a.hotkeyC)
	})
}

func (a *App) handleHotkeyEvents() {
	a.log.Debug("Starting hotkey event handler")
	for range a.hotkeyC {
		a.log.Debug("Hotkey event received - showing quick note")
		a.mainWindow.Show()
		a.quickNote.Show()
	}
	a.log.Debug("Hotkey event handler stopped")
}

func (a *App) setupGlobalHotkey() error {
	mapper := config.NewHotkeyMapper(a.log)
	mods, key, err := a.config.Hotkeys.QuickNote.Parse(mapper)
	if err != nil {
		a.log.Error("Failed to parse hotkey", "error", err)
		return fmt.Errorf("failed to parse hotkey: %w", err)
	}

	a.log.Debug("Creating hotkey", "modifiers", fmt.Sprintf("%v", mods), "key", fmt.Sprintf("%v", key))
	hk := hotkey.New(mods, key)

	a.log.Debug("Registering hotkey")
	if err := hk.Register(); err != nil {
		a.log.Error("Failed to register hotkey", "error", err)
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	a.log.Info("Global hotkey registered successfully", "hotkey", a.config.Hotkeys.QuickNote)

	// Start hotkey listener in a goroutine
	go func() {
		a.log.Debug("Starting hotkey event listener")
		keydownChan := hk.Keydown()
		if keydownChan == nil {
			a.log.Error("Keydown channel is nil")
			return
		}

		for {
			select {
			case _, ok := <-keydownChan:
				if !ok {
					a.log.Error("Keydown channel closed")
					return
				}
				a.log.Debug("Global hotkey triggered")
				select {
				case a.hotkeyC <- struct{}{}:
					a.log.Debug("Hotkey event sent to handler")
				default:
					a.log.Debug("Hotkey event dropped - handler busy")
				}
			}
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
