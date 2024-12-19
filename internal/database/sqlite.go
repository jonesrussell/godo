package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jonesrussell/godo/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

const (
	SchemaVersion = 1
	Schema        = `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT FALSE,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
)

// Config holds database configuration
type Config struct {
	Path string
}

// ensureDataDir creates the directory for the database if it doesn't exist
func ensureDataDir(dbPath string, log logger.Logger) error {
	dir := filepath.Dir(dbPath)
	if dir != "." {
		log.Debug("Creating database directory", "path", dir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Error("Failed to create database directory", "path", dir, "error", err)
			return err
		}
	}
	return nil
}

func NewSQLiteDB(dbPath string, log logger.Logger) (*sql.DB, error) {
	log.Info("Opening database", "path", dbPath)

	if err := ensureDataDir(dbPath, log); err != nil {
		log.Error("Failed to create database directory", "error", err)
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Error("Failed to open database", "error", err)
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Error("Failed to enable foreign keys", "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error("Database ping failed", "error", err)
		return nil, err
	}
	log.Info("Database connection successful")

	log.Info("Initializing database schema...")
	if err := initSchema(db, log); err != nil {
		log.Error("Schema initialization failed", "error", err)
		return nil, err
	}
	log.Info("Schema initialized successfully")

	return db, nil
}

func initSchema(db *sql.DB, log logger.Logger) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_version (
            version INTEGER PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Error("Failed to create schema_version table", "error", err)
		return err
	}

	var version int
	err = db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	if err != nil {
		log.Error("Failed to get schema version", "error", err)
		return err
	}

	log.Debug("Schema versions", "current", version, "target", SchemaVersion)

	if version < SchemaVersion {
		tx, err := db.Begin()
		if err != nil {
			log.Error("Failed to begin transaction", "error", err)
			return err
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Error("Failed to rollback transaction", "error", err)
			}
		}()

		if _, err := tx.Exec(Schema); err != nil {
			log.Error("Failed to create todos table", "error", err)
			return err
		}

		if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (?)", SchemaVersion); err != nil {
			log.Error("Failed to update schema version", "error", err)
			return err
		}

		if err := tx.Commit(); err != nil {
			log.Error("Failed to commit schema changes", "error", err)
			return err
		}

		log.Info("Schema updated", "version", SchemaVersion)
	}

	return nil
}

func TestConnection(db *sql.DB, log logger.Logger) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
	if err != nil {
		log.Error("Database test query failed", "error", err)
		return err
	}
	log.Debug("Current todo count", "count", count)
	return nil
}
