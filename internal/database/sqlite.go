package database

import (
	"database/sql"

	"github.com/jonesrussell/godo/pkg/database"
	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize schema
	if _, err := db.Exec(database.Schema); err != nil {
		return nil, err
	}

	return db, nil
}
