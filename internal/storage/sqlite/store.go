// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jonesrussell/godo/internal/domain/note"
	_ "modernc.org/sqlite"
)

// Store implements the note.Store interface using SQLite
type Store struct {
	db *sql.DB
}

// New creates a new SQLite store
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, &note.Error{
			Op:   "sqlite.New",
			Kind: note.StorageError,
			Msg:  "failed to open database",
			Err:  err,
		}
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, &note.Error{
			Op:   "sqlite.New",
			Kind: note.StorageError,
			Msg:  "failed to ping database",
			Err:  err,
		}
	}

	if err := createSchema(db); err != nil {
		db.Close()
		return nil, &note.Error{
			Op:   "sqlite.New",
			Kind: note.StorageError,
			Msg:  "failed to create schema",
			Err:  err,
		}
	}

	return &Store{db: db}, nil
}

// Add adds a new note
func (s *Store) Add(ctx context.Context, n *note.Note) error {
	if err := n.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO notes (id, content, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		n.ID,
		n.Content,
		n.Completed,
		n.CreatedAt.Unix(),
		n.UpdatedAt.Unix(),
	)
	if err != nil {
		return &note.Error{
			Op:   "Store.Add",
			Kind: note.StorageError,
			Msg:  "failed to insert note",
			Err:  err,
		}
	}

	return nil
}

// Get retrieves a note by ID
func (s *Store) Get(ctx context.Context, id string) (*note.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		WHERE id = ?
	`

	var n note.Note
	var createdAt, updatedAt int64
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&n.ID,
		&n.Content,
		&n.Completed,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, &note.Error{
			Op:   "Store.Get",
			Kind: note.NotFound,
			Msg:  fmt.Sprintf("note with id %s not found", id),
		}
	}
	if err != nil {
		return nil, &note.Error{
			Op:   "Store.Get",
			Kind: note.StorageError,
			Msg:  "failed to get note",
			Err:  err,
		}
	}

	n.CreatedAt = time.Unix(createdAt, 0)
	n.UpdatedAt = time.Unix(updatedAt, 0)
	return &n, nil
}

// List returns all notes
func (s *Store) List(ctx context.Context) ([]*note.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, &note.Error{
			Op:   "Store.List",
			Kind: note.StorageError,
			Msg:  "failed to query notes",
			Err:  err,
		}
	}
	defer rows.Close()

	var notes []*note.Note
	for rows.Next() {
		var n note.Note
		var createdAt, updatedAt int64
		err := rows.Scan(
			&n.ID,
			&n.Content,
			&n.Completed,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, &note.Error{
				Op:   "Store.List",
				Kind: note.StorageError,
				Msg:  "failed to scan note",
				Err:  err,
			}
		}
		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)
		notes = append(notes, &n)
	}

	if err := rows.Err(); err != nil {
		return nil, &note.Error{
			Op:   "Store.List",
			Kind: note.StorageError,
			Msg:  "error iterating notes",
			Err:  err,
		}
	}

	return notes, nil
}

// Update modifies an existing note
func (s *Store) Update(ctx context.Context, n *note.Note) error {
	if err := n.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE notes
		SET content = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		n.Content,
		n.Completed,
		n.UpdatedAt.Unix(),
		n.ID,
	)
	if err != nil {
		return &note.Error{
			Op:   "Store.Update",
			Kind: note.StorageError,
			Msg:  "failed to update note",
			Err:  err,
		}
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return &note.Error{
			Op:   "Store.Update",
			Kind: note.StorageError,
			Msg:  "failed to get affected rows",
			Err:  err,
		}
	}

	if rows == 0 {
		return &note.Error{
			Op:   "Store.Update",
			Kind: note.NotFound,
			Msg:  fmt.Sprintf("note with id %s not found", n.ID),
		}
	}

	return nil
}

// Delete removes a note by ID
func (s *Store) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return &note.Error{
			Op:   "Store.Delete",
			Kind: note.StorageError,
			Msg:  "failed to delete note",
			Err:  err,
		}
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return &note.Error{
			Op:   "Store.Delete",
			Kind: note.StorageError,
			Msg:  "failed to get affected rows",
			Err:  err,
		}
	}

	if rows == 0 {
		return &note.Error{
			Op:   "Store.Delete",
			Kind: note.NotFound,
			Msg:  fmt.Sprintf("note with id %s not found", id),
		}
	}

	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return &note.Error{
			Op:   "Store.Close",
			Kind: note.StorageError,
			Msg:  "failed to close database",
			Err:  err,
		}
	}
	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(ctx context.Context) (note.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, &note.Error{
			Op:   "Store.BeginTx",
			Kind: note.StorageError,
			Msg:  "failed to begin transaction",
			Err:  err,
		}
	}
	return &Transaction{tx: tx}, nil
}

// Transaction represents a database transaction
type Transaction struct {
	tx *sql.Tx
}

// Add adds a new note in the transaction
func (t *Transaction) Add(ctx context.Context, n *note.Note) error {
	if err := n.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO notes (id, content, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := t.tx.ExecContext(ctx, query,
		n.ID,
		n.Content,
		n.Completed,
		n.CreatedAt.Unix(),
		n.UpdatedAt.Unix(),
	)
	if err != nil {
		return &note.Error{
			Op:   "Transaction.Add",
			Kind: note.StorageError,
			Msg:  "failed to insert note",
			Err:  err,
		}
	}

	return nil
}

// Get retrieves a note by ID in the transaction
func (t *Transaction) Get(ctx context.Context, id string) (*note.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		WHERE id = ?
	`

	var n note.Note
	var createdAt, updatedAt int64
	err := t.tx.QueryRowContext(ctx, query, id).Scan(
		&n.ID,
		&n.Content,
		&n.Completed,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, &note.Error{
			Op:   "Transaction.Get",
			Kind: note.NotFound,
			Msg:  fmt.Sprintf("note with id %s not found", id),
		}
	}
	if err != nil {
		return nil, &note.Error{
			Op:   "Transaction.Get",
			Kind: note.StorageError,
			Msg:  "failed to get note",
			Err:  err,
		}
	}

	n.CreatedAt = time.Unix(createdAt, 0)
	n.UpdatedAt = time.Unix(updatedAt, 0)
	return &n, nil
}

// List returns all notes in the transaction
func (t *Transaction) List(ctx context.Context) ([]*note.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
	`

	rows, err := t.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, &note.Error{
			Op:   "Transaction.List",
			Kind: note.StorageError,
			Msg:  "failed to query notes",
			Err:  err,
		}
	}
	defer rows.Close()

	var notes []*note.Note
	for rows.Next() {
		var n note.Note
		var createdAt, updatedAt int64
		err := rows.Scan(
			&n.ID,
			&n.Content,
			&n.Completed,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, &note.Error{
				Op:   "Transaction.List",
				Kind: note.StorageError,
				Msg:  "failed to scan note",
				Err:  err,
			}
		}
		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)
		notes = append(notes, &n)
	}

	if err := rows.Err(); err != nil {
		return nil, &note.Error{
			Op:   "Transaction.List",
			Kind: note.StorageError,
			Msg:  "error iterating notes",
			Err:  err,
		}
	}

	return notes, nil
}

// Update modifies an existing note in the transaction
func (t *Transaction) Update(ctx context.Context, n *note.Note) error {
	if err := n.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE notes
		SET content = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := t.tx.ExecContext(ctx, query,
		n.Content,
		n.Completed,
		n.UpdatedAt.Unix(),
		n.ID,
	)
	if err != nil {
		return &note.Error{
			Op:   "Transaction.Update",
			Kind: note.StorageError,
			Msg:  "failed to update note",
			Err:  err,
		}
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return &note.Error{
			Op:   "Transaction.Update",
			Kind: note.StorageError,
			Msg:  "failed to get affected rows",
			Err:  err,
		}
	}

	if rows == 0 {
		return &note.Error{
			Op:   "Transaction.Update",
			Kind: note.NotFound,
			Msg:  fmt.Sprintf("note with id %s not found", n.ID),
		}
	}

	return nil
}

// Delete removes a note by ID in the transaction
func (t *Transaction) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := t.tx.ExecContext(ctx, query, id)
	if err != nil {
		return &note.Error{
			Op:   "Transaction.Delete",
			Kind: note.StorageError,
			Msg:  "failed to delete note",
			Err:  err,
		}
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return &note.Error{
			Op:   "Transaction.Delete",
			Kind: note.StorageError,
			Msg:  "failed to get affected rows",
			Err:  err,
		}
	}

	if rows == 0 {
		return &note.Error{
			Op:   "Transaction.Delete",
			Kind: note.NotFound,
			Msg:  fmt.Sprintf("note with id %s not found", id),
		}
	}

	return nil
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	if err := t.tx.Commit(); err != nil {
		return &note.Error{
			Op:   "Transaction.Commit",
			Kind: note.StorageError,
			Msg:  "failed to commit transaction",
			Err:  err,
		}
	}
	return nil
}

// Rollback aborts the transaction
func (t *Transaction) Rollback() error {
	if err := t.tx.Rollback(); err != nil {
		return &note.Error{
			Op:   "Transaction.Rollback",
			Kind: note.StorageError,
			Msg:  "failed to rollback transaction",
			Err:  err,
		}
	}
	return nil
}

func createSchema(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}
