// Package config handles application configuration management
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Configuration keys and defaults
const (
	// Environment settings
	EnvPrefix = "GODO"

	// Config paths
	DefaultConfigDir = "godo"

	// Default values
	DefaultAppID = "io.github.jonesrussell.godo"

	// Config keys
	KeyAppName    = "app.name"
	KeyAppVersion = "app.version"
	KeyAppID      = "app.id"
	KeyDBPath     = "database.path"
	KeyLogLevel   = "logger.level"
	KeyLogConsole = "logger.console"

	// DefaultMainWindowWidth and other default window dimensions
	DefaultMainWindowWidth  = 800
	DefaultMainWindowHeight = 600
	DefaultQuickNoteWidth   = 200
	DefaultQuickNoteHeight  = 100
)

// LogConfig holds logging configuration
type LogConfig struct {
	Level    string `mapstructure:"level"`
	Console  bool   `mapstructure:"console"`
	File     bool   `mapstructure:"file"`
	FilePath string `mapstructure:"file_path"`
}

// Config holds all application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Logger   LogConfig      `mapstructure:"logger"`
	Hotkeys  HotkeyConfig   `mapstructure:"hotkeys"`
	Database DatabaseConfig `mapstructure:"database"` // Deprecated - use Storage
	Storage  StorageConfig  `mapstructure:"storage"`
	UI       UIConfig       `mapstructure:"ui"`
	HTTP     HTTPConfig     `mapstructure:"http"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name             string `mapstructure:"name"`
	Version          string `mapstructure:"version"`
	ID               string `mapstructure:"id"`
	ForceKillTimeout int    `mapstructure:"force_kill_timeout"`
}

// HotkeyConfig holds hotkey configuration
type HotkeyConfig struct {
	QuickNote    HotkeyBinding `mapstructure:"quick_note"`
	MainWindow   HotkeyBinding `mapstructure:"main_window"`
	RetryDelayMs int           `mapstructure:"retry_delay_ms"`
	MaxRetries   int           `mapstructure:"max_retries"`
}

// HotkeyBinding represents a hotkey combination
type HotkeyBinding struct {
	Modifiers []string `mapstructure:"modifiers"`
	Key       string   `mapstructure:"key"`
}

// DatabaseConfig holds database configuration (deprecated - use StorageConfig)
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// StorageConfig holds storage configuration for different backends
type StorageConfig struct {
	Type   string       `mapstructure:"type"`
	SQLite SQLiteConfig `mapstructure:"sqlite"`
	API    APIConfig    `mapstructure:"api"`
}

// SQLiteConfig holds SQLite-specific configuration
type SQLiteConfig struct {
	FilePath string `mapstructure:"file_path"`
}

// APIConfig holds API-specific configuration
type APIConfig struct {
	BaseURL            string `mapstructure:"base_url"`
	Timeout            int    `mapstructure:"timeout_seconds"`
	RetryCount         int    `mapstructure:"retry_count"`
	RetryDelay         int    `mapstructure:"retry_delay_ms"`
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
}

// UIConfig holds UI-related configuration
type UIConfig struct {
	MainWindow WindowConfig `mapstructure:"main_window"`
	QuickNote  WindowConfig `mapstructure:"quick_note"`
}

// WindowConfig holds window-specific configuration
type WindowConfig struct {
	Width       int  `mapstructure:"width"`
	Height      int  `mapstructure:"height"`
	StartHidden bool `mapstructure:"start_hidden"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port              int `mapstructure:"port"`
	ReadTimeout       int `mapstructure:"read_timeout"`
	WriteTimeout      int `mapstructure:"write_timeout"`
	ReadHeaderTimeout int `mapstructure:"read_header_timeout"`
	IdleTimeout       int `mapstructure:"idle_timeout"`
	StartupTimeout    int `mapstructure:"startup_timeout"`
	ShutdownTimeout   int `mapstructure:"shutdown_timeout"`
}

// Logger interface for configuration
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	WithError(err error) Logger
}

// Provider handles configuration loading and validation
type Provider struct {
	paths      []string
	configName string
	configType string
	log        Logger
}

// Load reads and validates configuration
func (p *Provider) Load() (*Config, error) {
	v := viper.New()
	p.log.Info("starting config load")

	// Set up environment variables
	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	cfg := NewDefaultConfig()
	p.setDefaults(v, cfg)

	// Configure and read config file
	p.configureConfigFile(v)

	if err := v.ReadInConfig(); err != nil {
		p.log.Warn("config file read error", "error", err)
	} else {
		p.log.Info("config file loaded", "file", v.ConfigFileUsed())
	}

	// Bind environment variables explicitly
	if err := p.bindEnvironmentVariables(v); err != nil {
		return nil, err
	}

	p.log.Debug("after env binding",
		"app.name", v.GetString(KeyAppName),
		"database.path", v.GetString(KeyDBPath),
		"hotkeys.quick_note.modifiers", v.GetStringSlice("hotkeys.quick_note.modifiers"),
		"hotkeys.quick_note.key", v.GetString("hotkeys.quick_note.key"))

	// Unmarshal into struct
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		p.log.WithError(err).Error("unmarshal error")
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	p.log.Debug("after unmarshal",
		"app.name", cfg.App.Name,
		"database.path", cfg.Database.Path,
		"hotkeys.quick_note.modifiers", cfg.Hotkeys.QuickNote.Modifiers,
		"hotkeys.quick_note.key", cfg.Hotkeys.QuickNote.Key)

	// Validate and resolve paths
	if err := ValidateConfig(cfg); err != nil {
		p.log.WithError(err).Error("validation error")
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	p.log.Debug("validation passed")

	if err := p.ResolvePaths(cfg); err != nil {
		p.log.WithError(err).Error("path resolution error")
		return nil, fmt.Errorf("failed to resolve paths: %w", err)
	}

	p.log.Info("config load complete",
		"app.name", cfg.App.Name,
		"database.path", cfg.Database.Path)

	return cfg, nil
}

// setDefaults sets default values for configuration
func (p *Provider) setDefaults(v *viper.Viper, cfg *Config) {
	v.SetDefault(KeyAppName, cfg.App.Name)
	v.SetDefault(KeyAppVersion, cfg.App.Version)
	v.SetDefault(KeyAppID, cfg.App.ID)
	v.SetDefault(KeyDBPath, cfg.Database.Path)
	v.SetDefault(KeyLogLevel, cfg.Logger.Level)
	v.SetDefault(KeyLogConsole, cfg.Logger.Console)
	v.SetDefault("hotkeys.quick_note.modifiers", cfg.Hotkeys.QuickNote.Modifiers)
	v.SetDefault("hotkeys.quick_note.key", cfg.Hotkeys.QuickNote.Key)
	v.SetDefault("hotkeys.main_window.modifiers", cfg.Hotkeys.MainWindow.Modifiers)
	v.SetDefault("hotkeys.main_window.key", cfg.Hotkeys.MainWindow.Key)
	v.SetDefault("hotkeys.quick_note.retry_delay_ms", cfg.Hotkeys.RetryDelayMs)
	v.SetDefault("hotkeys.quick_note.max_retries", cfg.Hotkeys.MaxRetries)
	v.SetDefault("app.force_kill_timeout", cfg.App.ForceKillTimeout)
	v.SetDefault("http.startup_timeout", cfg.HTTP.StartupTimeout)
	v.SetDefault("http.shutdown_timeout", cfg.HTTP.ShutdownTimeout)
}

// configureConfigFile sets up the config file configuration
func (p *Provider) configureConfigFile(v *viper.Viper) {
	v.SetConfigType(p.configType)
	v.SetConfigName(p.configName)
	for _, path := range p.paths {
		v.AddConfigPath(path)
		p.log.Debug("added config path", "path", path)
	}
}

// bindEnvironmentVariables binds environment variables to configuration keys
func (p *Provider) bindEnvironmentVariables(v *viper.Viper) error {
	envBindings := map[string]string{
		KeyAppName:                          EnvPrefix + "_APP_NAME",
		KeyAppVersion:                       EnvPrefix + "_APP_VERSION",
		KeyAppID:                            EnvPrefix + "_APP_ID",
		KeyDBPath:                           EnvPrefix + "_DATABASE_PATH",
		KeyLogLevel:                         EnvPrefix + "_LOGGER_LEVEL",
		KeyLogConsole:                       EnvPrefix + "_LOGGER_CONSOLE",
		"hotkeys.quick_note.modifiers":      EnvPrefix + "_HOTKEYS_QUICK_NOTE_MODIFIERS",
		"hotkeys.quick_note.key":            EnvPrefix + "_HOTKEYS_QUICK_NOTE_KEY",
		"hotkeys.main_window.modifiers":     EnvPrefix + "_HOTKEYS_MAIN_WINDOW_MODIFIERS",
		"hotkeys.main_window.key":           EnvPrefix + "_HOTKEYS_MAIN_WINDOW_KEY",
		"hotkeys.quick_note.retry_delay_ms": EnvPrefix + "_HOTKEYS_QUICK_NOTE_RETRY_DELAY_MS",
		"hotkeys.quick_note.max_retries":    EnvPrefix + "_HOTKEYS_QUICK_NOTE_MAX_RETRIES",
		"app.force_kill_timeout":            EnvPrefix + "_APP_FORCE_KILL_TIMEOUT",
		"http.startup_timeout":              EnvPrefix + "_HTTP_STARTUP_TIMEOUT",
		"http.shutdown_timeout":             EnvPrefix + "_HTTP_SHUTDOWN_TIMEOUT",
	}

	for k, env := range envBindings {
		if err := v.BindEnv(k, env); err != nil {
			return err
		}
		if envVal := os.Getenv(env); envVal != "" {
			p.log.Debug("environment variable found", "key", env, "value", envVal)
			// Special handling for modifiers array
			if k == "hotkeys.quick_note.modifiers" || k == "hotkeys.main_window.modifiers" {
				p.log.Debug("processing modifiers array", "raw_value", envVal)
				// Remove brackets and quotes, then split by comma
				trimmed := strings.Trim(envVal, "[]")
				p.log.Debug("after trim brackets", "value", trimmed)
				parts := strings.Split(trimmed, ",")
				p.log.Debug("after split", "parts", parts)
				modifiers := make([]string, len(parts))
				for i, part := range parts {
					modifiers[i] = strings.Trim(part, "\" ")
					p.log.Debug("processed modifier", "index", i, "raw", part, "cleaned", modifiers[i])
				}
				p.log.Debug("setting modifiers", "final_value", modifiers)
				v.Set(k, modifiers)
			}
		}
	}
	return nil
}

// ResolvePaths resolves relative paths in the config to absolute paths
func (p *Provider) ResolvePaths(cfg *Config) error {
	// Skip path resolution for tests or when explicitly set to relative
	if os.Getenv("GODO_TEST_MODE") == "true" {
		return nil
	}

	// Expand environment variables in paths
	originalPath := cfg.Database.Path
	cfg.Database.Path = os.ExpandEnv(cfg.Database.Path)
	cfg.Logger.FilePath = os.ExpandEnv(cfg.Logger.FilePath)

	p.log.Debug("Database path resolution",
		"original", originalPath,
		"expanded", cfg.Database.Path,
		"is_absolute", filepath.IsAbs(cfg.Database.Path))

	if !filepath.IsAbs(cfg.Database.Path) {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		cfg.Database.Path = filepath.Join(userConfigDir, "godo", cfg.Database.Path)
		p.log.Debug("Made database path absolute", "final_path", cfg.Database.Path)
	}
	return nil
}

// ValidateConfig validates the configuration values
func ValidateConfig(cfg *Config) error {
	var validationErrors []string

	if cfg.App.Name == "" {
		validationErrors = append(validationErrors, "app name is required")
	}

	if !isValidLogLevel(cfg.Logger.Level) {
		validationErrors = append(validationErrors, "invalid log level: "+cfg.Logger.Level)
	}

	if cfg.HTTP.StartupTimeout <= 0 {
		validationErrors = append(validationErrors, "http.startup_timeout must be positive")
	}
	if cfg.HTTP.ShutdownTimeout <= 0 {
		validationErrors = append(validationErrors, "http.shutdown_timeout must be positive")
	}

	if len(validationErrors) > 0 {
		return &Error{
			Op:  "validate",
			Err: fmt.Errorf("validation failed: %s", strings.Join(validationErrors, "; ")),
		}
	}

	return nil
}

func isValidLogLevel(level string) bool {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	return validLevels[strings.ToLower(level)]
}

// NewDefaultConfig creates a new configuration with default values
func NewDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:             "Godo",
			Version:          "0.1.0",
			ID:               "io.github.jonesrussell/godo",
			ForceKillTimeout: 2,
		},
		Logger: LogConfig{
			Level:    "info",
			Console:  true,
			File:     true,
			FilePath: "logs/godo.log",
		},
		Database: DatabaseConfig{
			Path: "godo.db",
		},
		Hotkeys: HotkeyConfig{
			QuickNote: HotkeyBinding{
				Modifiers: []string{}, // Let config file or environment set this
				Key:       "",         // Let config file or environment set this
			},
			MainWindow: HotkeyBinding{
				Modifiers: []string{}, // Let config file or environment set this
				Key:       "",         // Let config file or environment set this
			},
			RetryDelayMs: 100,
			MaxRetries:   3,
		},
		UI: UIConfig{
			MainWindow: WindowConfig{
				Width:       DefaultMainWindowWidth,
				Height:      DefaultMainWindowHeight,
				StartHidden: false,
			},
			QuickNote: WindowConfig{
				Width:       DefaultQuickNoteWidth,
				Height:      DefaultQuickNoteHeight,
				StartHidden: false,
			},
		},
		HTTP: HTTPConfig{
			Port:              8080,
			ReadTimeout:       30,
			WriteTimeout:      30,
			ReadHeaderTimeout: 10,
			IdleTimeout:       120,
			StartupTimeout:    5,
			ShutdownTimeout:   5,
		},
	}
}

// Error represents a configuration operation error
type Error struct {
	Op  string
	Err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("config %s: %v", e.Op, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
