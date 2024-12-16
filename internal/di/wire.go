//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/hotkey"
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
	log.Println("Initializing services...")
	if err := a.initializeServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	log.Println("Services initialized successfully")

	// Start hotkey listener in a goroutine
	log.Println("Starting hotkey listener...")
	go func() {
		if err := a.hotkeyManager.Start(ctx); err != nil {
			log.Printf("Hotkey error: %v\n", err)
		}
	}()
	log.Println("Hotkey listener started")

	// Run the Bubble Tea program
	log.Println("Starting UI...")
	if err := a.program.Start(); err != nil {
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
func NewApp(todoService *service.TodoService, ui *ui.TodoUI) *App {
	program := tea.NewProgram(ui)

	// Create a closure that captures the UI instance
	showUI := func() {
		program.Send(struct{}{}) // We'll handle the actual message type in the UI
	}

	hotkeyManager := hotkey.New(showUI)

	return &App{
		todoService:   todoService,
		ui:            ui,
		program:       program,
		hotkeyManager: hotkeyManager,
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
	provideUI,
	NewApp,
)

// Add database provider
func NewSQLiteDB() (*sql.DB, error) {
	return database.NewSQLiteDB("./godo.db")
}

func provideUI(todoService *service.TodoService) *ui.TodoUI {
	return ui.New(todoService)
}

// InitializeApp sets up the dependency injection
func InitializeApp() (*App, error) {
	wire.Build(DefaultSet)
	return nil, nil
}
