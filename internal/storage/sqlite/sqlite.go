package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

// Store implements storage.Store using SQLite
type Store struct {
	db  *sql.DB
	log logger.Logger
}

// New creates a new SQLite store
func New(dbPath string, log logger.Logger) (*Store, error) {
	// Ensure directory exists
	if err := ensureDir(dbPath); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	log.Info("Opening database", "path", dbPath)

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	log.Info("Database connection successful")

	store := &Store{
		db:  db,
		log: log,
	}

	// Initialize schema
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// SaveNote saves a note to the database
func (s *Store) SaveNote(note string) error {
	_, err := s.db.Exec(`INSERT INTO notes (content) VALUES (?)`, note)
	if err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}
	s.log.Info("Note saved successfully")
	return nil
}

// GetNotes retrieves all notes from the database
func (s *Store) GetNotes() ([]string, error) {
	rows, err := s.db.Query(`SELECT content FROM notes ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var note string
		if err := rows.Scan(&note); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes: %w", err)
	}

	return notes, nil
}

// DeleteNote removes a note from the database
func (s *Store) DeleteNote(note string) error {
	result, err := s.db.Exec(`DELETE FROM notes WHERE content = ?`, note)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return storage.ErrTodoNotFound
	}

	return nil
}

// Clear removes all notes from the database
func (s *Store) Clear() error {
	_, err := s.db.Exec(`DELETE FROM notes`)
	if err != nil {
		return fmt.Errorf("failed to clear notes: %w", err)
	}
	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	s.log.Info("Database closed successfully")
	return nil
}

func (s *Store) initSchema() error {
	s.log.Info("Initializing database schema...")

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for i, migration := range migrations {
		s.log.Debug("Executing migration", "index", i, "query", migration)
		if _, err := s.db.Exec(migration); err != nil {
			s.log.Error("Migration failed", "index", i, "error", err)
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}

	s.log.Info("Schema initialized successfully")
	return nil
}

func ensureDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	return os.MkdirAll(dir, 0o755)
}
