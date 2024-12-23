package config

import (
	"path/filepath"

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

	// Set configuration type
	v.SetConfigType(p.configType)

	// Add config paths
	for _, path := range p.paths {
		v.AddConfigPath(path)
	}

	// Set config name
	v.SetConfigName(p.configName)

	// Set environment variable prefix
	v.SetEnvPrefix("GODO")
	v.AutomaticEnv()

	// Use the new function for default config
	cfg := NewDefaultConfig()

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		// If config file is not found, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil
		}
		return nil, err
	}

	// Unmarshal config
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
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
