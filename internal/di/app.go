// app.go
package di

import (
    "context"
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/jonesrussell/godo/internal/hotkey"
    "github.com/jonesrussell/godo/internal/logger"
    "github.com/jonesrussell/godo/internal/service"
    "github.com/jonesrussell/godo/internal/ui"
)

// App represents the main application
type App struct {
    todoService   *service.TodoService
    hotkeyManager *hotkey.HotkeyManager
    program       *tea.Program
    ui           *ui.TodoUI
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
    if err := a.initializeServices(ctx); err != nil {
        return fmt.Errorf("failed to initialize services: %w", err)
    }

    return a.startServices(ctx)
}

func (a *App) startServices(ctx context.Context) error {
    // Start hotkey listener
    go a.startHotkeyListener(ctx)
    
    // Start UI program
    go a.startUIProgram()

    // Keep the main thread alive
    <-ctx.Done()
    return nil
}

func (a *App) startHotkeyListener(ctx context.Context) {
    logger.Info("Starting hotkey listener...")
    if err := a.hotkeyManager.Start(ctx); err != nil {
        logger.Error("Hotkey error: %v", err)
    }
}

func (a *App) startUIProgram() {
    logger.Info("Starting UI program...")
    if err := a.program.Start(); err != nil {
        logger.Error("UI program error: %v", err)
    }
}