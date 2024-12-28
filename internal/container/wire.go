//go:build wireinject

package container

import (
	"fmt"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Provider Sets
var (
	// LoggingSet provides logging dependencies
	LoggingSet = wire.NewSet(
		ProvideLogger,
		wire.Bind(new(logger.Logger), new(*logger.ZapLogger)),
	)

	// ConfigSet provides configuration dependencies
	ConfigSet = wire.NewSet(
		ProvideAppName,
		ProvideAppVersion,
		ProvideAppID,
		ProvideDatabasePath,
		ProvideLogLevel,
		ProvideHotkeyBinding,
	)

	// StorageSet provides storage dependencies
	StorageSet = wire.NewSet(
		ProvideSQLiteStore,
		wire.Bind(new(storage.Store), new(*storage.SQLiteStore)),
	)

	// GUISet provides GUI dependencies
	GUISet = wire.NewSet(
		ProvideFyneApp,
		ProvideMainWindow,
		ProvideQuickNote,
		wire.Bind(new(gui.MainWindow), new(*mainwindow.Window)),
		wire.Bind(new(gui.QuickNote), new(*quicknote.Window)),
	)

	// AppSet provides application dependencies
	AppSet = wire.NewSet(
		app.New,
		wire.Bind(new(app.Application), new(*app.App)),
	)

	// HotkeySet provides hotkey dependencies
	HotkeySet = wire.NewSet(
		ProvideHotkeyManager,
		wire.Bind(new(hotkey.Manager), new(*hotkey.DefaultManager)),
	)

	// MockSet provides mock dependencies for testing
	MockSet = wire.NewSet(
		ProvideMockStore,
		ProvideMockMainWindow,
		ProvideMockQuickNote,
		ProvideMockHotkey,
		wire.Bind(new(gui.MainWindow), new(*gui.MockMainWindow)),
		wire.Bind(new(gui.QuickNote), new(*gui.MockQuickNote)),
		wire.Bind(new(hotkey.Manager), new(*hotkey.MockManager)),
	)
)

// Provider functions for common types
func ProvideAppName() common.AppName {
	return "Godo"
}

func ProvideAppVersion() common.AppVersion {
	return "1.0.0"
}

func ProvideAppID() common.AppID {
	return "com.jonesrussell.godo"
}

func ProvideDatabasePath() common.DatabasePath {
	return "godo.db"
}

func ProvideLogLevel() common.LogLevel {
	return "info"
}

// ProvideHotkeyBinding provides the default hotkey binding configuration
func ProvideHotkeyBinding() *common.HotkeyBinding {
	return &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "N",
	}
}

// ProvideLogger provides a zap logger instance
func ProvideLogger() (*logger.ZapLogger, func(), error) {
	config := &common.LogConfig{
		Level:       string(ProvideLogLevel()),
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	log, err := logger.NewZapLogger(config)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		// No cleanup needed for this logger implementation
	}

	return log, cleanup, nil
}

// ProvideSQLiteStore provides a SQLite store instance
func ProvideSQLiteStore(logger logger.Logger) (*storage.SQLiteStore, error) {
	return storage.NewSQLiteStore(string(ProvideDatabasePath()))
}

// ProvideHotkeyManager provides a hotkey manager instance
func ProvideHotkeyManager(binding *common.HotkeyBinding) (*hotkey.DefaultManager, error) {
	// Convert binding to hotkey format
	// Implementation needed
	return hotkey.NewManager(nil, 0) // Placeholder
}

// ProvideFyneApp provides a Fyne application instance
func ProvideFyneApp() fyne.App {
	fmt.Println("Creating Fyne application...")
	app := fyneapp.New()
	app.Settings().SetTheme(theme.DefaultTheme())
	return app
}

// ProvideMainWindow provides the main window instance
func ProvideMainWindow(store storage.Store, logger logger.Logger) *mainwindow.Window {
	return mainwindow.New(store, logger)
}

// ProvideQuickNote provides the quick note window instance
func ProvideQuickNote(store storage.Store, logger logger.Logger) *quicknote.Window {
	return quicknote.New(store, logger)
}

// InitializeApp initializes the application with all dependencies
func InitializeApp() (app.Application, func(), error) {
	wire.Build(
		ConfigSet,
		LoggingSet,
		StorageSet,
		GUISet,
		HotkeySet,
		AppSet,
	)
	return nil, nil, nil
}

// InitializeTestApp initializes the application with mock dependencies for testing
func InitializeTestApp() (*app.TestApp, func(), error) {
	wire.Build(
		LoggingSet,
		MockSet,
		wire.Struct(new(app.TestApp), "*"),
	)
	return nil, nil, nil
}

// Mock providers for testing
func ProvideMockStore() storage.Store {
	return storage.NewMockStore()
}

// ProvideMockMainWindow provides a mock main window for testing
func ProvideMockMainWindow() *gui.MockMainWindow {
	return &gui.MockMainWindow{}
}

// ProvideMockQuickNote provides a mock quick note window for testing
func ProvideMockQuickNote() *gui.MockQuickNote {
	return &gui.MockQuickNote{}
}

// ProvideMockHotkey provides a mock hotkey manager for testing
func ProvideMockHotkey() *hotkey.MockManager {
	return hotkey.NewMockManager()
}
