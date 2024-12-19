package container

import (
	"os"

	"github.com/jonesrussell/godo/internal/logger"
)

func provideEnvironment() string {
	env := os.Getenv("GODO_ENV")
	if env == "" {
		env = "development"
	}
	logger.Debug("Using environment", "env", env)
	return env
}
