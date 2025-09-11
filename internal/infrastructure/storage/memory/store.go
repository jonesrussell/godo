// Package memory provides an in-memory implementation of the storage.NoteStore interface
package memory

import (
	"context"
	"sync"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

// Store implements storage.NoteStore interface using in-memory storage
type Store struct {
	notes map[string]model.Note
	mu    sync.RWMutex
}

// New creates a new memory store
func New() *Store {
	return &Store{
		notes: make(map[string]model.Note),
	}
}

// Add adds a new note to the store
func (s *Store) Add(_ context.Context, note *model.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[note.ID]; exists {
		return storage.ErrDuplicateID
	}

	s.notes[note.ID] = *note
	return nil
}

// GetByID retrieves a note by its ID
func (s *Store) GetByID(_ context.Context, id string) (model.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	note, exists := s.notes[id]
	if !exists {
		return model.Note{}, storage.ErrNoteNotFound
	}

	return note, nil
}

// List returns all notes in the store
func (s *Store) List(_ context.Context) ([]model.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notes := make([]model.Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}

	return notes, nil
}

// Update updates an existing note
func (s *Store) Update(_ context.Context, note *model.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[note.ID]; !exists {
		return storage.ErrNoteNotFound
	}

	s.notes[note.ID] = *note
	return nil
}

// Delete removes a note from the store
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[id]; !exists {
		return storage.ErrNoteNotFound
	}

	delete(s.notes, id)
	return nil
}

// Close implements storage.NoteStore interface
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notes = make(map[string]model.Note)
	return nil
}
