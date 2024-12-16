package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
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

// initializeApplication handles app initialization
func initializeApplication() (*di.App, error) {
	logger.Info("Initializing dependency injection...")
	app, err := di.InitializeApp()
	if err != nil {
		return nil, err
	}
	logger.Info("Dependency injection initialized successfully")
	return app, nil
}

// runQuickNoteMode starts the application in quick-note mode
func runQuickNoteMode(app *di.App) {
	// Get the hotkey channel from the manager
	hotkeyEvents := app.GetHotkeyManager().GetEventChannel()

	// Start hotkey listener in background
	go func() {
		for {
			select {
			case <-hotkeyEvents:
				showQuickNote(app.GetTodoService())
			}
		}
	}()

	// Keep the application running
	select {}
}

// showQuickNote displays the quick-note UI
func showQuickNote(service *service.TodoService) {
	p := tea.NewProgram(ui.NewQuickNote(service))
	if _, err := p.Run(); err != nil {
		logger.Error("Quick note error: %v", err)
	}
}

func main() {
	flag.Parse()
	defer setupLogger()()

	logger.Info("Starting Godo application...")

	ctx := setupSignalHandler()

	app, err := initializeApplication()
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}

	if *fullUI {
		// Run full UI mode
		p := tea.NewProgram(ui.New(app.GetTodoService()))
		if _, err := p.Run(); err != nil {
			logger.Fatal("UI error: %v", err)
		}
	} else {
		// Run quick-note mode (default)
		runQuickNoteMode(app)
	}

	<-ctx.Done()
	logger.Info("Starting graceful shutdown...")
}
