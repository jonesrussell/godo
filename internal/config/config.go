package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

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

// Provider handles configuration loading and validation
type Provider struct {
	paths      []string
	configName string
	configType string
}

// NewProvider creates a new configuration provider
func NewProvider(paths []string, configName, configType string) *Provider {
	return &Provider{
		paths:      paths,
		configName: configName,
		configType: configType,
	}
}

// Load reads and validates configuration from files and environment
func (p *Provider) Load() (*Config, error) {
	v := viper.New()

	// Set up Viper
	v.SetConfigType(p.configType)
	for _, path := range p.paths {
		v.AddConfigPath(path)
	}
	v.SetConfigName(p.configName)

	// Environment variables
	v.SetEnvPrefix("GODO")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind specific environment variables
	if err := v.BindEnv("database.path", "GODO_DATABASE_PATH"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("logger.level", "GODO_LOGGER_LEVEL"); err != nil {
		return nil, err
	}

	// Load defaults first
	cfg := NewDefaultConfig()

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Return error only if it's not a missing file
			return nil, err
		}
		// Missing config file is ok, we'll use defaults
	}

	// Unmarshal the config
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}

	// Resolve paths
	if err := p.resolvePaths(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// resolvePaths resolves relative paths in the config to absolute paths
func (p *Provider) resolvePaths(cfg *Config) error {
	// Skip path resolution for tests or when explicitly set to relative
	if os.Getenv("GODO_TEST_MODE") == "true" {
		return nil
	}

	if !filepath.IsAbs(cfg.Database.Path) {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		cfg.Database.Path = filepath.Join(userConfigDir, "godo", cfg.Database.Path)
	}
	return nil
}

// ValidateConfig validates the configuration values
func ValidateConfig(cfg *Config) error {
	if cfg.App.Name == "" {
		return errors.New("app name is required")
	}

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[strings.ToLower(cfg.Logger.Level)] {
		return errors.New("invalid log level: " + cfg.Logger.Level)
	}

	return nil
}

// NewDefaultConfig creates a new configuration with default values
func NewDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "Godo",
			Version: "0.1.0",
			ID:      "io.github.jonesrussell.godo",
		},
		Database: DatabaseConfig{
			Path: "godo.db",
		},
		Logger: LoggerConfig{
			Level:   "info",
			Console: true,
		},
		Hotkeys: HotkeyConfig{
			QuickNote: "Ctrl+Alt+G",
		},
	}
}
