// providers.go
package di

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

func provideTodoRepository(db *sql.DB) repository.TodoRepository {
	return repository.NewSQLiteTodoRepository(db)
}

func NewSQLiteDB() (*sql.DB, error) {
	dbPath := "./godo.db"
	logger.Debug("Opening database at: %s", dbPath)
	return database.NewSQLiteDB(dbPath)
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

func provideHotkeyManager() (*hotkey.HotkeyManager, error) {
	logger.Debug("Initializing hotkey manager...")
	manager := hotkey.NewHotkeyManager()
	return manager, nil
}

func provideApp(
	todoService *service.TodoService,
	hotkeyManager *hotkey.HotkeyManager,
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
