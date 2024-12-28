// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"database/sql"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	_ "modernc.org/sqlite" // SQLite driver
)

// Store implements storage.Store using SQLite
type Store struct {
	db     *sql.DB
	logger logger.Logger
}

// New creates a new SQLite store
func New(path string, logger logger.Logger) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
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

// Add creates a new task in the store
func (s *Store) Add(task storage.Task) error {
	_, err := s.db.Exec(
		"INSERT INTO tasks (id, content, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		task.ID, task.Content, task.Done, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

// List returns all tasks in the store
func (s *Store) List() ([]storage.Task, error) {
	rows, err := s.db.Query("SELECT id, content, done, created_at, updated_at FROM tasks")
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

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}
