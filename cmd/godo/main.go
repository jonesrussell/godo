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
	"github.com/jonesrussell/godo/internal/logger"
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

	// Run the application
	if err := application.Run(ctx); err != nil {
		logger.Error("Application error: %v", err)
		return
	}

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

	// Add menu items
	mQuickNote := tray.AddMenuItem("Quick Note", "Add a quick note")
	mQuit := tray.AddMenuItem("Quit", "Quit the application")

	// Handle menu items
	go handleMenuItems(ctx, quickNote, mQuickNote, mQuit, tray, application)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	select {
	case sig := <-sigChan:
		logger.Info("Received signal: %v", sig)
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	// Cleanup
	if err := cleanup(application); err != nil {
		logger.Error("Error during cleanup: %v", err)
	}
}

func handleMenuItems(ctx context.Context, quickNote ui.QuickNoteUI, mQuickNote, mQuit *systray.MenuItem, tray internalsystray.Manager, application *app.App) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-mQuickNote.ClickedCh:
			if err := quickNote.Show(ctx); err != nil {
				logger.Error("Failed to show quick note: %v", err)
			}
		case <-mQuit.ClickedCh:
			tray.Quit()
			return
		case note := <-quickNote.GetInput():
			if note != "" {
				_, err := application.GetTodoService().CreateTodo(ctx, "Quick Note", note)
				if err != nil {
					logger.Error("Failed to create todo: %v", err)
				}
			}
		}
	}
}

func cleanup(application *app.App) error {
	logger.Info("Cleaning up application...")
	if err := application.Cleanup(); err != nil {
		logger.Error("Failed to cleanup: %v", err)
		return fmt.Errorf("cleanup failed: %w", err)
	}
	return nil
}
