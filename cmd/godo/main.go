package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/quicknote"
	internalsystray "github.com/jonesrussell/godo/internal/systray"
	"github.com/jonesrussell/godo/internal/ui"
)

func main() {
	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	env := os.Getenv("GODO_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		return
	}

	// Initialize logger with config
	if err := logger.InitializeWithConfig(cfg.Logging); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		return
	}

	// Create application instance with config
	application, err := app.InitializeAppWithConfig(cfg)
	if err != nil {
		logger.Error("Failed to initialize application: %v", err)
		return
	}
	defer cleanup(application)

	// Initialize QuickNoteUI
	quickNote, err := ui.NewQuickNoteUI()
	if err != nil {
		logger.Error("Failed to create quick note UI: %v", err)
		return
	}

	// Set up systray
	tray, err := internalsystray.SetupSystray()
	if err != nil {
		logger.Error("Failed to setup systray: %v", err)
		return
	}
	defer tray.Quit()

	// Add menu items in order
	mOpen := tray.AddMenuItem("Open", "Open Godo")
	mQuickNote := tray.AddMenuItem("Quick Note", "Add a quick note")
	systray.AddSeparator()
	mQuit := tray.AddMenuItem("Quit", "Quit the application")

	// Create error channel for goroutines
	errChan := make(chan error, 1)

	// Start hotkey handler
	go func() {
		if err := handleHotkeys(ctx, application.GetHotkeyManager(), quickNote); err != nil {
			errChan <- fmt.Errorf("hotkey handler error: %w", err)
		}
	}()

	// Start menu item handler
	go func() {
		if err := handleMenuItems(ctx, cancel, quickNote, mOpen, mQuickNote, mQuit, application); err != nil {
			errChan <- err
		}
	}()

	// Run the application in a goroutine
	go func() {
		if err := application.Run(ctx); err != nil {
			errChan <- fmt.Errorf("application error: %w", err)
		}
	}()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either signal, error, or context cancellation
	select {
	case sig := <-sigChan:
		logger.Info("Received signal: %v", sig)
	case err := <-errChan:
		logger.Error("Error occurred: %v", err)
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}
}

func handleMenuItems(ctx context.Context, cancel context.CancelFunc, quickNote quicknote.UI, mOpen, mQuickNote, mQuit *systray.MenuItem, application *app.App) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-mOpen.ClickedCh:
			// Handle opening the main application window
			logger.Debug("Open menu item clicked")
			// TODO: Implement opening main window
		case <-mQuickNote.ClickedCh:
			if err := quickNote.Show(ctx); err != nil {
				logger.Error("Failed to show quick note: %v", err)
				return fmt.Errorf("quick note error: %w", err)
			}
		case <-mQuit.ClickedCh:
			cancel()
			return nil
		case note := <-quickNote.GetInput():
			if note != "" {
				if _, err := application.GetTodoService().CreateTodo(ctx, "Quick Note", note); err != nil {
					logger.Error("Failed to create todo: %v", err)
					return fmt.Errorf("failed to create todo: %w", err)
				}
			}
		}
	}
}

func handleHotkeys(ctx context.Context, hotkeyManager hotkey.HotkeyManager, quickNote quicknote.UI) error {
	eventChan := hotkeyManager.GetEventChannel()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-eventChan:
			logger.Debug("Hotkey triggered")
			if err := quickNote.Show(ctx); err != nil {
				logger.Error("Failed to show quick note from hotkey: %v", err)
				return fmt.Errorf("quick note error: %w", err)
			}
		}
	}
}

func cleanup(application *app.App) {
	logger.Info("Cleaning up application...")
	if err := application.Cleanup(); err != nil {
		logger.Error("Failed to cleanup: %v", err)
	}
}
