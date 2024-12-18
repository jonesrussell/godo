// app.go
package app

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/quicknote"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// App represents the main application
type App struct {
	config      *config.Config
	todoService *service.TodoService
	program     *tea.Program
	todoUI      *ui.TodoUI
	quickNote   quicknote.UI
}

// NewApp creates a new App instance with all dependencies
func NewApp(
	cfg *config.Config,
	todoService *service.TodoService,
	program *tea.Program,
	todoUI *ui.TodoUI,
	quickNote quicknote.UI,
) *App {
	return &App{
		config:      cfg,
		todoService: todoService,
		program:     program,
		todoUI:      todoUI,
		quickNote:   quickNote,
	}
}

// GetTodoService returns the todo service instance
func (a *App) GetTodoService() *service.TodoService {
	return a.todoService
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	logger.Info("Starting application services...")

	if err := a.initializeServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Wait for context cancellation
	<-ctx.Done()
	return ctx.Err()
}

func (a *App) initializeServices(ctx context.Context) error {
	logger.Info("Initializing services...")

	if err := a.verifyDatabaseConnection(ctx); err != nil {
		return err
	}

	logger.Info("Services initialized successfully")
	return nil
}

func (a *App) verifyDatabaseConnection(ctx context.Context) error {
	testTodo, err := a.todoService.CreateTodo(ctx, "test", "Testing service initialization")
	if err != nil {
		logger.Error("Failed to verify database connection: %v", err)
		return fmt.Errorf("failed to verify database connection: %w", err)
	}

	if err := a.todoService.DeleteTodo(ctx, testTodo.ID); err != nil {
		logger.Error("Failed to cleanup test todo: %v", err)
		return fmt.Errorf("failed to cleanup test todo: %w", err)
	}

	return nil
}

// Cleanup performs any necessary cleanup before shutdown
func (a *App) Cleanup() error {
	logger.Info("Cleaning up application resources")
	// Add any cleanup logic here
	return nil
}
