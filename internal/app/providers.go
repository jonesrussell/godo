// providers.go
package app

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// DefaultSet defines the provider set for wire without config
var DefaultSet = wire.NewSet(
	provideTodoRepository,
	provideTodoService,
	provideUI,
	provideProgram,
	provideHotkeyManager,
	provideQuickNoteUI,
	provideApp,
)

// ConfiguredSet defines the provider set that requires configuration
var ConfiguredSet = wire.NewSet(
	DefaultSet,
	NewSQLiteDB,
)

// Provider functions
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
	return hotkey.NewHotkeyManager()
}

func provideQuickNoteUI() (ui.QuickNoteUI, error) {
	return ui.NewQuickNoteUI()
}

func provideApp(
	todoService *service.TodoService,
	hotkeyManager hotkey.HotkeyManager,
	program *tea.Program,
	ui *ui.TodoUI,
	quickNote ui.QuickNoteUI,
) *App {
	return NewApp(todoService, hotkeyManager, program, ui, quickNote)
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(cfg *config.Config) (*sql.DB, error) {
	logger.Debug("Opening database at: %s", cfg.Database.Path)
	return database.NewSQLiteDB(cfg.Database.Path)
}
