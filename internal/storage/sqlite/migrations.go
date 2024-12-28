// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"database/sql"
)

// migrationSet holds database migrations
type migrationSet struct {
	migrations []string
}

// newMigrationSet creates a new migration set with default migrations
func newMigrationSet() *migrationSet {
	return &migrationSet{
		migrations: []string{
			// Initial schema - create tasks table if it doesn't exist
			`DROP TABLE IF EXISTS tasks;
			CREATE TABLE tasks (
				id TEXT PRIMARY KEY,
				content TEXT NOT NULL,
				done BOOLEAN NOT NULL DEFAULT 0,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL
			);`,
		},
	}
}

// RunMigrations executes all migrations in the set
func (ms *migrationSet) RunMigrations(db *sql.DB) error {
	for _, migration := range ms.migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}
	return nil
}

// RunMigrations executes all database migrations
func RunMigrations(db *sql.DB) error {
	ms := newMigrationSet()
	return ms.RunMigrations(db)
}
