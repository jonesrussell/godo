// Package sqlite provides a SQLite implementation of the storage interfaces
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements storage.TaskTxStore using SQLite
type Store struct {
	db  *sql.DB
	log logger.Logger
}

// Verify Store implements TaskTxStore
var _ storage.TaskTxStore = (*Store)(nil)

// New creates a new SQLite store
func New(path string, log logger.Logger) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &Store{
		db:  db,
		log: log,
	}

	if err := store.init(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return store, nil
}

// init creates the tasks table if it doesn't exist
func (s *Store) init() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	`
	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}
	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// List returns all tasks
func (s *Store) List(ctx context.Context) ([]storage.Task, error) {
	query := `
		SELECT id, content, done, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
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
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// GetByID returns a task by its ID
func (s *Store) GetByID(ctx context.Context, id string) (*storage.Task, error) {
	query := `
		SELECT id, content, done, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`
	var task storage.Task
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Content,
		&task.Done,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

// Add creates a new task
func (s *Store) Add(ctx context.Context, task storage.Task) error {
	query := `
		INSERT INTO tasks (id, content, done, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.ExecContext(ctx, query,
		task.ID,
		task.Content,
		task.Done,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

// Update replaces an existing task
func (s *Store) Update(ctx context.Context, task storage.Task) error {
	query := `
		UPDATE tasks
		SET content = ?, done = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := s.db.ExecContext(ctx, query,
		task.Content,
		task.Done,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

// Delete removes a task
func (s *Store) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(ctx context.Context) (storage.TaskTx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &Tx{tx: tx}, nil
}

// Tx implements storage.TaskTx
type Tx struct {
	tx *sql.Tx
}

// List returns all tasks within the transaction
func (t *Tx) List(ctx context.Context) ([]storage.Task, error) {
	query := `
		SELECT id, content, done, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`
	rows, err := t.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
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
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// GetByID returns a task by its ID within the transaction
func (t *Tx) GetByID(ctx context.Context, id string) (*storage.Task, error) {
	query := `
		SELECT id, content, done, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`
	var task storage.Task
	err := t.tx.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Content,
		&task.Done,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

// Add creates a new task within the transaction
func (t *Tx) Add(ctx context.Context, task storage.Task) error {
	query := `
		INSERT INTO tasks (id, content, done, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := t.tx.ExecContext(ctx, query,
		task.ID,
		task.Content,
		task.Done,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

// Update replaces an existing task within the transaction
func (t *Tx) Update(ctx context.Context, task storage.Task) error {
	query := `
		UPDATE tasks
		SET content = ?, done = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := t.tx.ExecContext(ctx, query,
		task.Content,
		task.Done,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

// Delete removes a task within the transaction
func (t *Tx) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := t.tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

// Commit commits the transaction
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
