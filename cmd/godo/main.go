package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/di"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/ui"
)

func main() {
	// Load configuration
	env := os.Getenv("GODO_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger with config
	if err := logger.InitializeWithConfig(cfg.Logging); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Create application instance with config
	app, err := di.InitializeAppWithConfig(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start systray
	go systray.Run(func() {
		onSystrayReady(ctx, app)
	}, onSystrayExit)

	// Wait for signal
	select {
	case sig := <-sigChan:
		logger.Info("Received signal: %v", sig)
		cancel()
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	// Cleanup
	if err := cleanup(app); err != nil {
		logger.Error("Error during cleanup: %v", err)
	}
}

func onSystrayReady(ctx context.Context, app *di.App) {
	// Initialize QuickNoteUI
	quickNote, err := ui.NewQuickNoteUI()
	if err != nil {
		logger.Error("Failed to create quick note UI: %v", err)
		return
	}

	// Set up systray
	if err := ui.SetupSystray(); err != nil {
		logger.Error("Failed to setup systray: %v", err)
		// Continue running even if systray setup fails
	}

	mQuickNote := systray.AddMenuItem("Quick Note", "Add a quick note")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Handle menu items
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-mQuickNote.ClickedCh:
				if err := quickNote.Show(ctx); err != nil {
					logger.Error("Failed to show quick note: %v", err)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case note := <-quickNote.GetInput():
				if note != "" {
					_, err := app.GetTodoService().CreateTodo(ctx, "Quick Note", note)
					if err != nil {
						logger.Error("Failed to create todo: %v", err)
					}
				}
			}
		}
	}()

	// Start the application
	if err := app.Run(ctx); err != nil {
		logger.Error("Application error: %v", err)
	}
}

func onSystrayExit() {
	logger.Info("Systray exiting")
}

func cleanup(app *di.App) error {
	logger.Info("Cleaning up application...")
	if err := app.Cleanup(); err != nil {
		logger.Error("Failed to cleanup: %v", err)
		return fmt.Errorf("cleanup failed: %w", err)
	}
	return nil
}
