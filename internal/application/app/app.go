// Package app implements the main application logic for Godo.
package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/jonesrussell/godo/internal/application/app/hotkey"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/api"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/systray"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
	"github.com/jonesrussell/godo/internal/shared/config"
)

// Constants for configuration values
const (
	DefaultAPIPort     = 8080
	APIStartupTimeout  = 5 * time.Second
	APIShutdownTimeout = 5 * time.Second
)

// App represents the main application
type App struct {
	fyneApp     fyne.App
	mainWindow  gui.MainWindow
	quickNote   gui.QuickNote
	hotkey      hotkey.Manager
	apiRunner   *api.Runner
	config      *config.Config
	logger      logger.Logger
	taskService service.TaskService
	store       storage.TaskStore
}

// New creates a new application instance
func New(
	cfg *config.Config,
	log logger.Logger,
	taskService service.TaskService,
	mainWindow gui.MainWindow,
	quickNote gui.QuickNote,
	store storage.TaskStore,
) *App {
	fyneApp := app.New()

	var hotkeyManager hotkey.Manager
	if hkm, err := hotkey.NewLinuxManager(log); err != nil {
		log.Warn("Failed to create hotkey manager, continuing without hotkeys", "error", err)
		hotkeyManager = nil
	} else {
		hotkeyManager = hkm
		hotkeyManager.SetQuickNote(quickNote, &cfg.Hotkeys.QuickNote)
	}

	apiRunner := api.NewRunner(taskService, log, &cfg.HTTP)

	return &App{
		fyneApp:     fyneApp,
		mainWindow:  mainWindow,
		quickNote:   quickNote,
		hotkey:      hotkeyManager,
		apiRunner:   apiRunner,
		config:      cfg,
		logger:      log,
		taskService: taskService,
		store:       store,
	}
}

// setupSystray initializes the system tray
func (a *App) setupSystray() error {
	_, ok := a.fyneApp.(desktop.App)
	if !ok {
		return ErrDesktopFeaturesNotAvailable
	}

	return systray.SetupSystray(
		a.fyneApp,
		a.mainWindow.GetWindow(),
		a.quickNote,
		a.config.Logger.FilePath,
		a.config.Logger.FilePath+"-error",
	)
}

// setupHotkey sets up the global hotkey
func (a *App) setupHotkey() error {
	if a.hotkey == nil {
		a.logger.Warn("No hotkey manager available")
		return nil
	}

	if err := a.hotkey.Register(); err != nil {
		if errors.Is(err, hotkey.ErrWSL2NotSupported) {
			a.logger.Warn("Hotkey system not available in WSL2 environment")
			return nil
		}
		return fmt.Errorf("failed to setup hotkey system: %w", err)
	}
	return nil
}

// SetupUI initializes the user interface components
func (a *App) SetupUI() error {
	a.logger.Debug("Setting up UI components")

	// Set up systray
	if err := a.setupSystray(); err != nil {
		if errors.Is(err, ErrDesktopFeaturesNotAvailable) {
			a.logger.Warn("Desktop features not available, skipping systray setup")
		} else {
			a.logger.Warn("Failed to setup systray", "error", err)
		}
		// Continue without systray
	}

	// Show main window if not configured to start hidden
	if !a.config.UI.MainWindow.StartHidden {
		a.mainWindow.Show()
	}

	return nil
}

// Run starts the application
func (a *App) Run() {
	a.logger.Info("Starting Godo application")

	// Start API server
	a.apiRunner.Start(DefaultAPIPort)

	// Wait for API to be ready
	if !a.apiRunner.WaitForReady(APIStartupTimeout) {
		a.logger.Error("API server failed to start within timeout")
		return
	}

	// Setup UI
	if err := a.SetupUI(); err != nil {
		a.logger.Error("Failed to setup UI", "error", err)
		return
	}

	// Setup hotkey
	if err := a.setupHotkey(); err != nil {
		a.logger.Error("Failed to setup hotkey", "error", err)
		// Continue running even if hotkey fails
	}

	// Run the application (this blocks until quit)
	a.fyneApp.Run()

	// When we get here, the app was quit via GUI
	a.logger.Info("Application shutting down")
}

// Cleanup performs cleanup operations before shutdown
func (a *App) Cleanup() {
	a.logger.Info("Cleaning up application")

	// Stop API server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), APIShutdownTimeout)
	defer cancel()

	if err := a.apiRunner.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to stop API server", "error", err)
	}
}

// Quit performs cleanup and quits the application
func (a *App) Quit() {
	a.Cleanup()
	if a.fyneApp != nil {
		a.fyneApp.Quit()
	}
}

// Logger returns the application logger
func (a *App) Logger() logger.Logger {
	return a.logger
}

// Store returns the application store
func (a *App) Store() storage.TaskStore {
	return a.store
}
