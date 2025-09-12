//go:build wireinject
// +build wireinject

// Package container provides dependency injection container setup
package container

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"github.com/google/wire"
	"github.com/spf13/viper"

	"github.com/jonesrussell/godo/internal/application/core"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/repository"
	"github.com/jonesrussell/godo/internal/domain/service"
	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/quicknote"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/theme"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/factory"
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
		ProvideUnifiedStorage,
		ProvideNoteStoreAdapter,
		wire.Bind(new(storage.NoteStore), new(*storage.NoteStoreAdapter)),
	)

	// ServiceSet provides business logic
	ServiceSet = wire.NewSet(
		ProvideNoteRepositoryFromUnified,
		ProvideNoteService,
	)

	// UISet provides user interface components
	UISet = wire.NewSet(
		ProvideFyneApp,
		ProvideMainWindow,
		wire.Bind(new(gui.MainWindow), new(*mainwindow.Window)),
		wire.Bind(new(mainwindow.Interface), new(*mainwindow.Window)),
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
		core.New,
		wire.Bind(new(core.Application), new(*core.App)),
	)
)

// InitializeApp initializes the application with all dependencies
func InitializeApp() (core.Application, func(), error) {
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

	// Try to load config.yaml from executable directory and current directory
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add executable directory as first search path
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		v.AddConfigPath(execDir)
	}

	// Also try current directory as fallback
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
	}

	log, err := logger.New(logConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create logger: %w", err)
	}

	cleanup := func() {
		if zapLogger, ok := log.(*logger.ZapLogger); ok {
			// Sync is now safe since we use os.Stdout/os.Stderr directly
			if err := zapLogger.Sync(); err != nil {
				fmt.Printf("Failed to sync logger: %v\n", err)
			}
		}
	}

	return log, cleanup, nil
}

// Unified storage provider
func ProvideUnifiedStorage(log logger.Logger, cfg *config.Config) (domainstorage.UnifiedNoteStorage, func(), error) {
	if log == nil {
		return nil, nil, fmt.Errorf("logger is required")
	}
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is required")
	}

	// Convert config to storage config
	storageConfig := &domainstorage.StorageConfig{
		Type: domainstorage.StorageType(cfg.Storage.Type),
		SQLite: domainstorage.SQLiteConfig{
			FilePath: cfg.Storage.SQLite.FilePath,
		},
		API: domainstorage.APIConfig{
			BaseURL:    cfg.Storage.API.BaseURL,
			Timeout:    cfg.Storage.API.Timeout,
			RetryCount: cfg.Storage.API.RetryCount,
			RetryDelay: cfg.Storage.API.RetryDelay,
		},
	}

	// If storage config is not set, fall back to database config for backward compatibility
	if cfg.Storage.Type == "" {
		log.Debug("Storage type not configured, falling back to SQLite", "path", cfg.Database.Path)
		storageConfig.Type = domainstorage.StorageTypeSQLite
		storageConfig.SQLite.FilePath = cfg.Database.Path
	}

	log.Debug("Creating unified storage", "type", storageConfig.Type)
	store, err := factory.NewUnifiedStorage(storageConfig, log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create unified storage: %w", err)
	}

	cleanup := func() {
		if err := store.Close(); err != nil {
			log.Error("failed to close storage during cleanup", "error", err)
		}
	}

	return store, cleanup, nil
}

// Note store adapter provider
func ProvideNoteStoreAdapter(unifiedStore domainstorage.UnifiedNoteStorage) *storage.NoteStoreAdapter {
	return storage.NewNoteStoreAdapter(unifiedStore)
}

// Note repository provider from unified storage
func ProvideNoteRepositoryFromUnified(store domainstorage.UnifiedNoteStorage) repository.NoteRepository {
	return repository.NewNoteRepository(store)
}

// Note service provider
func ProvideNoteService(repo repository.NoteRepository, log logger.Logger) service.NoteService {
	return service.NewNoteService(repo, log)
}

// Fyne app provider
func ProvideFyneApp(cfg *config.Config) fyne.App {
	app := fyneapp.New()
	app.SetIcon(theme.AppIcon())
	return app
}

// Main window provider
func ProvideMainWindow(
	app fyne.App,
	store storage.NoteStore,
	log logger.Logger,
	cfg *config.Config,
) *mainwindow.Window {
	return mainwindow.New(app, store, log, cfg.UI.MainWindow)
}

// Quick note provider
func ProvideQuickNote(
	app fyne.App,
	store storage.NoteStore,
	log logger.Logger,
	cfg *config.Config,
) *quicknote.Window {
	return quicknote.New(app, store, log, cfg.UI.QuickNote)
}
