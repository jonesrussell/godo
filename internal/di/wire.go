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
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// App represents the main application
type App struct {
	todoService *service.TodoService
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	if err := a.initializeServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	ui := ui.New(a.todoService)
	p := tea.NewProgram(ui)

	// Run UI in a goroutine
	go func() {
		<-ctx.Done()
		p.Quit()
	}()

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start UI: %w", err)
	}

	return nil
}

func (a *App) initializeServices(ctx context.Context) error {
	// Verify database connection
	// We can do this by creating a test todo
	testTodo, err := a.todoService.CreateTodo(ctx, "Test Todo", "Testing service initialization")
	if err != nil {
		return fmt.Errorf("failed to verify database connection: %w", err)
	}

	// Clean up test todo
	if err := a.todoService.DeleteTodo(ctx, testTodo.ID); err != nil {
		return fmt.Errorf("failed to cleanup test todo: %w", err)
	}

	return nil
}

// NewApp creates a new App instance
func NewApp(todoService *service.TodoService) *App {
	return &App{
		todoService: todoService,
	}
}

// Provide TodoRepository interface implementation
func provideTodoRepository(db *sql.DB) repository.TodoRepository {
	return repository.NewSQLiteTodoRepository(db)
}

// Add provider set
var DefaultSet = wire.NewSet(
	NewSQLiteDB,
	provideTodoRepository,
	service.NewTodoService,
	NewApp,
)

// Add database provider
func NewSQLiteDB() (*sql.DB, error) {
	return database.NewSQLiteDB("./godo.db")
}

// InitializeApp sets up the dependency injection
func InitializeApp() (*App, error) {
	wire.Build(DefaultSet)
	return nil, nil
}
