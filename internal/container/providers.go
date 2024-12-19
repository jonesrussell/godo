package container

import (
	"path/filepath"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

// provideSQLite creates a new SQLite store
func provideSQLite(cfg *config.Config, log logger.Logger) (*sqlite.Store, func(), error) {
	dbPath := filepath.Join(cfg.Database.Path, "godo.db")
	store, err := sqlite.New(dbPath, log)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := store.Close(); err != nil {
			log.Error("Failed to close database", "error", err)
		}
	}

	return store, cleanup, err
}

// provideLogger provides a basic logger for initial config loading
func provideLogger() (logger.Logger, error) {
	return logger.NewZapLogger(&logger.Config{
		Level:   "info",
		Console: true,
	})
}
