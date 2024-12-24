//go:build wireinject

package container

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/storage"
	"go.uber.org/zap"
)

// ProvideLogger provides a zap logger instance
func ProvideLogger() (*zap.Logger, func(), error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = logger.Sync()
	}
	return logger, cleanup, nil
}

// ProvideFyneApp provides a Fyne application instance
func ProvideFyneApp() fyne.App {
	return fyneapp.New()
}

// ProvideStorage provides a storage instance
func ProvideStorage() storage.Store {
	return storage.NewMemoryStore()
}

// ProvideHotkeyManager provides the platform-specific hotkey manager
func ProvideHotkeyManager() app.HotkeyManager {
	return app.NewHotkeyManager()
}

var Set = wire.NewSet(
	ProvideLogger,
	ProvideFyneApp,
	ProvideStorage,
	ProvideHotkeyManager,
	app.New,
)

// InitializeApp creates a new App instance with all dependencies
func InitializeApp() (*app.App, func(), error) {
	wire.Build(Set)
	return nil, nil, nil
}
