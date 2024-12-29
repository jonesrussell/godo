//go:build wireinject && windows

package container

import (
	"fmt"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/google/wire"
	"github.com/jonesrussell/godo/internal/app"
	apphotkey "github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/options"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
	"golang.design/x/hotkey"
)

// Provider Sets
var (
	// CoreSet provides essential services that don't depend on UI or platform features
	CoreSet = wire.NewSet(
		BaseSet,
		LoggingSet,
		StorageSet,
		ConfigSet,
		ProvideCoreOptions,
	)

	// UISet provides UI components after core services are initialized
	UISet = wire.NewSet(
		ProvideFyneApp,
		ProvideGUIOptions,
		ProvideMainWindow,
		ProvideQuickNote,
		wire.Bind(new(gui.MainWindow), new(*mainwindow.Window)),
		wire.Bind(new(gui.QuickNote), new(*quicknote.Window)),
	)

	// BaseSet provides basic application metadata
	BaseSet = wire.NewSet(
		ProvideAppName,
		ProvideAppVersion,
		ProvideAppID,
	)

	// LoggingSet provides logging dependencies
	LoggingSet = wire.NewSet(
		ProvideLogLevel,
		ProvideLogOutputPaths,
		ProvideErrorOutputPaths,
		ProvideLoggerOptions,
		ProvideLogger,
		wire.Bind(new(logger.Logger), new(*logger.ZapLogger)),
	)

	// StorageSet provides storage dependencies
	StorageSet = wire.NewSet(
		ProvideDatabasePath,
		ProvideSQLiteStore,
		wire.Bind(new(storage.TaskStore), new(*sqlite.Store)),
	)

	// ConfigSet provides configuration dependencies
	ConfigSet = wire.NewSet(
		ProvideConfig,
	)

	// HTTPSet provides HTTP server dependencies
	HTTPSet = wire.NewSet(
		ProvideHTTPPort,
		ProvideReadTimeout,
		ProvideWriteTimeout,
		ProvideHeaderTimeout,
		ProvideIdleTimeout,
		ProvideHTTPOptions,
		ProvideHTTPConfig,
	)

	// HotkeySet provides hotkey dependencies
	HotkeySet = wire.NewSet(
		ProvideModifierKeys,
		ProvideKeyCode,
		ProvideHotkeyOptions,
		ProvideHotkeyManager,
		wire.Bind(new(apphotkey.Manager), new(*apphotkey.DefaultManager)),
	)

	// AppSet provides application dependencies
	AppSet = wire.NewSet(
		ProvideAppOptions,
		wire.Struct(new(app.Params), "*"),
		app.New,
		wire.Bind(new(app.Application), new(*app.App)),
	)

	// TestSet provides mock dependencies for testing
	TestSet = wire.NewSet(
		BaseSet,
		LoggingSet,
		ProvideMockStore,
		ProvideMockMainWindow,
		ProvideMockQuickNote,
		ProvideMockHotkey,
		ProvideMockFyneApp,
		HTTPSet,
		wire.Bind(new(gui.MainWindow), new(*gui.MockMainWindow)),
		wire.Bind(new(gui.QuickNote), new(*gui.MockQuickNote)),
		wire.Bind(new(apphotkey.Manager), new(*apphotkey.MockManager)),
	)
)

// Provider functions for options structs
func ProvideCoreOptions(
	logger logger.Logger,
	store storage.TaskStore,
	config *config.Config,
) *options.CoreOptions {
	return &options.CoreOptions{
		Logger: logger,
		Store:  store,
		Config: config,
	}
}

func ProvideGUIOptions(
	app fyne.App,
	mainWindow *mainwindow.Window,
	quickNote *quicknote.Window,
) *options.GUIOptions {
	return &options.GUIOptions{
		App:        app,
		MainWindow: mainWindow,
		QuickNote:  quickNote,
	}
}

func ProvideLoggerOptions(
	level common.LogLevel,
	output common.LogOutputPaths,
	errorOutput common.ErrorOutputPaths,
) *options.LoggerOptions {
	return &options.LoggerOptions{
		Level:       level,
		Output:      output,
		ErrorOutput: errorOutput,
	}
}

func ProvideHotkeyOptions(
	modifiers common.ModifierKeys,
	key common.KeyCode,
) *options.HotkeyOptions {
	return &options.HotkeyOptions{
		Modifiers: modifiers,
		Key:       key,
	}
}

func ProvideHTTPOptions(
	port common.HTTPPort,
	readTimeout common.ReadTimeoutSeconds,
	writeTimeout common.WriteTimeoutSeconds,
	readHeaderTimeout common.HeaderTimeoutSeconds,
	idleTimeout common.IdleTimeoutSeconds,
) *options.HTTPOptions {
	return &options.HTTPOptions{
		Config: &common.HTTPConfig{
			Port:              port.Int(),
			ReadTimeout:       int(readTimeout),
			WriteTimeout:      int(writeTimeout),
			ReadHeaderTimeout: int(readHeaderTimeout),
			IdleTimeout:       int(idleTimeout),
		},
	}
}

func ProvideAppOptions(
	core *options.CoreOptions,
	gui *options.GUIOptions,
	http *options.HTTPOptions,
	hotkey *options.HotkeyOptions,
	name common.AppName,
	version common.AppVersion,
	id common.AppID,
) *options.AppOptions {
	return &options.AppOptions{
		Core:    core,
		GUI:     gui,
		HTTP:    http,
		Hotkey:  hotkey,
		Name:    name,
		Version: version,
		ID:      id,
	}
}

// InitializeApp initializes the application with all dependencies
func InitializeApp() (app.Application, func(), error) {
	wire.Build(
		CoreSet,   // First initialize core services
		UISet,     // Then UI components
		HotkeySet, // Then platform-specific features
		HTTPSet,   // Then HTTP server config
		AppSet,    // Finally the main app
	)
	return nil, nil, nil
}

// InitializeTestApp initializes the application with mock dependencies for testing
func InitializeTestApp() (*app.TestApp, func(), error) {
	wire.Build(
		TestSet,
		wire.Struct(new(app.TestApp), "*"),
	)
	return nil, nil, nil
}

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
func ProvideLogger(opts *options.LoggerOptions) (*logger.ZapLogger, func(), error) {
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
func ProvideSQLiteStore(log logger.Logger) (*sqlite.Store, func(), error) {
	dbPath := string(ProvideDatabasePath())
	store, err := sqlite.New(dbPath, log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create store: %w", err)
	}

	cleanup := func() {
		store.Close()
	}

	return store, cleanup, nil
}

// ProvideHotkeyManager provides a hotkey manager instance using options
func ProvideHotkeyManager(opts *options.HotkeyOptions) (*apphotkey.DefaultManager, error) {
	// Convert string modifiers to hotkey.Modifier
	modifiers := make([]hotkey.Modifier, 0, len(opts.Modifiers.Slice()))
	for _, mod := range opts.Modifiers.Slice() {
		switch mod {
		case "Ctrl":
			modifiers = append(modifiers, hotkey.ModCtrl)
		case "Alt":
			modifiers = append(modifiers, hotkey.ModAlt)
		case "Shift":
			modifiers = append(modifiers, hotkey.ModShift)
		}
	}

	// Convert key string to hotkey.Key
	var key hotkey.Key
	switch opts.Key.String() {
	case "N":
		key = hotkey.KeyN
	// Add other key mappings as needed
	default:
		return nil, fmt.Errorf("unsupported key: %s", opts.Key)
	}

	return apphotkey.NewManager(modifiers, key)
}

// ProvideFyneApp provides a Fyne application instance
func ProvideFyneApp() fyne.App {
	fmt.Println("Creating Fyne application...")
	app := fyneapp.New()
	app.Settings().SetTheme(theme.DefaultTheme())
	return app
}

// ProvideMainWindow provides a main window instance
func ProvideMainWindow(app fyne.App, store storage.TaskStore, logger logger.Logger, cfg *config.Config) *mainwindow.Window {
	return mainwindow.New(app, store, logger, cfg.UI.MainWindow)
}

// ProvideQuickNote provides a quick note window instance
func ProvideQuickNote(app fyne.App, store storage.TaskStore, logger logger.Logger, cfg *config.Config) *quicknote.Window {
	return quicknote.New(app, store, logger, cfg.UI.QuickNote)
}

// ProvideHTTPConfig provides HTTP configuration using options
func ProvideHTTPConfig(opts *options.HTTPOptions) *common.HTTPConfig {
	return &common.HTTPConfig{
		Port:              opts.Config.Port,
		ReadTimeout:       opts.Config.ReadTimeout,
		WriteTimeout:      opts.Config.WriteTimeout,
		ReadHeaderTimeout: opts.Config.ReadHeaderTimeout,
		IdleTimeout:       opts.Config.IdleTimeout,
	}
}

// Mock providers for testing
func ProvideMockStore() storage.TaskStore {
	return storage.NewMockStore()
}

// ProvideMockConfig provides a mock configuration for testing
func ProvideMockConfig() (*config.Config, error) {
	return config.NewDefaultConfig(), nil
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
func ProvideMockHotkey() *apphotkey.MockManager {
	return apphotkey.NewMockManager()
}

// ProvideMockFyneApp provides a mock Fyne app for testing
func ProvideMockFyneApp() fyne.App {
	return test.NewApp()
}

// ProvideLogOutputPaths provides the default log output paths
func ProvideLogOutputPaths() common.LogOutputPaths {
	return common.LogOutputPaths{"stdout"}
}

// ProvideErrorOutputPaths provides the default error output paths
func ProvideErrorOutputPaths() common.ErrorOutputPaths {
	return common.ErrorOutputPaths{"stderr"}
}

// ProvideModifierKeys provides the default hotkey modifiers
func ProvideModifierKeys() common.ModifierKeys {
	return common.ModifierKeys{"Ctrl", "Shift"}
}

// ProvideKeyCode provides the default hotkey key code
func ProvideKeyCode() common.KeyCode {
	return common.KeyCode("N")
}

// Provider functions for HTTP configuration
func ProvideHTTPPort() common.HTTPPort {
	return common.HTTPPort(8080)
}

func ProvideReadTimeout() common.ReadTimeoutSeconds {
	return common.ReadTimeoutSeconds(30)
}

func ProvideWriteTimeout() common.WriteTimeoutSeconds {
	return common.WriteTimeoutSeconds(30)
}

func ProvideHeaderTimeout() common.HeaderTimeoutSeconds {
	return common.HeaderTimeoutSeconds(10)
}

func ProvideIdleTimeout() common.IdleTimeoutSeconds {
	return common.IdleTimeoutSeconds(120)
}

// ProvideStoreAdapter provides a store adapter instance
func ProvideStoreAdapter(store storage.TaskStore) *storage.StoreAdapter {
	return storage.NewStoreAdapter(store)
}

// ProvideConfig provides the application configuration
func ProvideConfig() (*config.Config, error) {
	provider := config.NewProvider(
		[]string{".", "./configs"},
		"default",
		"yaml",
	)
	return provider.Load()
}
