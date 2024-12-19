package container

import (
	"os"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

func provideSQLite(cfg *config.Config, log logger.Logger) (*sqlite.Store, error) {
	return sqlite.New(cfg.Database.Path, log)
}

func provideLogger() (logger.Logger, error) {
	defaultConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}
	return logger.New(defaultConfig)
}

func provideEnvironment(log logger.Logger) string {
	env := os.Getenv("GODO_ENV")
	if env == "" {
		env = "development"
	}
	log.Debug("Using environment", "env", env)
	return env
}
