// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	_ "modernc.org/sqlite" // SQLite driver
)

var (
	// ErrEmptyID is returned when a task ID is empty
	ErrEmptyID = errors.New("task ID cannot be empty")
	// ErrInvalidPath is returned when the database path is invalid
	ErrInvalidPath = errors.New("invalid database path")
	// ErrStoreClosed is returned when attempting to use a closed store
	ErrStoreClosed = errors.New("store is closed")
)

// Store implements storage.Store using SQLite
type Store struct {
	db     *sql.DB
	logger logger.Logger
	closed bool
}

// New creates a new SQLite store
func New(path string, logger logger.Logger) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		return nil, ErrInvalidPath
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Configure SQLite for better concurrency
	if _, err := db.Exec(`PRAGMA journal_mode = WAL`); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set journal mode: %w", err)
	}
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000`); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	store := &Store{
		db:     db,
		logger: logger,
	}

	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

// validateID checks if a task ID is valid
func (s *Store) validateID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrEmptyID
	}
	return nil
}

// checkClosed returns an error if the store is closed
func (s *Store) checkClosed() error {
	if s.closed {
		return ErrStoreClosed
	}
	return nil
}

// Add creates a new task in the store
func (s *Store) Add(task storage.Task) error {
	if err := s.checkClosed(); err != nil {
		return err
	}
	if err := s.validateID(task.ID); err != nil {
		return err
	}

	_, err := s.db.Exec(
		"INSERT INTO tasks (id, content, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		task.ID, task.Content, task.Done, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

// List returns all tasks in the store
func (s *Store) List() ([]storage.Task, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT id, content, done, created_at, updated_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []storage.Task
	for rows.Next() {
		var task storage.Task
		if err := rows.Scan(
			&task.ID,
			&task.Content,
			&task.Done,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// Update modifies an existing task
func (s *Store) Update(task storage.Task) error {
	if err := s.checkClosed(); err != nil {
		return err
	}
	if err := s.validateID(task.ID); err != nil {
		return err
	}

	result, err := s.db.Exec(
		"UPDATE tasks SET content = ?, done = ?, updated_at = ? WHERE id = ?",
		task.Content, task.Done, time.Now(), task.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	return nil
}

// Delete removes a task by ID
func (s *Store) Delete(id string) error {
	if err := s.checkClosed(); err != nil {
		return err
	}
	if err := s.validateID(id); err != nil {
		return err
	}

	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	return nil
}

// GetByID retrieves a task by its ID
func (s *Store) GetByID(id string) (*storage.Task, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}
	if err := s.validateID(id); err != nil {
		return nil, err
	}

	var task storage.Task
	err := s.db.QueryRow(
		"SELECT id, content, done, created_at, updated_at FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.Content, &task.Done, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, storage.ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	return s.db.Close()
}
