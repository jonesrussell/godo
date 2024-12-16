// app.go
package di

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
	hotkeyManager *hotkey.HotkeyManager
	program       *tea.Program
	ui            *ui.TodoUI
}

// GetTodoService returns the todo service instance
func (a *App) GetTodoService() *service.TodoService {
	return a.todoService
}

// GetHotkeyManager returns the hotkey manager instance
func (a *App) GetHotkeyManager() *hotkey.HotkeyManager {
	return a.hotkeyManager
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

	// Wait for context cancellation
	<-ctx.Done()
	return ctx.Err()
}

func (a *App) initializeServices(ctx context.Context) error {
	logger.Info("Initializing services...")

	// Verify database connection
	testTodo, err := a.todoService.CreateTodo(ctx, "Test Todo", "Testing service initialization")
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
