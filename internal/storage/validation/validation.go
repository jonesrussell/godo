// Package validation provides validation functions for storage operations
package validation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jonesrussell/godo/internal/domain/note"
	"github.com/jonesrussell/godo/internal/storage/types"
)

// Store wraps another store and adds validation
type Store struct {
	store types.Store
}

// New creates a new validation store wrapper
func New(store types.Store) *Store {
	return &Store{store: store}
}

// Add validates and adds a new note
func (s *Store) Add(ctx context.Context, n *note.Note) error {
	if err := validateNote(n); err != nil {
		return fmt.Errorf("invalid note: %w", err)
	}
	return s.store.Add(ctx, n)
}

// Get retrieves a note by ID
func (s *Store) Get(ctx context.Context, id string) (*note.Note, error) {
	if err := validateID(id); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	return s.store.Get(ctx, id)
}

// List retrieves all notes
func (s *Store) List(ctx context.Context) ([]*note.Note, error) {
	return s.store.List(ctx)
}

// Update validates and updates an existing note
func (s *Store) Update(ctx context.Context, n *note.Note) error {
	if err := validateNote(n); err != nil {
		return fmt.Errorf("invalid note: %w", err)
	}
	return s.store.Update(ctx, n)
}

// Delete removes a note by ID
func (s *Store) Delete(ctx context.Context, id string) error {
	if err := validateID(id); err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.store.Delete(ctx, id)
}

// Close closes the underlying store
func (s *Store) Close() error {
	return s.store.Close()
}

// BeginTx begins a new transaction
func (s *Store) BeginTx(ctx context.Context) (types.Transaction, error) {
	return s.store.BeginTx(ctx)
}

// validateNote validates a note
func validateNote(n *note.Note) error {
	if n == nil {
		return fmt.Errorf("note cannot be nil")
	}

	if err := validateID(n.ID); err != nil {
		return err
	}

	if strings.TrimSpace(n.Content) == "" {
		return fmt.Errorf("content cannot be empty")
	}

	if n.CreatedAt.IsZero() {
		return fmt.Errorf("created_at must be set")
	}

	if n.UpdatedAt.IsZero() {
		return fmt.Errorf("updated_at must be set")
	}

	if n.UpdatedAt.Before(n.CreatedAt) {
		return fmt.Errorf("updated_at cannot be before created_at")
	}

	if n.UpdatedAt.After(time.Now().Add(time.Second)) {
		return fmt.Errorf("updated_at cannot be in the future")
	}

	return nil
}

// validateID validates a note ID
func validateID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id cannot be empty")
	}
	return nil
}
