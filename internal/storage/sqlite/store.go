// Package sqlite provides SQLite-based implementation of the storage interface
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jonesrussell/godo/internal/storage/types"
	_ "modernc.org/sqlite"
)

// Store implements the types.Store interface using SQLite
type Store struct {
	db *sql.DB
}

// Transaction implements the types.Transaction interface
type Transaction struct {
	tx *sql.Tx
}

// New creates a new SQLite store instance
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &Store{db: db}
	if err := store.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return store, nil
}

// initialize creates the necessary tables if they don't exist
func (s *Store) initialize() error {
	query := `
		CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}

	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(ctx context.Context) (types.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &Transaction{tx: tx}, nil
}

// Add adds a new note
func (s *Store) Add(ctx context.Context, note types.Note) error {
	query := `
		INSERT INTO notes (id, content, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		note.ID,
		note.Content,
		note.Completed,
		note.CreatedAt,
		note.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add note: %w", err)
	}

	return nil
}

// Get retrieves a note by ID
func (s *Store) Get(ctx context.Context, id string) (types.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		WHERE id = ?
	`

	var note types.Note
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&note.ID,
		&note.Content,
		&note.Completed,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return types.Note{}, fmt.Errorf("note with ID %s not found", id)
	}

	if err != nil {
		return types.Note{}, fmt.Errorf("failed to get note: %w", err)
	}

	return note, nil
}

// List retrieves all notes
func (s *Store) List(ctx context.Context) ([]types.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}
	defer rows.Close()

	var notes []types.Note
	for rows.Next() {
		var note types.Note
		err := rows.Scan(
			&note.ID,
			&note.Content,
			&note.Completed,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes: %w", err)
	}

	return notes, nil
}

// Update updates an existing note
func (s *Store) Update(ctx context.Context, note types.Note) error {
	query := `
		UPDATE notes
		SET content = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`

	note.UpdatedAt = time.Now().Unix()
	result, err := s.db.ExecContext(ctx, query,
		note.Content,
		note.Completed,
		note.UpdatedAt,
		note.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note with ID %s not found", note.ID)
	}

	return nil
}

// Delete removes a note by ID
func (s *Store) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note with ID %s not found", id)
	}

	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

// Transaction methods

// Add adds a new note in the transaction
func (tx *Transaction) Add(ctx context.Context, note types.Note) error {
	query := `
		INSERT INTO notes (id, content, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := tx.tx.ExecContext(ctx, query,
		note.ID,
		note.Content,
		note.Completed,
		note.CreatedAt,
		note.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add note in transaction: %w", err)
	}

	return nil
}

// Get retrieves a note by ID in the transaction
func (tx *Transaction) Get(ctx context.Context, id string) (types.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		WHERE id = ?
	`

	var note types.Note
	err := tx.tx.QueryRowContext(ctx, query, id).Scan(
		&note.ID,
		&note.Content,
		&note.Completed,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return types.Note{}, fmt.Errorf("note with ID %s not found", id)
	}

	if err != nil {
		return types.Note{}, fmt.Errorf("failed to get note in transaction: %w", err)
	}

	return note, nil
}

// List retrieves all notes in the transaction
func (tx *Transaction) List(ctx context.Context) ([]types.Note, error) {
	query := `
		SELECT id, content, completed, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
	`

	rows, err := tx.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes in transaction: %w", err)
	}
	defer rows.Close()

	var notes []types.Note
	for rows.Next() {
		var note types.Note
		err := rows.Scan(
			&note.ID,
			&note.Content,
			&note.Completed,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note in transaction: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes in transaction: %w", err)
	}

	return notes, nil
}

// Update updates an existing note in the transaction
func (tx *Transaction) Update(ctx context.Context, note types.Note) error {
	query := `
		UPDATE notes
		SET content = ?, completed = ?, updated_at = ?
		WHERE id = ?
	`

	note.UpdatedAt = time.Now().Unix()
	result, err := tx.tx.ExecContext(ctx, query,
		note.Content,
		note.Completed,
		note.UpdatedAt,
		note.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update note in transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected in transaction: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note with ID %s not found", note.ID)
	}

	return nil
}

// Delete removes a note by ID in the transaction
func (tx *Transaction) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := tx.tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note in transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected in transaction: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note with ID %s not found", id)
	}

	return nil
}

// Commit commits the transaction
func (tx *Transaction) Commit() error {
	if err := tx.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back the transaction
func (tx *Transaction) Rollback() error {
	if err := tx.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}
