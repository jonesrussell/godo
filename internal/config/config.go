package config

import (
	"os"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Hotkeys  HotkeyConfig   `mapstructure:"hotkeys"`
	Database DatabaseConfig `mapstructure:"database"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	ID      string `mapstructure:"id"`
}

// LoggerConfig holds logger-specific configuration
type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Console  bool   `mapstructure:"console"`
	File     bool   `mapstructure:"file"`
	FilePath string `mapstructure:"file_path"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// Load reads configuration from the specified file
func (c *Config) Load(configPath string) error {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return v.Unmarshal(c)
}

// NewConfig creates a new Config instance with defaults
func NewConfig(log logger.Logger, configPath string) (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:    "Godo",
			Version: "0.1.0",
			ID:      "io.github.jonesrussell.godo",
		},
		Logger: LoggerConfig{
			Level:    "info",
			Console:  true,
			File:     false,
			FilePath: "",
		},
		Hotkeys: HotkeyConfig{
			QuickNote: "Ctrl+Alt+G",
		},
		Database: DatabaseConfig{
			Path: "./data", // Default database path
		},
	}

	// Try to load config file
	if err := cfg.Load(configPath); err != nil {
		// If file not found, use defaults and log info
		if os.IsNotExist(err) {
			log.Info("config file not found, using defaults", "path", configPath)
			return cfg, nil
		}
		// Other errors should be returned
		return nil, err
	}

	return cfg, nil
}
