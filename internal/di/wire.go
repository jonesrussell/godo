// wire.go
//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/config"
)

// DefaultSet defines the provider set for wire
var DefaultSet = wire.NewSet(
	NewSQLiteDB,
	provideTodoRepository,
	provideTodoService,
	provideUI,
	provideProgram,
	provideHotkeyManager,
	provideApp,
)

// InitializeApp sets up the dependency injection
func InitializeApp() (*App, error) {
	wire.Build(DefaultSet)
	return &App{}, nil
}

func InitializeAppWithConfig(cfg *config.Config) (*App, error) {
	wire.Build(DefaultSet)
	return &App{}, nil
}
