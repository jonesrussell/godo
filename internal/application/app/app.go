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
	"fyne.io/fyne/v2/widget"

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

	// Unified window management
	testWindow *testWindow
	windowMu   sync.Mutex
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

	// Now set up hotkey manager with factory that can access the App instance
	log.Info("Creating hotkey manager", "config", fmt.Sprintf("%+v", cfg.Hotkeys))
	if hkm, err := hotkey.NewUnifiedManager(log, &cfg.Hotkeys); err != nil {
		log.Warn("Failed to create hotkey manager, continuing without hotkeys", "error", err)
		app.hotkey = nil
	} else {
		log.Info("Hotkey manager created successfully")
		app.hotkey = hkm

		// Create a factory function for the quick note window
		quickNoteFactory := func() hotkey.QuickNoteService {
			log.Debug("Creating test window via factory")
			// Use the unified window management
			return app.getOrCreateTestWindow()
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

// testWindow is a simple test window for debugging
type testWindow struct {
	window fyne.Window
	log    logger.Logger
	closed bool
}

func (t *testWindow) Show() {
	t.log.Debug("testWindow.Show() called")

	// Check if window was closed and needs to be recreated
	if t.closed || t.window == nil {
		t.log.Debug("testWindow.Show() - window was closed, cannot show")
		return
	}

	// Window is already shown during creation, no additional operations needed
	t.log.Debug("testWindow.Show() - window already visible")
	t.log.Debug("testWindow.Show() - completed")
}

func (t *testWindow) Hide() {
	if t.window != nil && !t.closed {
		fyne.Do(func() {
			t.window.Hide()
		})
	}
}

// Close marks the window as closed
func (t *testWindow) Close() {
	t.closed = true
	if t.window != nil {
		fyne.Do(func() {
			t.window.Close()
		})
	}
}

// quickNoteWrapper wraps the quick note creation to handle on-demand creation
type quickNoteWrapper struct {
	app    *App // Changed to *App to access App instance
	store  storage.TaskStore
	log    logger.Logger
	config config.WindowConfig
	window gui.QuickNote
}

func (w *quickNoteWrapper) Show() {
	w.log.Debug("quickNoteWrapper.Show() called")

	if w.window == nil {
		w.log.Debug("Creating test window for systray")
		// Use the unified window management with error handling
		defer func() {
			if r := recover(); r != nil {
				w.log.Error("Panic in window creation", "error", r)
			}
		}()

		w.window = w.app.getOrCreateTestWindow()
		if w.window == nil {
			w.log.Error("Failed to create test window for systray")
			return
		}
		w.log.Debug("Systray window creation completed")
	}

	// Check if window is valid before trying to show it
	if w.window != nil {
		w.log.Debug("Systray window already visible")
	} else {
		w.log.Error("Systray window is nil after creation")
	}
}

func (w *quickNoteWrapper) Hide() {
	if w.window != nil {
		fyne.Do(func() {
			w.window.Hide()
		})
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

// getOrCreateTestWindow ensures a single test window instance
func (a *App) getOrCreateTestWindow() *testWindow {
	a.windowMu.Lock()
	defer a.windowMu.Unlock()

	// Check if window exists and is not closed
	if a.testWindow != nil && !a.testWindow.closed && a.testWindow.window != nil {
		a.logger.Debug("Reusing existing test window")
		return a.testWindow
	}

	// Window doesn't exist or was closed, create a new one
	a.logger.Debug("Creating unified test window")

	// Create window following Fyne best practices from latest documentation
	var window fyne.Window
	var createErr error

	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("Panic during window creation", "error", r)
			createErr = fmt.Errorf("panic during window creation: %v", r)
		}
	}()

	// Always use fyne.DoAndWait() but handle the case where we're already on the UI thread
	// by using a timeout to detect if we're already on the main thread
	done := make(chan bool, 1)

	go func() {
		fyne.DoAndWait(func() {
			a.logger.Debug("Inside fyne.DoAndWait - creating window")
			window = a.fyneApp.NewWindow("Test Window - Unified")
			if window == nil {
				a.logger.Error("NewWindow returned nil")
				createErr = fmt.Errorf("NewWindow returned nil")
				return
			}

			label := widget.NewLabel("Test Window - Unified Works!")
			window.SetContent(label)

			// Set size before showing (as per latest Fyne docs)
			window.Resize(fyne.NewSize(300, 200))

			// Show the window (as per latest Fyne docs)
			window.Show()
			a.logger.Debug("Window created and shown successfully")
		})
		done <- true
	}()

	// Wait for window creation with timeout
	select {
	case <-done:
		a.logger.Debug("Window creation completed")
	case <-time.After(100 * time.Millisecond):
		a.logger.Debug("Window creation timed out, likely already on UI thread")
		// If we timeout, we're probably already on the UI thread, so create directly
		window = a.fyneApp.NewWindow("Test Window - Unified")
		if window == nil {
			a.logger.Error("NewWindow returned nil")
			createErr = fmt.Errorf("NewWindow returned nil")
		} else {
			label := widget.NewLabel("Test Window - Unified Works!")
			window.SetContent(label)
			window.Resize(fyne.NewSize(300, 200))
			window.Show()
			a.logger.Debug("Window created and shown successfully (direct)")
		}
	}

	if createErr != nil {
		a.logger.Error("Failed to create window", "error", createErr)
		return nil
	}

	if window == nil {
		a.logger.Error("Window is nil after creation")
		return nil
	}

	a.testWindow = &testWindow{
		window: window,
		log:    a.logger,
		closed: false,
	}

	// Set up close handler to mark window as closed when user closes it
	fyne.Do(func() {
		window.SetCloseIntercept(func() {
			a.logger.Debug("Test window close intercepted")
			a.testWindow.closed = true
			window.Hide() // Hide instead of close to prevent crashes
		})
	})

	a.logger.Debug("Unified test window created successfully")

	return a.testWindow
}
