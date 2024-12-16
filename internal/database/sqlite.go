package database

import (
	"database/sql"

	"github.com/jonesrussell/godo/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	logger.Info("Opening database at: %s", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error("Failed to open database: %v", err)
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		logger.Error("Error pinging database: %v", err)
		return nil, err
	}
	logger.Info("Database connection successful")

	// Initialize schema
	logger.Info("Initializing database schema...")
	if err := initSchema(db); err != nil {
		logger.Error("Failed to initialize schema: %v", err)
		return nil, err
	}
	logger.Info("Schema initialized successfully")

	return db, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(Schema)
	if err != nil {
		logger.Error("Failed to create todos table: %v", err)
		return err
	}

	return nil
}
