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
				completed BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
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
