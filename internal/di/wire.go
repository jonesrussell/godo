//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"database/sql"

	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
)

// App represents the main application
type App struct {
	todoService *service.TodoService
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	// TODO: Implement application logic
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
