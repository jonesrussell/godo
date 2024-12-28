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

// Options structs for complex providers
type LoggerOptions struct {
	Level       common.LogLevel
	Output      common.OutputPaths
	ErrorOutput common.OutputPaths
}

type HTTPOptions struct {
	Port              common.HTTPPort
	ReadTimeout       common.TimeoutSeconds
	WriteTimeout      common.TimeoutSeconds
	ReadHeaderTimeout common.TimeoutSeconds
	IdleTimeout       common.TimeoutSeconds
}

type HotkeyOptions struct {
	Modifiers common.ModifierKeys
	Key       common.KeyCode
}

// Provider Sets
var (
	// BaseSet provides basic application metadata
	BaseSet = wire.NewSet(
		ProvideAppName,
		ProvideAppVersion,
		ProvideAppID,
	)

	// LoggingSet provides logging dependencies
	LoggingSet = wire.NewSet(
		wire.Struct(new(LoggerOptions), "*"),
		ProvideLogger,
		wire.Bind(new(logger.Logger), new(*logger.ZapLogger)),
	)

	// StorageSet provides storage dependencies
	StorageSet = wire.NewSet(
		ProvideDatabasePath,
		ProvideSQLiteStore,
		wire.Bind(new(storage.Store), new(*storage.SQLiteStore)),
	)

	// HTTPSet provides HTTP server dependencies
	HTTPSet = wire.NewSet(
		wire.Struct(new(HTTPOptions), "*"),
		ProvideHTTPConfig,
	)

	// HotkeySet provides hotkey dependencies
	HotkeySet = wire.NewSet(
		wire.Struct(new(HotkeyOptions), "*"),
		ProvideHotkeyManager,
		wire.Bind(new(hotkey.Manager), new(*hotkey.DefaultManager)),
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

	// TestSet provides mock dependencies for testing
	TestSet = wire.NewSet(
		ProvideMockStore,
		ProvideMockMainWindow,
		ProvideMockQuickNote,
		ProvideMockHotkey,
		wire.Bind(new(storage.Store), new(*storage.MockStore)),
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

// ProvideLogger provides a zap logger instance using options
func ProvideLogger(opts *LoggerOptions) (*logger.ZapLogger, func(), error) {
	config := &common.LogConfig{
		Level:       opts.Level.String(),
		Output:      opts.Output.Slice(),
		ErrorOutput: opts.ErrorOutput.Slice(),
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

// ProvideHotkeyManager provides a hotkey manager instance using options
func ProvideHotkeyManager(opts *HotkeyOptions) (*hotkey.DefaultManager, error) {
	binding := &common.HotkeyBinding{
		Modifiers: opts.Modifiers.Slice(),
		Key:       opts.Key.String(),
	}
	return hotkey.NewManager(binding, 0)
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

// ProvideHTTPConfig provides HTTP configuration using options
func ProvideHTTPConfig(opts *HTTPOptions) *common.HTTPConfig {
	return &common.HTTPConfig{
		Port:              opts.Port.Int(),
		ReadTimeout:       int(opts.ReadTimeout),
		WriteTimeout:      int(opts.WriteTimeout),
		ReadHeaderTimeout: int(opts.ReadHeaderTimeout),
		IdleTimeout:       int(opts.IdleTimeout),
	}
}

// InitializeApp initializes the application with all dependencies
func InitializeApp() (app.Application, func(), error) {
	wire.Build(
		BaseSet,
		LoggingSet,
		StorageSet,
		HTTPSet,
		HotkeySet,
		GUISet,
		AppSet,
	)
	return nil, nil, nil
}

// InitializeTestApp initializes the application with mock dependencies for testing
func InitializeTestApp() (*app.TestApp, func(), error) {
	wire.Build(
		LoggingSet,
		TestSet,
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
