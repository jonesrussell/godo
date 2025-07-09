// Package app implements the main application logic for Godo.
package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/jonesrussell/godo/internal/application/app/hotkey"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/api"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/systray"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/platform"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
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
	log.Info("Creating hotkey manager", "config", fmt.Sprintf("%+v", cfg.Hotkeys))
	if hkm, err := hotkey.NewUnifiedManager(log, &cfg.Hotkeys); err != nil {
		log.Warn("Failed to create hotkey manager, continuing without hotkeys", "error", err)
		hotkeyManager = nil
	} else {
		log.Info("Hotkey manager created successfully")
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

	// Use config file path or default
	logPath := a.config.Logger.FilePath
	if logPath == "" {
		logPath = "logs/godo.log"
	}
	errorLogPath := logPath + "-error"

	return systray.SetupSystray(
		a.fyneApp,
		a.mainWindow.GetWindow(),
		a.quickNote,
		logPath,
		errorLogPath,
		a.logger,
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

	// Start the hotkey manager to begin listening for events
	if err := a.hotkey.Start(); err != nil {
		a.logger.Error("Failed to start hotkey manager", "error", err)
		return fmt.Errorf("failed to start hotkey manager: %w", err)
	}

	a.logger.Info("Hotkey system started successfully")
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

	// Debug information for troubleshooting
	a.logger.Info("Environment details",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid(),
		"display", os.Getenv("DISPLAY"),
		"username", os.Getenv("USERNAME"),
		"sessionname", os.Getenv("SESSIONNAME"),
		"headless", platform.IsHeadless(),
		"wsl2", platform.IsWSL2(),
		"supports_gui", platform.SupportsGUI())

	// Start API server
	a.apiRunner.Start(a.config.HTTP.Port)

	// Wait for API to be ready
	timeout := time.Duration(a.config.HTTP.StartupTimeout) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if !a.apiRunner.WaitForReady(timeout) {
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

	// Stop hotkey manager with timeout
	if a.hotkey != nil {
		a.logger.Info("Stopping hotkey manager")
		hotkeyCtx, hotkeyCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer hotkeyCancel()

		// Use a goroutine to stop hotkey with timeout
		done := make(chan error, 1)
		go func() {
			done <- a.hotkey.Stop()
		}()

		select {
		case err := <-done:
			if err != nil {
				a.logger.Error("Failed to stop hotkey manager", "error", err)
			} else {
				a.logger.Info("Hotkey manager stopped")
			}
		case <-hotkeyCtx.Done():
			a.logger.Warn("Hotkey manager stop timed out, forcing cleanup")
		}
	}

	// Stop API server with timeout
	timeout := time.Duration(a.config.HTTP.ShutdownTimeout) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

func (a *App) ForceKillTimeout() time.Duration {
	return time.Duration(a.config.App.ForceKillTimeout) * time.Second
}
