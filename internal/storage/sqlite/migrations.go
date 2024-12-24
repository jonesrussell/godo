package sqlite

import (
	"database/sql"

	"github.com/jonesrussell/godo/internal/logger"
)

//nolint:gochecknoglobals // migrations need to be package-level for database initialization
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
func RunMigrations(db *sql.DB, log logger.Logger) error {
	log.Info("Running database migrations")

	for i, migration := range migrations {
		log.Debug("Executing migration",
			"index", i,
			"query", migration)

		if _, err := db.Exec(migration); err != nil {
			log.Error("Migration failed",
				"index", i,
				"error", err)
			return err
		}
	}

	log.Info("Database migrations completed successfully")
	return nil
}
