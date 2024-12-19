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

func InitializeApp() (*app.App, func(), error) {
	wire.Build(
		provideLogger,
		config.Load,
		provideSQLite,
		app.NewApp,
		wire.Bind(new(storage.Store), new(*sqlite.Store)),
	)
	return nil, nil, nil
}
