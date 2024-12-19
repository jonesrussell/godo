package config

import (
	"os"
	"strings"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/spf13/viper"
	"golang.design/x/hotkey"
)

// Config holds the application configuration
type Config struct {
	App      AppConfig        `mapstructure:"app"`
	Database DatabaseConfig   `mapstructure:"database"`
	Logging  common.LogConfig `mapstructure:"logging"`
	Hotkeys  HotkeysConfig    `mapstructure:"hotkeys"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Path         string `mapstructure:"path"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// HotkeysConfig holds hotkey configurations
type HotkeysConfig struct {
	QuickNote HotkeyConfig `mapstructure:"quick_note"`
}

// HotkeyConfig holds configuration for a single hotkey
type HotkeyConfig struct {
	Modifiers []string `mapstructure:"modifiers"`
	Key       string   `mapstructure:"key"`
}

// String returns a string representation of the hotkey
func (h HotkeyConfig) String() string {
	return strings.Join(append(h.Modifiers, h.Key), "+")
}

// ToHotkey converts the config to a hotkey.Hotkey
func (h HotkeyConfig) ToHotkey() (*hotkey.Hotkey, error) {
	var mods []hotkey.Modifier
	for _, m := range h.Modifiers {
		switch strings.ToLower(m) {
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		}
	}

	var key hotkey.Key
	switch strings.ToUpper(h.Key) {
	case "G":
		key = hotkey.KeyG
		// Add more keys as needed
	}

	return hotkey.New(mods, key), nil
}

func Load(log logger.Logger) (*Config, error) {
	env := getEnv()
	log.Info("Loading configuration", "environment", env)

	v := viper.New()

	// Set default values first
	v.SetDefault("database.max_open_conns", 1)
	v.SetDefault("database.max_idle_conns", 1)

	// Configure Viper
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.AddConfigPath(".")
	v.SetEnvPrefix("GODO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind environment variables
	v.BindEnv("database.path", "GODO_DATABASE_PATH")
	v.BindEnv("database.max_open_conns", "GODO_DATABASE_MAX_OPEN_CONNS")
	v.BindEnv("database.max_idle_conns", "GODO_DATABASE_MAX_IDLE_CONNS")
	v.BindEnv("logging.level", "GODO_LOGGING_LEVEL")

	// Load default config
	v.SetConfigName("default")
	if err := v.ReadInConfig(); err != nil {
		log.Error("Failed to read default config", "error", err)
		return nil, err
	}

	// Load environment specific config
	if env != "development" {
		v.SetConfigName(env)
		if err := v.MergeInConfig(); err != nil {
			log.Error("Failed to merge environment config", "error", err)
			return nil, err
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Error("Failed to unmarshal config", "error", err)
		return nil, err
	}

	// Ensure default values are set if not provided
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 1
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 1
	}

	return &config, nil
}

func getEnv() string {
	if env := os.Getenv("GODO_ENV"); env != "" {
		return env
	}
	return "development"
}
