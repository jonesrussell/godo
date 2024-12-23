package config

import (
	"errors"
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

// Provider holds configuration paths and name
type Provider struct {
	paths      []string
	configName string
	configType string
}

// NewProvider creates a new Provider instance
func NewProvider(paths []string, name, configType string) *Provider {
	return &Provider{
		paths:      paths,
		configName: name,
		configType: configType,
	}
}

// NewDefaultConfig returns a new Config instance with default values
func NewDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "Godo",
			Version: "0.1.0",
			ID:      "io.github.jonesrussell.godo",
		},
		Logger: LoggerConfig{
			Level:   "info",
			Console: true,
		},
		Database: DatabaseConfig{
			Path: "godo.db",
		},
	}
}

// Load reads configuration from the specified file
func (p *Provider) Load() (*Config, error) {
	v := viper.New()

	v.SetConfigType(p.configType)

	for _, path := range p.paths {
		v.AddConfigPath(path)
	}

	v.SetConfigName(p.configName)
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

	cfg := NewDefaultConfig()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil
		}
		return nil, err
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ValidateConfig validates the configuration values
func ValidateConfig(cfg *Config) error {
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

// Add this function
func NewConfig(configPath string) (*Config, error) {
	provider := NewProvider(
		[]string{filepath.Dir(configPath)},
		filepath.Base(configPath),
		"yaml",
	)
	return provider.Load()
}
