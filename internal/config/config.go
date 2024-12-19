package config

import (
	"os"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	yaml "gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	App      AppConfig        `yaml:"app"`
	Database DatabaseConfig   `yaml:"database"`
	Hotkeys  HotkeyConfig     `yaml:"hotkeys"`
	Logging  common.LogConfig `yaml:"logging"`
	UI       UIConfig         `yaml:"ui"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type DatabaseConfig struct {
	Path         string `yaml:"path"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type HotkeyConfig struct {
	QuickNote *common.HotkeyBinding `yaml:"quick_note"`
	OpenApp   *common.HotkeyBinding `yaml:"open_app"`
}

type UIConfig struct {
	QuickNote QuickNoteConfig `yaml:"quick_note"`
}

type QuickNoteConfig struct {
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
	Title  string `yaml:"title"`
}

// Load loads the configuration from files and environment
func Load(log logger.Logger) (*Config, error) {
	env := "development"
	config := &Config{}

	// Load default config
	if err := loadConfigFile(config, "configs/default.yaml"); err != nil {
		log.Error("Failed loading default config", "error", err)
		return nil, err
	}

	// Load environment-specific config if it exists
	envConfig := "configs/" + env + ".yaml"
	if _, err := os.Stat(envConfig); err == nil {
		if err := loadConfigFile(config, envConfig); err != nil {
			log.Error("Failed loading environment config",
				"env", env,
				"error", err)
			return nil, err
		}
	}

	// Override with environment variables
	if err := loadEnvOverrides(config); err != nil {
		log.Error("Failed loading environment overrides", "error", err)
		return nil, err
	}

	return config, nil
}

func loadConfigFile(config *Config, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(config)
}

// loadEnvOverrides loads configuration overrides from environment variables
func loadEnvOverrides(config *Config) error {
	// Example: GODO_DATABASE_PATH overrides database.path
	if path := os.Getenv("GODO_DATABASE_PATH"); path != "" {
		config.Database.Path = path
	}

	if level := os.Getenv("GODO_LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}

	return nil
}
