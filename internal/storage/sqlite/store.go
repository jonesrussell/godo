// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"context"
	"database/sql"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"

	_ "modernc.org/sqlite" // SQLite driver
)

// Store implements storage.TaskStore using SQLite
type Store struct {
	db     *sql.DB
	logger logger.Logger
}

// New creates a new SQLite store
func New(path string, log logger.Logger) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &Store{
		db:     db,
		logger: log,
	}

	if migErr := RunMigrations(db); migErr != nil {
		db.Close()
		return nil, migErr
	}

	return store, nil
}

// Add creates a new task in the store
func (s *Store) Add(ctx context.Context, task *storage.Task) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO tasks (id, content, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		task.ID, task.Content, task.Done, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

// GetByID retrieves a task by its ID
func (s *Store) GetByID(ctx context.Context, id string) (storage.Task, error) {
	var task storage.Task
	err := s.db.QueryRowContext(ctx,
		"SELECT id, content, done, created_at, updated_at FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.Content, &task.Done, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return storage.Task{}, &errors.NotFoundError{ID: id}
	}
	return task, err
}

// Update modifies an existing task
func (s *Store) Update(ctx context.Context, task *storage.Task) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE tasks SET content = ?, done = ?, updated_at = ? WHERE id = ?",
		task.Content, task.Done, task.UpdatedAt, task.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &errors.NotFoundError{ID: task.ID}
	}
	return nil
}

// Delete removes a task by ID
func (s *Store) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &errors.NotFoundError{ID: id}
	}
	return nil
}

// List returns all tasks
func (s *Store) List(ctx context.Context) ([]storage.Task, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, content, done, created_at, updated_at FROM tasks ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []storage.Task
	for rows.Next() {
		var task storage.Task
		scanErr := rows.Scan(&task.ID, &task.Content, &task.Done, &task.CreatedAt, &task.UpdatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(ctx context.Context) (storage.TaskTx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

// Transaction implements storage.TaskTx
type Transaction struct {
	tx *sql.Tx
}

// Add creates a new task in the transaction
func (t *Transaction) Add(ctx context.Context, task *storage.Task) error {
	_, err := t.tx.ExecContext(ctx,
		"INSERT INTO tasks (id, content, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		task.ID, task.Content, task.Done, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

// List returns all tasks in the transaction
func (t *Transaction) List(ctx context.Context) ([]storage.Task, error) {
	rows, err := t.tx.QueryContext(ctx, "SELECT id, content, done, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []storage.Task
	for rows.Next() {
		var task storage.Task
		if scanErr := rows.Scan(
			&task.ID,
			&task.Content,
			&task.Done,
			&task.CreatedAt,
			&task.UpdatedAt,
		); scanErr != nil {
			return nil, scanErr
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// GetByID returns a task by its ID in the transaction
func (t *Transaction) GetByID(ctx context.Context, id string) (storage.Task, error) {
	var task storage.Task
	err := t.tx.QueryRowContext(ctx,
		"SELECT id, content, done, created_at, updated_at FROM tasks WHERE id = ?",
		id,
	).Scan(
		&task.ID,
		&task.Content,
		&task.Done,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return storage.Task{}, &errors.NotFoundError{ID: id}
	}
	if err != nil {
		return storage.Task{}, err
	}
	return task, nil
}

// Update modifies an existing task in the transaction
func (t *Transaction) Update(ctx context.Context, task *storage.Task) error {
	result, err := t.tx.ExecContext(ctx,
		"UPDATE tasks SET content = ?, done = ?, updated_at = ? WHERE id = ?",
		task.Content, task.Done, task.UpdatedAt, task.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &errors.NotFoundError{ID: task.ID}
	}
	return nil
}

// Delete removes a task by ID in the transaction
func (t *Transaction) Delete(ctx context.Context, id string) error {
	result, err := t.tx.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &errors.NotFoundError{ID: id}
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
