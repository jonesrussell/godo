// providers.go
package app

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/quicknote"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
	"go.uber.org/zap"
)

// ProvideDefaultSet returns the Wire provider set for the application
func ProvideDefaultSet() wire.ProviderSet {
	return wire.NewSet(
		provideRepository,
		provideTodoService,
		provideTodoUI,
		provideQuickNoteUI,
		provideFyneApp,
	)
}

// provideRepository initializes and returns a SQLite-backed TodoRepository
func provideRepository(cfg *config.Config) (repository.TodoRepository, error) {
	logger.Debug("Creating SQLite repository", "path", cfg.Database.Path)
	db, err := database.NewSQLiteDB(cfg.Database.Path)
	if err != nil {
		return nil, err
	}
	return repository.NewSQLiteTodoRepository(db), nil
}

// provideTodoService creates a new TodoService instance
func provideTodoService(repo repository.TodoRepository) *service.TodoService {
	logger.Debug("Creating TodoService")
	return service.NewTodoService(repo)
}

// provideTodoUI creates a new TodoUI instance
func provideTodoUI(svc *service.TodoService, app fyne.App) *ui.TodoUI {
	logger.Debug("Creating TodoUI")
	window := app.NewWindow("Godo")
	return ui.NewTodoUI(*svc, window)
}

// provideQuickNoteUI creates a new QuickNote UI instance
func provideQuickNoteUI() (quicknote.UI, error) {
	logger.Debug("Creating QuickNoteUI")
	return quicknote.New()
}

// provideFyneApp creates a new Fyne application instance
func provideFyneApp() fyne.App {
	logger.Debug("Creating Fyne application")
	return fyneapp.New()
}

// provideLogger creates a new zap logger instance
func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	log, err := logger.InitializeWithConfig(cfg.Logging)
	if err != nil {
		return nil, err
	}
	return log, nil
}
