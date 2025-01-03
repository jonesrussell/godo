// Package memory provides an in-memory implementation of the storage interface
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/godo/internal/storage/types"
)

// Store implements an in-memory storage
type Store struct {
	mu    sync.RWMutex
	notes map[string]types.Note
}

// New creates a new memory store
func New() *Store {
	return &Store{
		notes: make(map[string]types.Note),
	}
}

// Add stores a new note
func (s *Store) Add(ctx context.Context, note types.Note) error {
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
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[id]; !exists {
		return fmt.Errorf("note with ID %s not found", id)
	}

	delete(s.notes, id)
	return nil
}

// Close is a no-op for memory store
func (s *Store) Close() error {
	return nil
}

// BeginTx begins a new transaction
func (s *Store) BeginTx(ctx context.Context) (types.Transaction, error) {
	return nil, fmt.Errorf("transactions not supported in memory store")
}
