// Package app implements the main application logic for Godo.
package core

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

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/api"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/quicknote"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/theme"
	"github.com/jonesrussell/godo/internal/infrastructure/hotkey"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/platform"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

// App represents the main application
type App struct {
	fyneApp     fyne.App
	mainWindow  gui.MainWindow
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
		hotkey:      nil, // Will be set up below
		apiRunner:   apiRunner,
		config:      cfg,
		logger:      log,
		taskService: taskService,
		store:       store,
	}

	// Create quick note window during initialization
	// Since we're on the main goroutine, we can call Fyne functions directly
	log.Debug("Creating quick note window during app initialization")
	windowConfig := config.WindowConfig{
		Width:  cfg.UI.QuickNote.Width,
		Height: cfg.UI.QuickNote.Height,
	}

	app.quickNoteWindow = quicknote.New(fyneApp, store, log, windowConfig)
	app.quickNoteWindow.Initialize(fyneApp, log)
	log.Debug("Quick note window created during initialization")

	// Now set up hotkey manager with the simplified interface
	log.Info("Creating hotkey manager", "config", fmt.Sprintf("%+v", cfg.Hotkeys))
	if hkm, err := hotkey.NewManager(log, &cfg.Hotkeys); err != nil {
		log.Warn("Failed to create hotkey manager, continuing without hotkeys", "error", err)
		app.hotkey = nil
	} else {
		log.Info("Hotkey manager created successfully")
		app.hotkey = hkm

		// Use the simplified interface - pass the quick note service directly
		app.hotkey.SetQuickNote(app.quickNoteWindow, &cfg.Hotkeys.QuickNote)
	}

	return app
}

// setupSystray initializes the system tray using the simplified Fyne approach
func (a *App) setupSystray() error {
	a.logger.Debug("Setting up systray")

	// Since we're on the main goroutine, we can call Fyne functions directly
	desk, ok := a.fyneApp.(desktop.App)
	if !ok {
		a.logger.Warn("Desktop features not available for systray")
		return nil // Not an error, just not supported
	}

	// Create menu using the standard Fyne menu API
	m := fyne.NewMenu("Godo",
		fyne.NewMenuItem("Show", func() {
			a.logger.Debug("Systray Show menu item tapped")
			// Ensure window operations happen on UI thread
			fyne.Do(func() {
				a.mainWindow.Show()
				a.mainWindow.GetWindow().RequestFocus()
			})
		}),
		fyne.NewMenuItem("Quick Note", func() {
			a.logger.Debug("Systray Quick Note menu item tapped")
			// Ensure window operations happen on UI thread
			fyne.Do(func() {
				quickNoteWindow := a.getQuickNoteWindow()
				if quickNoteWindow != nil {
					quickNoteWindow.Show()
				}
			})
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			a.logger.Debug("Systray Quit menu item tapped")
			// Quit can be called from any thread
			a.Quit()
		}),
	)

	// Set the system tray menu and icon on the UI thread
	fyne.Do(func() {
		desk.SetSystemTrayMenu(m)

		// Set the system tray icon
		icon := theme.AppIcon()
		if icon != nil {
			desk.SetSystemTrayIcon(icon)
			a.logger.Debug("Systray icon set successfully")
		} else {
			a.logger.Warn("Failed to get systray icon - systray will have no icon")
		}
	})

	a.logger.Info("Systray setup completed")
	return nil
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
		a.logger.Warn("Failed to setup systray", "error", err)
		// Continue without systray
	}

	// Show main window if not configured to start hidden
	if !a.config.UI.MainWindow.StartHidden {
		// Since we're on the main goroutine, we can call Fyne functions directly
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

	// Start API server (this is safe to do from any thread)
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

	// Setup UI (this contains UI thread operations)
	if err := a.SetupUI(); err != nil {
		a.logger.Error("Failed to setup UI", "error", err)
		return
	}

	// Setup hotkey (this is safe to do from any thread)
	if err := a.setupHotkey(); err != nil {
		a.logger.Error("Failed to setup hotkey", "error", err)
		// Continue running even if hotkey fails
	}

	// Run the application (this blocks until quit)
	// This MUST be called from the main thread
	a.fyneApp.Run()

	// When we get here, the app was quit via GUI
	a.logger.Info("Application shutting down")
}

// Cleanup performs cleanup operations before shutdown
func (a *App) Cleanup() {
	a.logger.Info("Cleaning up application")

	// Clean up quick note window
	if a.quickNoteWindow != nil {
		a.logger.Info("Cleaning up quick note window")
		// Use fyne.DoAndWait to ensure this runs on the main thread
		fyne.DoAndWait(func() {
			a.quickNoteWindow.Hide() // Hide the window if it's visible
		})
	}

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
	// First perform cleanup
	a.Cleanup()

	// Then quit the application on the UI thread
	if a.fyneApp != nil {
		// Determine if we're already on the UI thread
		// If we are, just call Quit directly; if not, use fyne.Do
		done := make(chan struct{})
		timeout := time.After(2 * time.Second)

		go func() {
			defer close(done)
			fyne.Do(func() {
				a.fyneApp.Quit()
			})
		}()

		// Wait for quit with timeout
		select {
		case <-done:
			a.logger.Info("Application quit completed")
		case <-timeout:
			a.logger.Warn("Application quit timed out, forcing exit")
		}
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
