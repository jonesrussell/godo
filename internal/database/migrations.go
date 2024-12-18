package database

import (
	"database/sql"

	"github.com/jonesrussell/godo/internal/logger"
)

// Migrations holds all database migrations
var migrations = []string{
	`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`,
}

// RunMigrations executes all database migrations
func RunMigrations(db *sql.DB) error {
	logger.Info("Running database migrations")

	for i, migration := range migrations {
		logger.Debug("Executing migration",
			"index", i,
			"query", migration)

		if _, err := db.Exec(migration); err != nil {
			logger.Error("Migration failed",
				"index", i,
				"error", err)
			return err
		}
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
