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

// Application constants
const (
	DefaultAppName    = "Godo"
	DefaultAppVersion = "0.1.0"
	DefaultAppID      = "io.github.jonesrussell.godo"
	DefaultLogLevel   = "info"
	DefaultDBPath     = "godo.db"
	DefaultLogFile    = "logs/godo.log"
)

// Provider Sets - Organized by concern
var (
	// ConfigSet provides the single source of truth for configuration
	ConfigSet = wire.NewSet(
		ProvideConfig,
	)

	// CoreSet provides essential services
	CoreSet = wire.NewSet(
		ConfigSet,
		LoggingSet,
		StorageSet,
		ServiceSet,
		ProvideCoreOptions,
	)

	// LoggingSet provides logging infrastructure
	LoggingSet = wire.NewSet(
		ProvideLogger,
	)

	// StorageSet provides data persistence
	StorageSet = wire.NewSet(
		ProvideSQLiteStore,
		wire.Bind(new(storage.TaskStore), new(*sqlite.Store)),
	)

	// ServiceSet provides business logic
	ServiceSet = wire.NewSet(
		ProvideTaskService,
	)

	// UISet provides user interface components
	UISet = wire.NewSet(
		ProvideFyneApp,
		ProvideMainWindow,
		ProvideQuickNote,
		ProvideGUIOptions,
		wire.Bind(new(gui.MainWindow), new(*mainwindow.Window)),
		wire.Bind(new(mainwindow.Interface), new(*mainwindow.Window)),
		wire.Bind(new(gui.QuickNote), new(*quicknote.Window)),
	)

	// HotkeySet provides hotkey functionality
	HotkeySet = wire.NewSet(
		ProvideHotkeyManager,
		ProvideHotkeyOptions,
	)

	// APISet provides HTTP API server
	APISet = wire.NewSet(
		ProvideAPIServer,
		ProvideAPIRunner,
	)

	// AppSet provides the main application
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
		CoreSet,   // Configuration and core services
		UISet,     // User interface
		HotkeySet, // Platform-specific features
		APISet,    // API server
		AppSet,    // Main application
	)
	return nil, nil, nil
}

// Configuration provider - single source of truth
func ProvideConfig() (*config.Config, error) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    DefaultAppName,
			Version: DefaultAppVersion,
			ID:      DefaultAppID,
		},
		Logger: common.LogConfig{
			Level:       DefaultLogLevel,
			Console:     true,
			File:        true,
			FilePath:    DefaultLogFile,
			Output:      []string{"stdout", DefaultLogFile},
			ErrorOutput: []string{"stderr"},
		},
		Hotkeys: config.HotkeyConfig{
			QuickNote: common.HotkeyBinding{
				Modifiers: []string{"Ctrl", "Shift"},
				Key:       "G",
			},
		},
		Database: config.DatabaseConfig{
			Path: DefaultDBPath,
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

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Core options provider
func ProvideCoreOptions(
	logger logger.Logger,
	store storage.TaskStore,
	cfg *config.Config,
) (*options.CoreOptions, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if store == nil {
		return nil, fmt.Errorf("store is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	return &options.CoreOptions{
		Logger: logger,
		Store:  store,
		Config: cfg,
	}, nil
}

// GUI options provider
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

// App options provider
func ProvideAppOptions(
	core *options.CoreOptions,
	gui *options.GUIOptions,
) (*options.AppOptions, error) {
	if core == nil {
		return nil, fmt.Errorf("core options are required")
	}
	if gui == nil {
		return nil, fmt.Errorf("GUI options are required")
	}

	return &options.AppOptions{
		Core:    core,
		GUI:     gui,
		Name:    common.AppName(core.Config.App.Name),
		Version: common.AppVersion(core.Config.App.Version),
		ID:      common.AppID(core.Config.App.ID),
	}, nil
}

// Logger provider
func ProvideLogger(cfg *config.Config) (logger.Logger, func(), error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is required")
	}

	logConfig := &common.LogConfig{
		Level:       cfg.Logger.Level,
		Console:     cfg.Logger.Console,
		File:        cfg.Logger.File,
		FilePath:    cfg.Logger.FilePath,
		Output:      cfg.Logger.Output,
		ErrorOutput: cfg.Logger.ErrorOutput,
	}

	log, err := logger.New(logConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create logger: %w", err)
	}

	cleanup := func() {
		if zapLogger, ok := log.(*logger.ZapLogger); ok {
			if err := zapLogger.Sync(); err != nil {
				fmt.Printf("Failed to sync logger: %v\n", err)
			}
		}
	}

	return log, cleanup, nil
}

// SQLite store provider
func ProvideSQLiteStore(log logger.Logger, cfg *config.Config) (*sqlite.Store, func(), error) {
	if log == nil {
		return nil, nil, fmt.Errorf("logger is required")
	}
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is required")
	}

	store, err := sqlite.New(cfg.Database.Path, log)
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

// Task service provider
func ProvideTaskService(store storage.TaskStore, log logger.Logger) service.TaskService {
	return service.NewTaskService(store, log)
}

// Fyne app provider
func ProvideFyneApp(cfg *config.Config) fyne.App {
	app := fyneapp.NewWithID(cfg.App.ID)
	app.Settings().SetTheme(theme.DefaultTheme())
	return app
}

// Main window provider
func ProvideMainWindow(
	app fyne.App,
	store storage.TaskStore,
	log logger.Logger,
	cfg *config.Config,
) *mainwindow.Window {
	return mainwindow.New(app, store, log, cfg.UI.MainWindow)
}

// Quick note window provider
func ProvideQuickNote(
	app fyne.App,
	store storage.TaskStore,
	log logger.Logger,
	cfg *config.Config,
	mainWindow mainwindow.Interface,
) *quicknote.Window {
	return quicknote.New(app, store, log, cfg.UI.QuickNote, mainWindow)
}

// Hotkey manager provider
func ProvideHotkeyManager(
	log logger.Logger,
	cfg *config.Config,
	quickNote *quicknote.Window,
) (apphotkey.Manager, error) {
	if log == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if quickNote == nil {
		return nil, fmt.Errorf("quick note window is required")
	}

	if err := validateHotkeyConfig(&cfg.Hotkeys.QuickNote); err != nil {
		return nil, fmt.Errorf("invalid hotkey configuration: %w", err)
	}

	manager, err := apphotkey.New(quickNote, &cfg.Hotkeys.QuickNote, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create hotkey manager: %w", err)
	}

	return manager, nil
}

// Hotkey options provider
func ProvideHotkeyOptions(cfg *config.Config) (*options.HotkeyOptions, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	binding := &cfg.Hotkeys.QuickNote
	if len(binding.Modifiers) == 0 {
		return nil, fmt.Errorf("at least one modifier key is required")
	}
	if binding.Key == "" {
		return nil, fmt.Errorf("key is required")
	}

	return &options.HotkeyOptions{
		Modifiers: common.ModifierKeys(binding.Modifiers),
		Key:       common.KeyCode(binding.Key),
	}, nil
}

// API server provider
func ProvideAPIServer(
	store storage.TaskStore,
	taskService service.TaskService,
	log logger.Logger,
) *api.Server {
	return api.NewServer(store, taskService, log)
}

// API runner provider
func ProvideAPIRunner(
	store storage.TaskStore,
	taskService service.TaskService,
	log logger.Logger,
) *api.Runner {
	return api.NewRunner(store, taskService, log, &common.HTTPConfig{
		Port:              8080, // TODO: Get from config
		ReadTimeout:       30,
		WriteTimeout:      30,
		ReadHeaderTimeout: 10,
		IdleTimeout:       60,
	})
}

// Validation functions
func validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is required")
	}

	if cfg.UI.MainWindow.Width <= 0 || cfg.UI.MainWindow.Height <= 0 {
		return fmt.Errorf("invalid window dimensions")
	}

	if err := validateHotkeyConfig(&cfg.Hotkeys.QuickNote); err != nil {
		return fmt.Errorf("invalid hotkey configuration: %w", err)
	}

	return nil
}

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
