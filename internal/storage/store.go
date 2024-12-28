// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"database/sql"
	"sync"
	"time"

	_ "modernc.org/sqlite" // SQLite driver for database connectivity
)

// Task represents a todo task
type Task struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store defines the interface for data storage
type Store interface {
	List() ([]Task, error)
	Add(task Task) error
	Update(task Task) error
	Delete(id string) error
	GetByID(id string) (*Task, error)
	Close() error
}

// MemoryStore provides an in-memory task storage implementation
type MemoryStore struct {
	tasks []Task
	mu    sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make([]Task, 0),
	}
}

// Add stores a new task
func (s *MemoryStore) Add(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
	return nil
}

// List returns all stored tasks
func (s *MemoryStore) List() ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks, nil
}

// Update modifies an existing task
func (s *MemoryStore) Update(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == task.ID {
			s.tasks[i] = task
			return nil
		}
	}
	return ErrTaskNotFound
}

// Delete removes a task by ID
func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return ErrTaskNotFound
}

// Close is a no-op for memory store
func (s *MemoryStore) Close() error {
	return nil
}

// GetByID retrieves a task by its ID
func (s *MemoryStore) GetByID(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, ErrTaskNotFound
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

// List returns all tasks
func (s *SQLiteStore) List() ([]Task, error) {
	rows, err := s.db.Query(`
		SELECT id, content, done, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Content, &task.Done, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// Add creates a new task
func (s *SQLiteStore) Add(task Task) error {
	_, err := s.db.Exec(`
		INSERT INTO tasks (id, content, done, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, task.ID, task.Content, task.Done, task.CreatedAt, task.UpdatedAt)
	return err
}

// Update modifies an existing task
func (s *SQLiteStore) Update(task Task) error {
	result, err := s.db.Exec(`
		UPDATE tasks
		SET content = ?, done = ?, updated_at = ?
		WHERE id = ?
	`, task.Content, task.Done, time.Now(), task.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// Delete removes a task
func (s *SQLiteStore) Delete(id string) error {
	result, err := s.db.Exec(`
		DELETE FROM tasks
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func initTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	return err
}
