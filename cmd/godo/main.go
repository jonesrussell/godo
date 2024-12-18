package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/di"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/ui"
)

func main() {
	// Initialize logger
	if err := logger.Initialize(); err != nil {
		logger.Fatal("Failed to initialize logger: %v", err)
	}

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create application instance
	app, err := di.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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

	// Set up systray icon and menu
	systray.SetIcon(getIcon())
	systray.SetTitle("Godo")
	systray.SetTooltip("Godo - Quick Notes & Tasks")

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
	logger.Info("Cleaning up...")
	// Add any necessary cleanup here
	return nil
}

func getIcon() []byte {
	// Replace with your actual icon data
	// For now, return a minimal icon
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		// ... rest of icon data ...
	}
}
