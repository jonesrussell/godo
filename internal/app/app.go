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

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	// Initialize quick note feature
	if err := a.quickNote.Show(ctx); err != nil {
		a.logger.Error("Failed to initialize quick note", zap.Error(err))
		return err
	}

	// Run the main UI
	return a.ui.Run()
}

// Cleanup performs any necessary cleanup before shutdown
func (a *App) Cleanup() error {
	a.logger.Info("Cleaning up application resources")
	// Add any cleanup logic here
	return nil
}
