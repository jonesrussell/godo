// wire.go
//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/service"
)

// ConfiguredSet is the provider set for the application
var ConfiguredSet = wire.NewSet(
	// Core dependencies
	provideLogger,
	provideRepository,
	provideFyneApp,

	// Service layer
	provideTodoService,
	wire.Bind(new(service.TodoServicer), new(*service.TodoService)),

	// UI components
	provideQuickNoteUI,
	provideTodoUI,

	// Application
	NewApp,
)

// InitializeAppWithConfig sets up the dependency injection with configuration
func InitializeAppWithConfig(cfg *config.Config) (*App, error) {
	wire.Build(ConfiguredSet)
	return &App{}, nil
}
