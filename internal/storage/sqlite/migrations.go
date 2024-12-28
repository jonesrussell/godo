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
			`CREATE TABLE IF NOT EXISTS tasks (
				id TEXT PRIMARY KEY,
				title TEXT NOT NULL,
				description TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				completed_at TIMESTAMP
			)`,
			`-- Migration to update existing tasks table
			ALTER TABLE tasks ADD COLUMN IF NOT EXISTS description TEXT;
			ALTER TABLE tasks ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP;
			ALTER TABLE tasks DROP COLUMN IF EXISTS completed;`,
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
