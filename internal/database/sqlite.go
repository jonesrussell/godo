package database

import (
	"database/sql"
	"log"

	"github.com/jonesrussell/godo/pkg/database"
	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	log.Printf("Opening database at: %s", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}
	log.Println("Database connection successful")

	// Initialize schema
	log.Println("Initializing database schema...")
	if _, err := db.Exec(database.Schema); err != nil {
		log.Printf("Error initializing schema: %v", err)
		return nil, err
	}
	log.Println("Schema initialized successfully")

	return db, nil
}
