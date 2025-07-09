//go:build wireinject && (windows || linux)

package container

import (
	"fmt"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/google/wire"

	"github.com/jonesrussell/godo/internal/application/app"
	"github.com/jonesrussell/godo/internal/domain/repository"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/quicknote"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/sqlite"
	"github.com/jonesrussell/godo/internal/shared/common"
	"github.com/jonesrussell/godo/internal/shared/config"
)

// Application constants
const (
	DefaultAppName    = "Godo"
	DefaultAppVersion = "0.1.0"
	DefaultAppID      = "io.github.jonesrussell.godo"
	DefaultLogLevel   = "info"
	DefaultDBPath     = "godo.db"
	DefaultLogFile    = "logs/godo.log"
)

// Provider Sets - Organized by concern
var (
	// ConfigSet provides the single source of truth for configuration
	ConfigSet = wire.NewSet(
		ProvideConfig,
	)

	// LoggingSet provides logging infrastructure
	LoggingSet = wire.NewSet(
		ProvideLogger,
	)

	// StorageSet provides data persistence
	StorageSet = wire.NewSet(
		ProvideSQLiteStore,
		wire.Bind(new(storage.TaskStore), new(*sqlite.Store)),
	)

	// ServiceSet provides business logic
	ServiceSet = wire.NewSet(
		ProvideTaskRepository,
		ProvideTaskService,
	)

	// UISet provides user interface components
	UISet = wire.NewSet(
		ProvideFyneApp,
		ProvideMainWindow,
		ProvideQuickNote,
		wire.Bind(new(gui.MainWindow), new(*mainwindow.Window)),
		wire.Bind(new(mainwindow.Interface), new(*mainwindow.Window)),
		wire.Bind(new(gui.QuickNote), new(*quicknote.Window)),
	)

	// CoreSet provides essential services
	CoreSet = wire.NewSet(
		ConfigSet,
		LoggingSet,
		StorageSet,
		ServiceSet,
	)

	// AppSet provides the main application
	AppSet = wire.NewSet(
		app.New,
		wire.Bind(new(app.Application), new(*app.App)),
	)
)

// InitializeApp initializes the application with all dependencies
func InitializeApp() (app.Application, func(), error) {
	wire.Build(
		CoreSet, // Configuration and core services
		UISet,   // User interface
		AppSet,  // Main application
	)
	return nil, nil, nil
}

// Configuration provider - single source of truth
func ProvideConfig() (*config.Config, error) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    DefaultAppName,
			Version: DefaultAppVersion,
			ID:      DefaultAppID,
		},
		Logger: common.LogConfig{
			Level:       DefaultLogLevel,
			Console:     true,
			File:        true,
			FilePath:    DefaultLogFile,
			Output:      []string{"stdout", DefaultLogFile},
			ErrorOutput: []string{"stderr"},
		},
		Hotkeys: config.HotkeyConfig{
			QuickNote: common.HotkeyBinding{
				Modifiers: []string{"Ctrl", "Shift"},
				Key:       "G",
			},
		},
		Database: config.DatabaseConfig{
			Path: DefaultDBPath,
		},
		UI: config.UIConfig{
			MainWindow: config.WindowConfig{
				Width:       800,
				Height:      600,
				StartHidden: false,
			},
			QuickNote: config.WindowConfig{
				Width:       400,
				Height:      300,
				StartHidden: true,
			},
		},
		HTTP: common.HTTPConfig{
			Port: 8080,
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Logger provider
func ProvideLogger(cfg *config.Config) (logger.Logger, func(), error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is required")
	}

	logConfig := &common.LogConfig{
		Level:       cfg.Logger.Level,
		Console:     cfg.Logger.Console,
		File:        cfg.Logger.File,
		FilePath:    cfg.Logger.FilePath,
		Output:      cfg.Logger.Output,
		ErrorOutput: cfg.Logger.ErrorOutput,
	}

	log, err := logger.New(logConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create logger: %w", err)
	}

	cleanup := func() {
		if zapLogger, ok := log.(*logger.ZapLogger); ok {
			if err := zapLogger.Sync(); err != nil {
				fmt.Printf("Failed to sync logger: %v\n", err)
			}
		}
	}

	return log, cleanup, nil
}

// SQLite store provider
func ProvideSQLiteStore(log logger.Logger, cfg *config.Config) (*sqlite.Store, func(), error) {
	if log == nil {
		return nil, nil, fmt.Errorf("logger is required")
	}
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is required")
	}

	store, err := sqlite.New(cfg.Database.Path, log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create store: %w", err)
	}

	cleanup := func() {
		if err := store.Close(); err != nil {
			log.Error("failed to close store during cleanup", "error", err)
		}
	}

	return store, cleanup, nil
}

// Task repository provider
func ProvideTaskRepository(store storage.TaskStore) repository.TaskRepository {
	return repository.NewTaskRepository(store)
}

// Task service provider (now uses repository)
func ProvideTaskService(repo repository.TaskRepository, log logger.Logger) service.TaskService {
	return service.NewTaskService(repo, log)
}

// Fyne app provider
func ProvideFyneApp(cfg *config.Config) fyne.App {
	app := fyneapp.NewWithID(cfg.App.ID)
	app.Settings().SetTheme(theme.DefaultTheme())
	return app
}

// Main window provider
func ProvideMainWindow(
	app fyne.App,
	store storage.TaskStore,
	log logger.Logger,
	cfg *config.Config,
) *mainwindow.Window {
	return mainwindow.New(app, store, log, cfg.UI.MainWindow)
}

// Quick note window provider
func ProvideQuickNote(
	app fyne.App,
	store storage.TaskStore,
	log logger.Logger,
	cfg *config.Config,
) *quicknote.Window {
	return quicknote.New(app, store, log, cfg.UI.QuickNote)
}

// Validation functions
func validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is required")
	}

	if cfg.UI.MainWindow.Width <= 0 || cfg.UI.MainWindow.Height <= 0 {
		return fmt.Errorf("invalid window dimensions")
	}

	if err := validateHotkeyConfig(&cfg.Hotkeys.QuickNote); err != nil {
		return fmt.Errorf("invalid hotkey configuration: %w", err)
	}

	return nil
}

func validateHotkeyConfig(binding *common.HotkeyBinding) error {
	if binding == nil {
		return fmt.Errorf("hotkey binding is required")
	}
	if len(binding.Modifiers) == 0 {
		return fmt.Errorf("at least one modifier key is required")
	}
	if binding.Key == "" {
		return fmt.Errorf("key is required")
	}
	return nil
}
