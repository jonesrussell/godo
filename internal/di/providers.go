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

func NewApp(todoService *service.TodoService, ui *ui.TodoUI) (*App, error) {
    program := tea.NewProgram(ui)
    hotkeyManager := createHotkeyManager(program)
    
    return &App{
        todoService:   todoService,
        hotkeyManager: hotkeyManager,
        program:       program,
        ui:           ui,
    }, nil
}

func createHotkeyManager(program *tea.Program) *hotkey.HotkeyManager {
    showUI := func() {
        logger.Debug("Hotkey triggered - showing UI")
        program.Send(struct{}{})
    }
    return hotkey.New(showUI)
}

func provideTodoRepository(db *sql.DB) repository.TodoRepository {
    return repository.NewSQLiteTodoRepository(db)
}

func NewSQLiteDB() (*sql.DB, error) {
    return database.NewSQLiteDB("./godo.db")
}

func provideUI(todoService *service.TodoService) *ui.TodoUI {
    return ui.New(todoService)
}