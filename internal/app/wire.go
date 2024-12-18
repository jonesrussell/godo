// wire.go
//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/config"
)

// DefaultSet defines the provider set for wire without config
var DefaultSet = wire.NewSet(
	provideTodoRepository,
	provideTodoService,
	provideUI,
	provideProgram,
	provideHotkeyManager,
	provideApp,
)

// ConfiguredSet defines the provider set that requires configuration
var ConfiguredSet = wire.NewSet(
	DefaultSet,
	NewSQLiteDB,
)

// InitializeApp sets up the dependency injection
//
// Deprecated: Use InitializeAppWithConfig instead. This function will be removed in a future version.
// It uses hardcoded defaults and doesn't support proper configuration management.
// Migration guide:
//  1. Load configuration using config.Load()
//  2. Pass the config to InitializeAppWithConfig()
//
// Example:
//
//	cfg, err := config.Load("development")
//	if err != nil {
//	    // handle error
//	}
//	app, err := di.InitializeAppWithConfig(cfg)
func InitializeApp() (*App, error) {
	// For backwards compatibility, use default config
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: "./godo.db",
		},
	}
	return InitializeAppWithConfig(cfg)
}

// InitializeAppWithConfig sets up the dependency injection with configuration
func InitializeAppWithConfig(cfg *config.Config) (*App, error) {
	wire.Build(ConfiguredSet)
	return &App{}, nil
}
