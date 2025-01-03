// Package mock provides mock implementations for testing
package mock

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/godo/internal/storage/types"
)

// Store implements a mock storage for testing
type Store struct {
	mu    sync.RWMutex
	notes map[string]types.Note
	err   error
}

// New creates a new mock store
func New() *Store {
	return &Store{
		notes: make(map[string]types.Note),
	}
}

// Add adds a note to the store
func (s *Store) Add(ctx context.Context, note types.Note) error {
	if s.err != nil {
		return s.err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[note.ID]; exists {
		return fmt.Errorf("note with ID %s already exists", note.ID)
	}

	s.notes[note.ID] = note
	return nil
}

// Get retrieves a note by ID
func (s *Store) Get(ctx context.Context, id string) (types.Note, error) {
	if s.err != nil {
		return types.Note{}, s.err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	note, exists := s.notes[id]
	if !exists {
		return types.Note{}, fmt.Errorf("note with ID %s not found", id)
	}

	return note, nil
}

// List returns all notes
func (s *Store) List(ctx context.Context) ([]types.Note, error) {
	if s.err != nil {
		return nil, s.err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	notes := make([]types.Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}

	return notes, nil
}

// Update modifies an existing note
func (s *Store) Update(ctx context.Context, note types.Note) error {
	if s.err != nil {
		return s.err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[note.ID]; !exists {
		return fmt.Errorf("note with ID %s not found", note.ID)
	}

	s.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID
func (s *Store) Delete(ctx context.Context, id string) error {
	if s.err != nil {
		return s.err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[id]; !exists {
		return fmt.Errorf("note with ID %s not found", id)
	}

	delete(s.notes, id)
	return nil
}

// Close is a no-op for mock store
func (s *Store) Close() error {
	if s.err != nil {
		return s.err
	}
	return nil
}

// BeginTx begins a new transaction
func (s *Store) BeginTx(ctx context.Context) (types.Transaction, error) {
	if s.err != nil {
		return nil, s.err
	}
	return nil, fmt.Errorf("transactions not supported in mock store")
}

// SetError sets an error to be returned by all operations
func (s *Store) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = err
}

// Reset clears all notes and errors
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes = make(map[string]types.Note)
	s.err = nil
}
