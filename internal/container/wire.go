//go:build wireinject

package container

import (
	"fmt"

	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

// Remove the build constraint from the generated file
//go:generate wire

// Wire sets for dependency injection
var defaultSet = wire.NewSet(
	provideLogger,
	provideConfig,
	provideSQLite,
	provideMainWindow,
	app.NewApp,
	wire.Bind(new(storage.Store), new(*sqlite.Store)),
)

// provideConfig creates a new config provider and loads the configuration
func provideConfig() (*config.Config, error) {
	provider := config.NewProvider(
		[]string{
			".",              // Current directory
			"./configs",      // Project configs directory
			"~/.config/godo", // User config directory
		},
		"default",
		"yaml",
	)
	return provider.Load()
}

// provideLogger provides a basic logger for initial config loading
func provideLogger() (logger.Logger, error) {
	return logger.New(&common.LogConfig{
		Level:       "info",
		Console:     true,
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	})
}

// provideSQLite creates a new SQLite store
func provideSQLite(cfg *config.Config, log logger.Logger) (*sqlite.Store, func(), error) {
	if cfg.Database.Path == "" {
		return nil, nil, fmt.Errorf("invalid database path: path cannot be empty")
	}

	store, err := sqlite.New(cfg.Database.Path, log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SQLite store: %w", err)
	}

	cleanup := func() {
		if err := store.Close(); err != nil {
			log.Error("Failed to close SQLite store", "error", err)
		}
	}

	return store, cleanup, nil
}

// provideMainWindow creates the main application window
func provideMainWindow(store storage.Store, log logger.Logger) *mainwindow.Window {
	return mainwindow.New(store, log)
}

// InitializeApp creates a new application instance with all dependencies wired
func InitializeApp() (*app.App, func(), error) {
	wire.Build(defaultSet)
	return nil, nil, nil
}
