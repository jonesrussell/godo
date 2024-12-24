package app

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/gui/systray"
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
}

func NewApp(cfg *config.Config, store storage.Store, log logger.Logger) *App {
	fyneApp := fyneapp.NewWithID("io.github.jonesrussell.godo")
	mainWindow := fyneApp.NewWindow(cfg.App.Name)

	app := &App{
		fyneApp:    fyneApp,
		mainWindow: mainWindow,
		store:      store,
		config:     cfg,
		log:        log,
		systray:    systray.New(fyneApp, log),
	}

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
	a.setupLifecycleLogging()
	a.setupSystemTray()
	a.setupMainWindow()
	if err := a.setupGlobalHotkey(); err != nil {
		a.log.Error("Failed to setup global hotkey", "error", err)
	}
}

func (a *App) Run() {
	a.fyneApp.Run()
}

func (a *App) setupLifecycleLogging() {
	a.fyneApp.Lifecycle().SetOnStarted(func() {
		a.log.Info("Lifecycle: Started")
	})
	a.fyneApp.Lifecycle().SetOnStopped(func() {
		a.log.Info("Lifecycle: Stopped")
	})
}

func (a *App) setupSystemTray() {
	a.log.Debug("Loading system tray icon")
	systrayIcon := assets.GetSystrayIconResource()
	appIcon := assets.GetAppIconResource()
	a.fyneApp.SetIcon(appIcon)

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

	// Setup menu first
	a.systray.Setup(menu)

	// Then set icon
	if a.systray.IsReady() {
		a.systray.SetIcon(systrayIcon)
		a.log.Info("System tray setup complete")
	} else {
		a.log.Warn("System tray not ready, icon not set")
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

func (a *App) setupGlobalHotkey() error {
	hk := hotkey.New([]hotkey.Modifier{
		config.ModCtrl,
		config.ModAlt,
	}, hotkey.KeyG)

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
