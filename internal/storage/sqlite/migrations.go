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
			// Initial schema
			`CREATE TABLE IF NOT EXISTS tasks (
				id TEXT PRIMARY KEY,
				title TEXT NOT NULL,
				description TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				completed_at TIMESTAMP
			);`,
			// Add new columns if they don't exist
			`PRAGMA foreign_keys=off;
			BEGIN TRANSACTION;
			
			CREATE TABLE IF NOT EXISTS _tasks_new (
				id TEXT PRIMARY KEY,
				title TEXT NOT NULL,
				description TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				completed_at TIMESTAMP
			);
			
			INSERT INTO _tasks_new (id, title)
			SELECT id, title FROM tasks;
			
			DROP TABLE IF EXISTS tasks;
			ALTER TABLE _tasks_new RENAME TO tasks;
			
			COMMIT;
			PRAGMA foreign_keys=on;`,
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
