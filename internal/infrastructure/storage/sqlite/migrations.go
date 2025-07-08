// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"database/sql"
)

// RunMigrations applies all database migrations
func RunMigrations(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	`
	_, err := db.Exec(query)
	return err
}
