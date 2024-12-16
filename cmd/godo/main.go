package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/di"
	"github.com/jonesrussell/godo/internal/icon"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

var (
	fullUI = flag.Bool("ui", false, "Launch full todo management interface")
)

// setupSignalHandler creates a signal handler
func setupSignalHandler(parentCtx context.Context) context.Context {
	// Use the parent context
	ctx, cancel := context.WithCancel(parentCtx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer cancel() // Ensure cancel is called when the goroutine exits
		select {
		case sig := <-sigChan:
			logger.Info("Received signal: %v", sig)
			cancel()
		case <-parentCtx.Done():
			// Parent context was cancelled
		}
	}()

	return ctx
}

// showQuickNote displays a minimal quick-note window
func showQuickNote(service *service.TodoService) {
	qn := ui.NewQuickNote(service)
	qn.Show()
}

// runBackgroundService handles hotkey events and background operations
func runBackgroundService(ctx context.Context, app *di.App) {
	logger.Info("Starting background service...")

	// Start the hotkey manager with context
	if err := app.GetHotkeyManager().Start(ctx); err != nil {
		logger.Error("Failed to start hotkey manager: %v", err)
		return
	}

	hotkeyEvents := app.GetHotkeyManager().GetEventChannel()
	logger.Info("Listening for hotkey events (Ctrl+Alt+G)...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Background service shutting down...")
			if err := app.GetHotkeyManager().Cleanup(); err != nil {
				logger.Error("Error cleaning up hotkey: %v", err)
			}
			return
		case <-hotkeyEvents:
			logger.Info("Hotkey triggered - showing quick note")
			showQuickNote(app.GetTodoService())
			logger.Debug("Quick note window closed")
		}
	}
}

// onReady is called when systray is ready
func onReady(ctx context.Context, app *di.App, cancel context.CancelFunc) func() {
	return func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("Godo")
		systray.SetTooltip("Quick Todo Manager")

		systray.AddSeparator()
		mOpen := systray.AddMenuItem("Open Manager", "Open todo manager")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit application")

		// Start background service
		go runBackgroundService(ctx, app)

		// Handle menu items
		for {
			select {
			case <-mQuit.ClickedCh:
				cancel()
				systray.Quit()
				return
			case <-mOpen.ClickedCh:
				showFullUI(app.GetTodoService())
			case <-ctx.Done():
				systray.Quit()
				return
			}
		}
	}
}

// onExit is called when systray is quitting
func onExit() {
	logger.Info("Cleaning up...")
	os.Exit(0) // Ensure the process exits
}

// showFullUI displays the full todo management interface
func showFullUI(service *service.TodoService) {
	p := tea.NewProgram(
		ui.New(service),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		logger.Error("UI error: %v", err)
	}
}

func main() {
	// Initialize logger first
	cleanup := logger.Initialize()
	defer cleanup()

	flag.Parse()

	logger.Info("Starting Godo application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handler with the cancellable context
	sigCtx := setupSignalHandler(ctx)

	app, err := di.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}

	if *fullUI {
		showFullUI(app.GetTodoService())
	} else {
		// Run in system tray with quit handler
		go systray.Run(onReady(sigCtx, app, cancel), onExit)

		// Wait for context cancellation
		<-sigCtx.Done()
		logger.Info("Starting graceful shutdown...")
	}
}
