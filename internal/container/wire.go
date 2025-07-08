//go:build wireinject && (windows || linux)

package container

import (
	"fmt"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/google/wire"

	"github.com/jonesrussell/godo/internal/api"
	"github.com/jonesrussell/godo/internal/app"
	apphotkey "github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/options"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

// Provider Sets
var (
	// CoreSet provides essential services that don't depend on UI or platform features
	CoreSet = wire.NewSet(
		BaseSet,
		LoggingSet,
		StorageSet,
		ServiceSet,
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
		wire.Bind(new(mainwindow.Interface), new(*mainwindow.Window)),
		wire.Bind(new(gui.QuickNote), new(*quicknote.Window)),
	)

	// APISet provides HTTP API server dependencies
	APISet = wire.NewSet(
		ProvideAPIConfig,
		ProvideAPIServer,
		ProvideAPIRunner,
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

	// ServiceSet provides business logic services
	ServiceSet = wire.NewSet(
		ProvideTaskService,
	)

	// ConfigSet provides configuration dependencies
	ConfigSet = wire.NewSet(
		ProvideConfig,
	)

	// HotkeySet provides hotkey dependencies
	HotkeySet = wire.NewSet(
		ProvideModifierKeys,
		ProvideKeyCode,
		ProvideHotkeyOptions,
		ProvideHotkeyManager,
	)

	// AppSet provides application dependencies
	AppSet = wire.NewSet(
		ProvideAppOptions,
		wire.Struct(new(app.Params), "*"),
		app.New,
		wire.Bind(new(app.Application), new(*app.App)),
	)
)

// InitializeApp initializes the application with all dependencies
func InitializeApp() (app.Application, func(), error) {
	wire.Build(
		CoreSet,   // First initialize core services
		UISet,     // Then UI components
		HotkeySet, // Then platform-specific features
		APISet,    // Then API server
		AppSet,    // Finally the main app
	)
	return nil, nil, nil
}

// Provider functions for options structs
func ProvideCoreOptions(
	logger logger.Logger,
	store storage.TaskStore,
	config *config.Config,
) (*options.CoreOptions, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if store == nil {
		return nil, fmt.Errorf("store is required")
	}
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	return &options.CoreOptions{
		Logger: logger,
		Store:  store,
		Config: config,
	}, nil
}

func ProvideGUIOptions(
	app fyne.App,
	mainWindow *mainwindow.Window,
	quickNote *quicknote.Window,
) (*options.GUIOptions, error) {
	if app == nil {
		return nil, fmt.Errorf("fyne app is required")
	}
	if mainWindow == nil {
		return nil, fmt.Errorf("main window is required")
	}
	if quickNote == nil {
		return nil, fmt.Errorf("quick note window is required")
	}
	return &options.GUIOptions{
		App:        app,
		MainWindow: mainWindow,
		QuickNote:  quickNote,
	}, nil
}

func ProvideLoggerOptions(
	level common.LogLevel,
	output common.LogOutputPaths,
	errorOutput common.ErrorOutputPaths,
) (*options.LoggerOptions, error) {
	if level == "" {
		return nil, fmt.Errorf("log level is required")
	}
	if len(output) == 0 {
		return nil, fmt.Errorf("log output paths are required")
	}
	if len(errorOutput) == 0 {
		return nil, fmt.Errorf("error output paths are required")
	}
	return &options.LoggerOptions{
		Level:       level,
		Output:      output,
		ErrorOutput: errorOutput,
	}, nil
}

func ProvideHotkeyOptions(
	modifiers common.ModifierKeys,
	key common.KeyCode,
) (*options.HotkeyOptions, error) {
	if len(modifiers) == 0 {
		return nil, fmt.Errorf("at least one modifier key is required")
	}
	if key == "" {
		return nil, fmt.Errorf("key code is required")
	}
	return &options.HotkeyOptions{
		Modifiers: modifiers,
		Key:       key,
	}, nil
}

func ProvideAppOptions(
	core *options.CoreOptions,
	gui *options.GUIOptions,
	name common.AppName,
	version common.AppVersion,
	id common.AppID,
) (*options.AppOptions, error) {
	if core == nil {
		return nil, fmt.Errorf("core options are required")
	}
	if gui == nil {
		return nil, fmt.Errorf("GUI options are required")
	}
	if name == "" {
		return nil, fmt.Errorf("app name is required")
	}
	if version == "" {
		return nil, fmt.Errorf("app version is required")
	}
	if id == "" {
		return nil, fmt.Errorf("app ID is required")
	}
	return &options.AppOptions{
		Core:    core,
		GUI:     gui,
		Name:    name,
		Version: version,
		ID:      id,
	}, nil
}

// Provider functions for common types
func ProvideAppName() (common.AppName, error) {
	name := common.AppName("Godo")
	if name == "" {
		return "", fmt.Errorf("app name cannot be empty")
	}
	return name, nil
}

func ProvideAppVersion() (common.AppVersion, error) {
	version := common.AppVersion("1.0.0")
	if version == "" {
		return "", fmt.Errorf("app version cannot be empty")
	}
	return version, nil
}

func ProvideAppID() (common.AppID, error) {
	id := common.AppID("com.jonesrussell.godo")
	if id == "" {
		return "", fmt.Errorf("app ID cannot be empty")
	}
	return id, nil
}

func ProvideDatabasePath() (common.DatabasePath, error) {
	path := common.DatabasePath("godo.db")
	if path == "" {
		return "", fmt.Errorf("database path cannot be empty")
	}
	return path, nil
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
func ProvideLogger(opts *options.LoggerOptions) (logger.Logger, func(), error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("logger options are required")
	}

	// Create a LogConfig that supports file logging
	logConfig := &common.LogConfig{
		Level:       string(opts.Level),
		Console:     true,
		File:        true,
		FilePath:    "logs/godo.log",
		Output:      opts.Output.Slice(),
		ErrorOutput: opts.ErrorOutput.Slice(),
	}

	log, err := logger.New(logConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Return a cleanup function that syncs the logger
	cleanup := func() {
		if zapLogger, ok := log.(*logger.ZapLogger); ok {
			if err := zapLogger.Sync(); err != nil {
				// Log sync errors but don't fail the application
				fmt.Printf("Failed to sync logger: %v\n", err)
			}
		}
	}

	return log, cleanup, nil
}

// ProvideSQLiteStore provides a SQLite store instance
func ProvideSQLiteStore(log logger.Logger, dbPath common.DatabasePath) (*sqlite.Store, func(), error) {
	if log == nil {
		return nil, nil, fmt.Errorf("logger is required")
	}
	if dbPath == "" {
		return nil, nil, fmt.Errorf("database path is required")
	}

	store, err := sqlite.New(string(dbPath), log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create store: %w", err)
	}

	cleanup := func() {
		if err := store.Close(); err != nil {
			log.Error("failed to close store during cleanup", "error", err)
		}
	}

	return store, cleanup, nil
}

// ProvideHotkeyManager provides a hotkey manager instance
func ProvideHotkeyManager(log logger.Logger, cfg *config.Config, quickNote *quicknote.Window) (apphotkey.Manager, error) {
	if log == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if quickNote == nil {
		return nil, fmt.Errorf("quick note window is required")
	}

	// Validate hotkey config
	if err := validateHotkeyConfig(&cfg.Hotkeys.QuickNote); err != nil {
		return nil, fmt.Errorf("invalid hotkey configuration: %w", err)
	}

	manager, err := apphotkey.New(quickNote, &cfg.Hotkeys.QuickNote, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create hotkey manager: %w", err)
	}

	return manager, nil
}

// validateHotkeyConfig validates the hotkey configuration
func validateHotkeyConfig(binding *common.HotkeyBinding) error {
	if binding == nil {
		return fmt.Errorf("hotkey binding is required")
	}
	if len(binding.Modifiers) == 0 {
		return fmt.Errorf("at least one modifier key is required")
	}
	if binding.Key == "" {
		return fmt.Errorf("key is required")
	}
	return nil
}

// ProvideFyneApp provides a Fyne application instance
func ProvideFyneApp(cfg *config.Config) fyne.App {
	fmt.Println("Creating Fyne application...")
	app := fyneapp.NewWithID(string(cfg.App.ID))
	app.Settings().SetTheme(theme.DefaultTheme())
	return app
}

// ProvideMainWindow provides a main window instance
func ProvideMainWindow(app fyne.App, store storage.TaskStore, logger logger.Logger, cfg *config.Config) *mainwindow.Window {
	return mainwindow.New(app, store, logger, cfg.UI.MainWindow)
}

// ProvideQuickNote provides a quick note window instance
func ProvideQuickNote(app fyne.App, store storage.TaskStore, logger logger.Logger, cfg *config.Config, mainWindow mainwindow.Interface) *quicknote.Window {
	return quicknote.New(app, store, logger, cfg.UI.QuickNote, mainWindow)
}

// ProvideModifierKeys provides the hotkey modifiers from config
func ProvideModifierKeys(cfg *config.Config) common.ModifierKeys {
	return common.ModifierKeys(cfg.Hotkeys.QuickNote.Modifiers)
}

// ProvideKeyCode provides the hotkey key code from config
func ProvideKeyCode(cfg *config.Config) common.KeyCode {
	return common.KeyCode(cfg.Hotkeys.QuickNote.Key)
}

// ProvideConfig provides the application configuration
func ProvideConfig() (*config.Config, error) {
	// Create a default configuration
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Godo",
			Version: "0.1.0",
			ID:      "io.github.jonesrussell.godo",
		},
		Logger: common.LogConfig{
			Level:       "info",
			Console:     true,
			File:        false,
			FilePath:    "",
			Output:      []string{"stdout"},
			ErrorOutput: []string{"stderr"},
		},
		Hotkeys: config.HotkeyConfig{
			QuickNote: common.HotkeyBinding{
				Modifiers: []string{"Ctrl", "Shift"},
				Key:       "N",
			},
		},
		Database: config.DatabaseConfig{
			Path: "godo.db",
		},
		UI: config.UIConfig{
			MainWindow: config.WindowConfig{
				Width:       800,
				Height:      600,
				StartHidden: false,
			},
			QuickNote: config.WindowConfig{
				Width:       400,
				Height:      300,
				StartHidden: true,
			},
		},
	}

	// Validate the configuration
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validateConfig performs validation of the entire configuration
func validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is required")
	}

	// Validate UI configuration
	if cfg.UI.MainWindow.Width <= 0 || cfg.UI.MainWindow.Height <= 0 {
		return fmt.Errorf("invalid window dimensions")
	}

	// Validate hotkey configuration
	if err := validateHotkeyConfig(&cfg.Hotkeys.QuickNote); err != nil {
		return fmt.Errorf("invalid hotkey configuration: %w", err)
	}

	return nil
}

// ProvideLogOutputPaths provides the default log output paths
func ProvideLogOutputPaths() common.LogOutputPaths {
	return common.LogOutputPaths{"stdout", "logs/godo.log"}
}

// ProvideErrorOutputPaths provides the default error output paths
func ProvideErrorOutputPaths() common.ErrorOutputPaths {
	return common.ErrorOutputPaths{"stderr", "logs/godo-error.log"}
}

// Provider functions for API components
func ProvideAPIConfig() *api.ServerConfig {
	return api.NewServerConfig()
}

func ProvideAPIServer(store storage.TaskStore, taskService service.TaskService, log logger.Logger) *api.Server {
	return api.NewServer(store, taskService, log)
}

func ProvideAPIRunner(store storage.TaskStore, taskService service.TaskService, log logger.Logger, cfg *api.ServerConfig) *api.Runner {
	return api.NewRunner(store, taskService, log, &common.HTTPConfig{
		Port:              8080, // TODO: Get from config
		ReadTimeout:       int(cfg.ReadTimeout.Seconds()),
		WriteTimeout:      int(cfg.WriteTimeout.Seconds()),
		ReadHeaderTimeout: int(cfg.ReadHeaderTimeout.Seconds()),
		IdleTimeout:       int(cfg.IdleTimeout.Seconds()),
	})
}

// ProvideTaskService provides a TaskService instance
func ProvideTaskService(store storage.TaskStore, log logger.Logger) service.TaskService {
	return service.NewTaskService(store, log)
}
