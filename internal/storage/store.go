// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// Store defines the interface for data storage
type Store interface {
	Save(key string, value interface{}) error
	Load(key string) (interface{}, error)
	Delete(key string) error
	Close() error
	Add(key string, value interface{}) error
}

// SQLiteStore implements Store using SQLite
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite store
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Initialize tables
	if err := initTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Save(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT OR REPLACE INTO items (key, value)
		VALUES (?, ?)
	`, key, string(data))
	return err
}

func (s *SQLiteStore) Load(key string) (interface{}, error) {
	var data string
	err := s.db.QueryRow(`
		SELECT value FROM items WHERE key = ?
	`, key).Scan(&data)
	if err != nil {
		return nil, err
	}

	var value interface{}
	err = json.Unmarshal([]byte(data), &value)
	return value, err
}

func (s *SQLiteStore) Delete(key string) error {
	_, err := s.db.Exec(`
		DELETE FROM items WHERE key = ?
	`, key)
	return err
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) Add(key string, value interface{}) error {
	return s.Save(key, value)
}

func initTables(db *sql.DB) error {
	// Create tables if they don't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}
