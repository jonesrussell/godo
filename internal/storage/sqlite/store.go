// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	_ "modernc.org/sqlite" // SQLite driver
)

// Store implements the storage.Store interface using SQLite
type Store struct {
	db     *sql.DB
	logger logger.Logger
}

// New creates a new SQLite store instance
func New(dbPath string, logger logger.Logger) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &Store{
		db:     db,
		logger: logger,
	}

	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return store, nil
}

// migrate creates the necessary database tables
func (s *Store) migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);
	`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}

	return nil
}

// Add adds a new task
func (s *Store) Add(ctx context.Context, task storage.Task) error {
	query := `
		INSERT INTO tasks (id, title, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		task.ID,
		task.Title,
		task.Completed,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

// Get retrieves a task by ID
func (s *Store) Get(ctx context.Context, id string) (storage.Task, error) {
	query := `
		SELECT id, title, completed, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	var task storage.Task
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return storage.Task{}, fmt.Errorf("task with ID %s not found", id)
	}
	if err != nil {
		return storage.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// List retrieves all tasks
func (s *Store) List(ctx context.Context) ([]storage.Task, error) {
	query := `
		SELECT id, title, completed, created_at, updated_at
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
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Update updates an existing task
func (s *Store) Update(ctx context.Context, task storage.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		task.Title,
		task.Completed,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	return nil
}

// Delete removes a task by ID
func (s *Store) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("task with ID %s not found", id)
	}

	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(ctx context.Context) (storage.TaskTx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &Transaction{tx: tx}, nil
}

// Transaction implements the storage.TaskTx interface
type Transaction struct {
	tx *sql.Tx
}

// Add adds a new task within the transaction
func (t *Transaction) Add(ctx context.Context, task storage.Task) error {
	query := `
		INSERT INTO tasks (id, title, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := t.tx.ExecContext(ctx, query,
		task.ID,
		task.Title,
		task.Completed,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

// Get retrieves a task by ID within the transaction
func (t *Transaction) Get(ctx context.Context, id string) (storage.Task, error) {
	query := `
		SELECT id, title, completed, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	var task storage.Task
	err := t.tx.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return storage.Task{}, fmt.Errorf("task with ID %s not found", id)
	}
	if err != nil {
		return storage.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// List retrieves all tasks within the transaction
func (t *Transaction) List(ctx context.Context) ([]storage.Task, error) {
	query := `
		SELECT id, title, completed, created_at, updated_at
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
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Update updates an existing task within the transaction
func (t *Transaction) Update(ctx context.Context, task storage.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := t.tx.ExecContext(ctx, query,
		task.Title,
		task.Completed,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	return nil
}

// Delete removes a task by ID within the transaction
func (t *Transaction) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := t.tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("task with ID %s not found", id)
	}

	return nil
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}
