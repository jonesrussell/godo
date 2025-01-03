// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"database/sql"
)

// RunMigrations applies all database migrations
func RunMigrations(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);
	`
	_, err := db.Exec(query)
	return err
}
