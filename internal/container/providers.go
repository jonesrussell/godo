package container

import (
	"path/filepath"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
)

// ProvideSQLite creates a new SQLite store
func ProvideSQLite(cfg *config.Config, log logger.Logger) (*sqlite.Store, error) {
	dbPath := filepath.Join(cfg.Database.Path, "godo.db")
	return sqlite.New(dbPath, log)
}
