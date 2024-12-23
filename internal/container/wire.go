//go:build wireinject
// +build wireinject

package container

import (
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

// Wire sets for dependency injection
var defaultSet = wire.NewSet(
	provideLogger,
	config.Load,
	provideSQLite,
	app.NewApp,
	wire.Bind(new(storage.Store), new(*sqlite.Store)),
)

// InitializeApp creates a new application instance with all dependencies wired
func InitializeApp() (*app.App, func(), error) {
	wire.Build(defaultSet)
	return nil, nil, nil
}
