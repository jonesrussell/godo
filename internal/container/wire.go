//go:build wireinject

package container

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// ProvideLogger provides a zap logger instance
func ProvideLogger() (logger.Logger, func(), error) {
	config := &common.LogConfig{
		Level:       "info",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	log, err := logger.New(config)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		// No cleanup needed for this logger implementation
	}

	return log, cleanup, nil
}

// ProvideFyneApp provides a Fyne application instance
func ProvideFyneApp() fyne.App {
	fmt.Println("Creating Fyne application...")

	// Set any required environment variables
	os.Setenv("FYNE_SCALE", "1.0") // Force 1.0 scale for now

	// Create the app
	app := fyneapp.NewWithID("com.jonesrussell.godo")
	if app == nil {
		fmt.Println("ERROR: Failed to create Fyne application")
		return nil
	}

	// Set theme and other settings
	app.Settings().SetTheme(theme.DefaultTheme())

	fmt.Println("Fyne application created successfully")
	return app
}

// ProvideStorage provides a storage instance
func ProvideStorage() storage.Store {
	return storage.NewMemoryStore()
}

// ProvideQuickNote provides a quick note window instance
func ProvideQuickNote(store storage.Store, logger logger.Logger) quicknote.Interface {
	return quicknote.New(store, logger)
}

// ProvideHotkeyManager provides the platform-specific hotkey manager
func ProvideHotkeyManager(quickNote quicknote.Interface) hotkey.Manager {
	return hotkey.New(quickNote)
}

var Set = wire.NewSet(
	ProvideLogger,
	ProvideFyneApp,
	ProvideStorage,
	ProvideQuickNote,
	ProvideHotkeyManager,
	wire.Bind(new(hotkey.QuickNoteService), new(quicknote.Interface)),
	app.New,
)

// InitializeApp creates a new App instance with all dependencies
func InitializeApp() (*app.App, func(), error) {
	wire.Build(Set)
	return nil, nil, nil
}
