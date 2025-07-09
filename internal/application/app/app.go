// Package app implements the main application logic for Godo.
package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/jonesrussell/godo/internal/application/app/hotkey"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/api"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/quicknote"
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

	// Quick note window management
	quickNoteWindow quicknote.Interface
	quickNoteMu     sync.Mutex
}

// New creates a new application instance
func New(
	fyneApp fyne.App,
	cfg *config.Config,
	log logger.Logger,
	taskService service.TaskService,
	mainWindow gui.MainWindow,
	store storage.TaskStore,
) *App {
	apiRunner := api.NewRunner(taskService, log, &cfg.HTTP)

	// Create the App instance first
	app := &App{
		fyneApp:     fyneApp,
		mainWindow:  mainWindow,
		quickNote:   nil, // Will be created on-demand via factory
		hotkey:      nil, // Will be set up below
		apiRunner:   apiRunner,
		config:      cfg,
		logger:      log,
		taskService: taskService,
		store:       store,
	}

	// Create quick note window on UI thread during initialization
	log.Debug("Creating quick note window during app initialization")
	windowConfig := config.WindowConfig{
		Width:  cfg.UI.QuickNote.Width,
		Height: cfg.UI.QuickNote.Height,
	}
	app.quickNoteWindow = quicknote.New(fyneApp, store, log, windowConfig)
	app.quickNoteWindow.Initialize(fyneApp, log)
	log.Debug("Quick note window created during initialization")

	// Now set up hotkey manager with factory that can access the App instance
	log.Info("Creating hotkey manager", "config", fmt.Sprintf("%+v", cfg.Hotkeys))
	if hkm, err := hotkey.NewManager(log, &cfg.Hotkeys); err != nil {
		log.Warn("Failed to create hotkey manager, continuing without hotkeys", "error", err)
		app.hotkey = nil
	} else {
		log.Info("Hotkey manager created successfully")
		app.hotkey = hkm

		// Create a factory function for the quick note window that just returns the existing instance
		quickNoteFactory := func() hotkey.QuickNoteService {
			log.Debug("Returning existing quick note window via factory")
			return app.quickNoteWindow
		}

		app.hotkey.SetQuickNoteFactory(quickNoteFactory, &cfg.Hotkeys.QuickNote)
	}

	return app
}

// setupSystray initializes the system tray
func (a *App) setupSystray() error {
	a.logger.Debug("Setting up systray")

	_, ok := a.fyneApp.(desktop.App)
	if !ok {
		a.logger.Error("Desktop features not available for systray")
		return ErrDesktopFeaturesNotAvailable
	}

	// Use config file path or default
	logPath := a.config.Logger.FilePath
	if logPath == "" {
		logPath = "logs/godo.log"
	}
	errorLogPath := logPath + "-error"

	// Create a wrapper for the quick note that creates it on demand
	quickNoteWrapper := &quickNoteWrapper{
		app: a,        // Pass the App instance instead of just fyneApp
		log: a.logger, // Make sure logger is set
	}

	a.logger.Debug("Quick note wrapper created for systray")

	return systray.SetupSystray(
		a.fyneApp,
		a.mainWindow.GetWindow(),
		quickNoteWrapper,
		logPath,
		errorLogPath,
		a.logger,
	)
}

// quickNoteWrapper wraps the quick note creation to handle on-demand creation
type quickNoteWrapper struct {
	app *App // Changed to *App to access App instance
	log logger.Logger
}

func (w *quickNoteWrapper) Show() {
	w.log.Debug("quickNoteWrapper.Show() called")

	// Use the unified quick note window management with error handling
	defer func() {
		if r := recover(); r != nil {
			w.log.Error("Panic in quick note window creation", "error", r)
		}
	}()

	quickNoteWindow := w.app.getQuickNoteWindow()
	if quickNoteWindow == nil {
		w.log.Error("Failed to get quick note window for systray")
		return
	}

	quickNoteWindow.Show()
	w.log.Debug("Systray quick note window shown")
}

func (w *quickNoteWrapper) Hide() {
	quickNoteWindow := w.app.getQuickNoteWindow()
	if quickNoteWindow != nil {
		quickNoteWindow.Hide()
	}
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

// getQuickNoteWindow returns the existing quick note window instance
func (a *App) getQuickNoteWindow() quicknote.Interface {
	a.logger.Debug("getQuickNoteWindow called", "goroutine_id", getGoroutineID())

	a.quickNoteMu.Lock()
	defer a.quickNoteMu.Unlock()

	if a.quickNoteWindow != nil {
		a.logger.Debug("Returning existing quick note window")
		return a.quickNoteWindow
	}

	a.logger.Error("Quick note window is nil - this should not happen")
	return nil
}

// getGoroutineID returns a simple goroutine ID for debugging
func getGoroutineID() string {
	return fmt.Sprintf("%p", &struct{}{})
}
