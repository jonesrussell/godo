// providers.go
package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/quicknote"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// DefaultSet defines the default provider set
var DefaultSet = wire.NewSet(
	provideTodoService,
	provideHotkeyManager,
	provideProgram,
	provideTodoUI,
	provideQuickNoteUI,
	provideRepository,
)

func provideTodoService(repo repository.TodoRepository) *service.TodoService {
	return service.NewTodoService(repo)
}

func provideRepository(cfg *config.Config) (repository.TodoRepository, error) {
	db, err := database.NewSQLiteDB(cfg.Database.Path)
	if err != nil {
		return nil, err
	}
	return repository.NewSQLiteTodoRepository(db), nil
}

func provideHotkeyManager() (hotkey.HotkeyManager, error) {
	return hotkey.NewHotkeyManager()
}

func provideProgram(todoUI *ui.TodoUI) *tea.Program {
	return tea.NewProgram(todoUI)
}

func provideTodoUI(svc *service.TodoService) *ui.TodoUI {
	return ui.New(svc)
}

func provideQuickNoteUI() (quicknote.UI, error) {
	return quicknote.New()
}
