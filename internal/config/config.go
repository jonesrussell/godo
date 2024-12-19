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
	if strings.EqualFold(h.Key, "CTRL+SPACE") {
		key = hotkey.KeyG
		// Add more keys as needed
	}

	return hotkey.New(mods, key), nil
}

// Load loads the configuration from files and environment
func Load(log logger.Logger) (*Config, error) {
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

	// Use loadHotkeys
	hotkeys, err := loadHotkeys(v)
	if err != nil {
		return nil, err
	}

	if err := bindEnvVariables(v, log); err != nil {
		return nil, err
	}

	// Load default config
	v.SetConfigName("default")
	if err := v.ReadInConfig(); err != nil {
		log.Error("Failed to read default config", "error", err)
		return nil, err
	}

	// Load environment specific config
	if env := getEnv(); env != "development" {
		v.SetConfigName(env)
		if err := v.MergeInConfig(); err != nil {
			log.Error("Failed to merge environment config", "error", err)
			return nil, err
		}
	}

	return &Config{
		App: AppConfig{
			Name:    v.GetString("app.name"),
			Version: v.GetString("app.version"),
		},
		Database: DatabaseConfig{
			Path:         v.GetString("database.path"),
			MaxOpenConns: v.GetInt("database.max_open_conns"),
			MaxIdleConns: v.GetInt("database.max_idle_conns"),
		},
		Hotkeys: HotkeysConfig{
			QuickNote: *hotkeys,
		},
		Logging: common.LogConfig{
			Level: v.GetString("logging.level"),
		},
	}, nil
}

func getEnv() string {
	if env := os.Getenv("GODO_ENV"); env != "" {
		return env
	}
	return "development"
}

// loadHotkeys loads hotkey configuration
func loadHotkeys(v *viper.Viper) (*HotkeyConfig, error) {
	h := &HotkeyConfig{}
	if err := v.UnmarshalKey("hotkeys.quick_note", h); err != nil {
		return nil, err
	}

	if strings.EqualFold(h.Key, "CTRL+SPACE") {
		h.Key = "CTRL+SPACE"
	}

	return h, nil
}

// bindEnvVariables binds environment variables to configuration
func bindEnvVariables(v *viper.Viper, log logger.Logger) error {
	envVars := []struct {
		key string
		env string
	}{
		{"database.path", "GODO_DATABASE_PATH"},
		{"database.max_open_conns", "GODO_DATABASE_MAX_OPEN_CONNS"},
		{"database.max_idle_conns", "GODO_DATABASE_MAX_IDLE_CONNS"},
		{"logging.level", "GODO_LOG_LEVEL"},
	}

	for _, ev := range envVars {
		if err := v.BindEnv(ev.key, ev.env); err != nil {
			log.Error("Failed to bind environment variable",
				"key", ev.key,
				"env", ev.env,
				"error", err)
			return err
		}
	}

	return nil
}
