// providers.go
package app

import (
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// NewSQLiteDB creates a new database connection using configuration
func NewSQLiteDB(cfg *config.Config) (*sql.DB, error) {
	logger.Debug("Opening database at: %s", cfg.Database.Path)
	return database.NewSQLiteDB(cfg.Database.Path)
}

func provideTodoRepository(db *sql.DB) repository.TodoRepository {
	return repository.NewSQLiteTodoRepository(db)
}

func provideTodoService(repo repository.TodoRepository) *service.TodoService {
	return service.NewTodoService(repo)
}

func provideUI(todoService *service.TodoService) *ui.TodoUI {
	return ui.New(todoService)
}

func provideProgram(ui *ui.TodoUI) *tea.Program {
	return tea.NewProgram(ui)
}

func provideHotkeyManager() (hotkey.HotkeyManager, error) {
	logger.Debug("Initializing hotkey manager...")
	manager, err := hotkey.NewHotkeyManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create hotkey manager: %w", err)
	}
	return manager, nil
}

func provideApp(
	todoService *service.TodoService,
	hotkeyManager hotkey.HotkeyManager,
	program *tea.Program,
	ui *ui.TodoUI,
) *App {
	logger.Debug("Creating application instance...")
	return &App{
		todoService:   todoService,
		hotkeyManager: hotkeyManager,
		program:       program,
		ui:            ui,
	}
}
