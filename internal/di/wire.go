//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/repository"
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

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	logger.Info("Initializing services...")
	if err := a.initializeServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	logger.Info("Services initialized successfully")

	// Start hotkey listener in a goroutine
	logger.Info("Starting hotkey listener...")
	go func() {
		if err := a.hotkeyManager.Start(ctx); err != nil {
			logger.Error("Hotkey error: %v", err)
		}
	}()
	logger.Info("Hotkey listener started")

	return nil
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

// NewApp creates a new App instance
func NewApp(todoService *service.TodoService, ui *ui.TodoUI) (*App, error) {
	logger.Debug("Creating new App instance")
	program := tea.NewProgram(ui)

	showUI := func() {
		logger.Debug("ShowUI callback triggered")
		program.Send(struct{}{})
	}

	hotkeyManager := hotkey.New(showUI)
	if hotkeyManager == nil {
		return nil, fmt.Errorf("failed to create hotkey manager")
	}

	return &App{
		todoService:   todoService,
		hotkeyManager: hotkeyManager,
		program:       program,
		ui:            ui,
	}, nil
}

// Provide TodoRepository interface implementation
func provideTodoRepository(db *sql.DB) repository.TodoRepository {
	logger.Debug("Creating new TodoRepository")
	return repository.NewSQLiteTodoRepository(db)
}

// Add provider set
var DefaultSet = wire.NewSet(
	NewSQLiteDB,
	provideTodoRepository,
	service.NewTodoService,
	provideUI,
	NewApp,
)

// Add database provider
func NewSQLiteDB() (*sql.DB, error) {
	logger.Debug("Creating new SQLite database connection")
	return database.NewSQLiteDB("./godo.db")
}

func provideUI(todoService *service.TodoService) *ui.TodoUI {
	logger.Debug("Creating new TodoUI")
	return ui.New(todoService)
}

// InitializeApp sets up the dependency injection
func InitializeApp() (*App, error) {
	wire.Build(DefaultSet)
	return nil, nil
}
