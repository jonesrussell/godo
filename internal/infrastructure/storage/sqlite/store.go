// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/errors"

	_ "modernc.org/sqlite" // SQLite driver
)

// Store implements storage.NoteStore using SQLite
type Store struct {
	db     *sql.DB
	logger logger.Logger
}

// New creates a new SQLite store
func New(path string, log logger.Logger) (*Store, error) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Error("Failed to create database directory", "dir", dir, "error", err)
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}
	log.Debug("Database directory ensured", "dir", dir)

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

// Add creates a new note in the store
func (s *Store) Add(ctx context.Context, note *model.Note) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO notes (id, content, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		note.ID, note.Content, note.Done, note.CreatedAt, note.UpdatedAt,
	)
	return err
}

// GetByID retrieves a note by its ID
func (s *Store) GetByID(ctx context.Context, id string) (model.Note, error) {
	var note model.Note
	err := s.db.QueryRowContext(ctx,
		"SELECT id, content, done, created_at, updated_at FROM notes WHERE id = ?",
		id,
	).Scan(&note.ID, &note.Content, &note.Done, &note.CreatedAt, &note.UpdatedAt)

	if err == sql.ErrNoRows {
		return model.Note{}, &errors.NotFoundError{ID: id}
	}
	return note, err
}

// Update modifies an existing note
func (s *Store) Update(ctx context.Context, note *model.Note) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE notes SET content = ?, done = ?, updated_at = ? WHERE id = ?",
		note.Content, note.Done, note.UpdatedAt, note.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &errors.NotFoundError{ID: note.ID}
	}
	return nil
}

// Delete removes a note by ID
func (s *Store) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM notes WHERE id = ?", id)
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

// List returns all notes
func (s *Store) List(ctx context.Context) ([]model.Note, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, content, done, created_at, updated_at FROM notes ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		var note model.Note
		scanErr := rows.Scan(&note.ID, &note.Content, &note.Done, &note.CreatedAt, &note.UpdatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(ctx context.Context) (storage.NoteTx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

// Transaction implements storage.NoteTx
type Transaction struct {
	tx *sql.Tx
}

// Add creates a new note in the transaction
func (t *Transaction) Add(ctx context.Context, note *model.Note) error {
	_, err := t.tx.ExecContext(ctx,
		"INSERT INTO notes (id, content, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		note.ID, note.Content, note.Done, note.CreatedAt, note.UpdatedAt,
	)
	return err
}

// List returns all notes in the transaction
func (t *Transaction) List(ctx context.Context) ([]model.Note, error) {
	rows, err := t.tx.QueryContext(ctx, "SELECT id, content, done, created_at, updated_at FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		var note model.Note
		if scanErr := rows.Scan(
			&note.ID,
			&note.Content,
			&note.Done,
			&note.CreatedAt,
			&note.UpdatedAt,
		); scanErr != nil {
			return nil, scanErr
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

// GetByID returns a note by its ID in the transaction
func (t *Transaction) GetByID(ctx context.Context, id string) (model.Note, error) {
	var note model.Note
	err := t.tx.QueryRowContext(ctx,
		"SELECT id, content, done, created_at, updated_at FROM notes WHERE id = ?",
		id,
	).Scan(
		&note.ID,
		&note.Content,
		&note.Done,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return model.Note{}, &errors.NotFoundError{ID: id}
	}
	if err != nil {
		return model.Note{}, err
	}
	return note, nil
}

// Update modifies an existing note in the transaction
func (t *Transaction) Update(ctx context.Context, note *model.Note) error {
	result, err := t.tx.ExecContext(ctx,
		"UPDATE notes SET content = ?, done = ?, updated_at = ? WHERE id = ?",
		note.Content, note.Done, note.UpdatedAt, note.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &errors.NotFoundError{ID: note.ID}
	}
	return nil
}

// Delete removes a note by ID in the transaction
func (t *Transaction) Delete(ctx context.Context, id string) error {
	result, err := t.tx.ExecContext(ctx, "DELETE FROM notes WHERE id = ?", id)
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
