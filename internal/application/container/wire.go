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
	"github.com/spf13/viper"

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

// Configuration provider - uses actual config system
func ProvideConfig() (*config.Config, error) {
	// Start with default config
	cfg := config.NewDefaultConfig()

	// Try to load config.yaml from current directory
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err == nil {
		// Config file found, unmarshal into our config
		if configErr := v.Unmarshal(cfg); configErr != nil {
			// If unmarshaling fails, keep defaults
			fmt.Printf("Failed to parse config file, using defaults: %v\n", configErr)
		} else {
			fmt.Printf("Config file loaded: %s\n", v.ConfigFileUsed())
		}
	} else {
		// No config file found, using defaults
		fmt.Printf("No config file found, using defaults: %v\n", err)
	}

	return cfg, nil
}

// Logger provider
func ProvideLogger(cfg *config.Config) (logger.Logger, func(), error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is required")
	}

	logConfig := &logger.LogConfig{
		Level:    cfg.Logger.Level,
		Console:  cfg.Logger.Console,
		File:     cfg.Logger.File,
		FilePath: cfg.Logger.FilePath,
		Output:   cfg.Logger.Output,
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
