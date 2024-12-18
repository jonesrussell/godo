// providers.go
package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/quicknote"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// DefaultSet defines the default provider set for wire dependency injection.
// It includes all necessary providers for creating a complete application instance.
var DefaultSet = wire.NewSet(
	provideTodoService,
	provideProgram,
	provideTodoUI,
	provideQuickNoteUI,
	provideRepository,
)

// provideTodoService creates a new TodoService instance with the given repository.
// It handles the business logic for todo operations.
func provideTodoService(repo repository.TodoRepository) *service.TodoService {
	logger.Debug("Creating TodoService")
	return service.NewTodoService(repo)
}

// provideRepository initializes and returns a SQLite-backed TodoRepository.
// It handles the database connection and repository setup.
func provideRepository(cfg *config.Config) (repository.TodoRepository, error) {
	logger.Debug("Creating SQLite repository", "path", cfg.Database.Path)
	db, err := database.NewSQLiteDB(cfg.Database.Path)
	if err != nil {
		return nil, err
	}
	return repository.NewSQLiteTodoRepository(db), nil
}

// provideProgram creates a new Bubble Tea program instance with the TodoUI.
// It sets up the terminal user interface program.
func provideProgram(todoUI *ui.TodoUI) *tea.Program {
	logger.Debug("Creating Bubble Tea program")
	return tea.NewProgram(todoUI)
}

// provideTodoUI creates a new TodoUI instance with the given service.
// It handles the terminal-based todo management interface.
func provideTodoUI(svc *service.TodoService) *ui.TodoUI {
	logger.Debug("Creating TodoUI")
	return ui.New(svc)
}

// provideQuickNoteUI creates a new platform-specific QuickNote UI instance.
// It handles the popup window for quick note capture.
func provideQuickNoteUI() (quicknote.UI, error) {
	logger.Debug("Creating QuickNoteUI")
	return quicknote.New()
}
