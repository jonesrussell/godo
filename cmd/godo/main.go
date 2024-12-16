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
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

var (
	fullUI = flag.Bool("ui", false, "Launch full todo management interface")
)

// setupLogger configures and manages logger lifecycle
func setupLogger() func() {
	return func() {
		if err := logger.Sync(); err != nil {
			os.Stderr.WriteString("Failed to sync logger: " + err.Error() + "\n")
		}
	}
}

// setupSignalHandler creates a signal handler
func setupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal: %v", sig)
		cancel()
	}()

	return ctx
}

// showQuickNote displays a minimal quick-note UI
func showQuickNote(service *service.TodoService) {
	p := tea.NewProgram(
		ui.NewQuickNote(service),
		tea.WithAltScreen(),        // Use alternate screen
		tea.WithMouseCellMotion(),  // Enable mouse support
		tea.WithoutSignalHandler(), // Don't handle signals
	)

	// Run in a goroutine to not block
	go func() {
		if _, err := p.Run(); err != nil {
			logger.Error("Quick note error: %v", err)
		}
	}()
}

// runBackgroundService handles hotkey events and background operations
func runBackgroundService(ctx context.Context, app *di.App) {
	logger.Info("Starting background service...")

	// Get hotkey channel
	hotkeyEvents := app.GetHotkeyManager().GetEventChannel()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Background service shutting down...")
			return
		case <-hotkeyEvents:
			logger.Debug("Received hotkey event")
			showQuickNote(app.GetTodoService())
		}
	}
}

// onReady is called when systray is ready
func onReady(ctx context.Context, app *di.App) func() {
	return func() {
		systray.SetIcon(nil) // Set your icon here
		systray.SetTitle("Godo")
		systray.SetTooltip("Quick Todo Manager")

		mOpen := systray.AddMenuItem("Open Manager", "Open todo manager")
		mQuit := systray.AddMenuItem("Quit", "Quit application")

		// Handle menu items
		go func() {
			for {
				select {
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				case <-mOpen.ClickedCh:
					showFullUI(app.GetTodoService())
				case <-ctx.Done():
					return
				}
			}
		}()

		// Start background service
		go runBackgroundService(ctx, app)
	}
}

// onExit is called when systray is quitting
func onExit() {
	logger.Info("Cleaning up...")
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
	flag.Parse()
	defer setupLogger()()

	logger.Info("Starting Godo application...")

	ctx := setupSignalHandler()

	app, err := di.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}

	if *fullUI {
		showFullUI(app.GetTodoService())
	} else {
		// Run in system tray
		go systray.Run(onReady(ctx, app), onExit)

		// Wait for context cancellation
		<-ctx.Done()
		logger.Info("Starting graceful shutdown...")
	}
}
