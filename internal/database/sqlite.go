package database

import (
    "database/sql"
    "fmt"

    "github.com/jonesrussell/godo/internal/logger"
    _ "github.com/mattn/go-sqlite3"
)

const (
    SchemaVersion = 1
    Schema = `
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

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
    logger.Info("Opening database at: %s", dbPath)

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        logger.Error("Failed to open database: %v", err)
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Set connection pool settings
    db.SetMaxOpenConns(1) // SQLite only supports one writer
    db.SetMaxIdleConns(1)

    // Enable foreign keys
    if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
        logger.Error("Failed to enable foreign keys: %v", err)
        return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
    }

    // Test the connection
    if err := db.Ping(); err != nil {
        logger.Error("Error pinging database: %v", err)
        return nil, fmt.Errorf("database ping failed: %w", err)
    }
    logger.Info("Database connection successful")

    // Initialize schema
    logger.Info("Initializing database schema...")
    if err := initSchema(db); err != nil {
        logger.Error("Failed to initialize schema: %v", err)
        return nil, fmt.Errorf("schema initialization failed: %w", err)
    }
    logger.Info("Schema initialized successfully")

    return db, nil
}

func initSchema(db *sql.DB) error {
    // Create schema_version table if it doesn't exist
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_version (
            version INTEGER PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        logger.Error("Failed to create schema_version table: %v", err)
        return fmt.Errorf("failed to create schema_version table: %w", err)
    }

    // Check current schema version
    var version int
    err = db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
    if err != nil {
        logger.Error("Failed to get schema version: %v", err)
        return fmt.Errorf("failed to get schema version: %w", err)
    }

    logger.Debug("Current schema version: %d, Target version: %d", version, SchemaVersion)

    if version < SchemaVersion {
        // Begin transaction for schema updates
        tx, err := db.Begin()
        if err != nil {
            logger.Error("Failed to begin transaction: %v", err)
            return fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()

        // Apply schema changes
        if _, err := tx.Exec(Schema); err != nil {
            logger.Error("Failed to create todos table: %v", err)
            return fmt.Errorf("failed to create todos table: %w", err)
        }

        // Update schema version
        if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (?)", SchemaVersion); err != nil {
            logger.Error("Failed to update schema version: %v", err)
            return fmt.Errorf("failed to update schema version: %w", err)
        }

        // Commit transaction
        if err := tx.Commit(); err != nil {
            logger.Error("Failed to commit schema changes: %v", err)
            return fmt.Errorf("failed to commit schema changes: %w", err)
        }

        logger.Info("Schema updated to version %d", SchemaVersion)
    }

    return nil
} 