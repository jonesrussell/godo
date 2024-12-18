// app.go
package app

import (
	"context"

	"github.com/jonesrussell/godo/internal/quicknote"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
	"go.uber.org/zap"
)

// App represents the main application
type App struct {
	todoService service.TodoServicer
	quickNote   quicknote.UI
	ui          *ui.TodoUI
	logger      *zap.Logger
}

// NewApp creates a new App instance with all dependencies
func NewApp(
	todoService service.TodoServicer,
	quickNote quicknote.UI,
	todoUI *ui.TodoUI,
	logger *zap.Logger,
) *App {
	return &App{
		todoService: todoService,
		quickNote:   quickNote,
		ui:          todoUI,
		logger:      logger,
	}
}

// GetTodoService returns the todo service instance
func (a *App) GetTodoService() service.TodoServicer {
	return a.todoService
}

// GetQuickNote returns the quick note UI instance
func (a *App) GetQuickNote() quicknote.UI {
	return a.quickNote
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	// Start the UI in the main goroutine
	return a.ui.Run()
}

// Cleanup performs any necessary cleanup before shutdown
func (a *App) Cleanup() error {
	a.logger.Info("Cleaning up application resources")
	// Add any cleanup logic here
	return nil
}
