//go:build wireinject
// +build wireinject

package container

import (
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

// BootstrapLogger provides a basic logger for initial config loading
func BootstrapLogger() (logger.Logger, error) {
	return logger.NewZapLogger(&logger.Config{
		Level:   "info",
		Console: true,
	})
}

var defaultSet = wire.NewSet(
	BootstrapLogger,
	config.Load,
	ProvideSQLite,
	app.NewApp,
	wire.Bind(new(storage.Store), new(*sqlite.Store)),
)

func InitializeApp() (*app.App, func(), error) {
	wire.Build(defaultSet)
	return nil, nil, nil
}
