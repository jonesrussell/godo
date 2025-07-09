//go:build wireinject
// +build wireinject

// Package container provides dependency injection container setup
package container

import (
	"fmt"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/google/wire"

	"github.com/jonesrussell/godo/internal/application/app"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/repository"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/quicknote"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/sqlite"
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
		Logger: config.LogConfig{
			Level:   DefaultLogLevel,
			Console: true,
			File:    true,
			Output:  []string{"stdout"},
		},
		Hotkeys: config.HotkeyConfig{
			QuickNote: config.HotkeyBinding{
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
		HTTP: config.HTTPConfig{
			Port:              8080,
			ReadTimeout:       30,
			WriteTimeout:      30,
			ReadHeaderTimeout: 10,
			IdleTimeout:       120,
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

	logConfig := &logger.LogConfig{
		Level:   cfg.Logger.Level,
		Console: cfg.Logger.Console,
		File:    cfg.Logger.File,
		Output:  cfg.Logger.Output,
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

// Task service provider
func ProvideTaskService(repo repository.TaskRepository, log logger.Logger) service.TaskService {
	return service.NewTaskService(repo, log)
}

// Fyne app provider
func ProvideFyneApp(cfg *config.Config) fyne.App {
	app := fyneapp.New()
	app.SetIcon(theme.ComputerIcon())
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

// Quick note provider
func ProvideQuickNote(
	app fyne.App,
	store storage.TaskStore,
	log logger.Logger,
	cfg *config.Config,
) *quicknote.Window {
	return quicknote.New(app, store, log, cfg.UI.QuickNote)
}

// validateConfig validates the configuration
func validateConfig(cfg *config.Config) error {
	if cfg.App.Name == "" {
		return fmt.Errorf("app name is required")
	}
	if cfg.Logger.Level == "" {
		return fmt.Errorf("log level is required")
	}
	if err := validateHotkeyConfig(&cfg.Hotkeys.QuickNote); err != nil {
		return fmt.Errorf("invalid hotkey configuration: %w", err)
	}
	return nil
}

// validateHotkeyConfig validates hotkey configuration
func validateHotkeyConfig(binding *config.HotkeyBinding) error {
	if len(binding.Modifiers) == 0 {
		return fmt.Errorf("at least one modifier is required")
	}
	if binding.Key == "" {
		return fmt.Errorf("key is required")
	}
	for _, mod := range binding.Modifiers {
		switch mod {
		case "Ctrl", "Alt", "Shift":
			continue
		default:
			return fmt.Errorf("invalid modifier: %s", mod)
		}
	}
	return nil
}
