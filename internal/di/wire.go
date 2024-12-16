// wire.go
//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
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
