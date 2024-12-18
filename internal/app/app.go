// app.go
package app

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// App represents the main application
type App struct {
	todoService   *service.TodoService
	hotkeyManager hotkey.HotkeyManager
	program       *tea.Program
	ui            *ui.TodoUI
	quickNote     ui.QuickNoteUI
}

// NewApp creates a new application instance
func NewApp(
	todoService *service.TodoService,
	hotkeyManager hotkey.HotkeyManager,
	program *tea.Program,
	ui *ui.TodoUI,
	quickNote ui.QuickNoteUI,
) *App {
	return &App{
		todoService:   todoService,
		hotkeyManager: hotkeyManager,
		program:       program,
		ui:            ui,
		quickNote:     quickNote,
	}
}

// GetTodoService returns the todo service instance
func (a *App) GetTodoService() *service.TodoService {
	return a.todoService
}

// GetHotkeyManager returns the hotkey manager instance
func (a *App) GetHotkeyManager() hotkey.HotkeyManager {
	return a.hotkeyManager
}

// GetProgram returns the Bubble Tea program instance
func (a *App) GetProgram() *tea.Program {
	return a.program
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	logger.Info("Starting application services...")

	if err := a.initializeServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Start hotkey manager
	if err := a.hotkeyManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start hotkey manager: %w", err)
	}

	// Start background service to handle hotkey events
	go func() {
		hotkeyEvents := a.hotkeyManager.GetEventChannel()
		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopping hotkey listener...")
				return
			case <-hotkeyEvents:
				logger.Info("Hotkey triggered - showing quick note")
				// Handle quick note through platform-specific UI
				// This will be implemented separately
				a.handleQuickNote(ctx)
			}
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	return ctx.Err()
}

func (a *App) handleQuickNote(ctx context.Context) error {
	logger.Info("Quick note triggered")

	if err := a.quickNote.Show(ctx); err != nil {
		return fmt.Errorf("failed to show quick note: %w", err)
	}

	return nil
}

func (a *App) initializeServices(ctx context.Context) error {
	logger.Info("Initializing services...")

	// Verify database connection
	testTodo, err := a.todoService.CreateTodo(ctx, "test", "Testing service initialization")
	if err != nil {
		logger.Error("Failed to verify database connection: %v", err)
		return fmt.Errorf("failed to verify database connection: %w", err)
	}

	// Clean up test todo
	if err := a.todoService.DeleteTodo(ctx, testTodo.ID); err != nil {
		logger.Error("Failed to cleanup test todo: %v", err)
		return fmt.Errorf("failed to cleanup test todo: %w", err)
	}

	logger.Info("Services initialized successfully")
	return nil
}

func (a *App) Cleanup() error {
	logger.Info("Cleaning up application resources")
	// Add any cleanup logic here
	return nil
}
