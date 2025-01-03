package mock

import (
	"context"
	"sync"

	"github.com/jonesrussell/godo/internal/storage/errors"
)

// MockStore implements the Store interface for testing
type MockStore struct {
	sync.RWMutex
	notes map[string]Note
	err   error
}

// New creates a new mock store
func New() *MockStore {
	return &MockStore{
		notes: make(map[string]Note),
	}
}

// List retrieves all notes
func (s *MockStore) List(_ context.Context) ([]Note, error) {
	s.RLock()
	defer s.RUnlock()

	if s.err != nil {
		return nil, s.err
	}

	notes := make([]Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}

	return notes, nil
}

// Get retrieves a note by ID
func (s *MockStore) Get(_ context.Context, id string) (Note, error) {
	s.RLock()
	defer s.RUnlock()

	if s.err != nil {
		return Note{}, s.err
	}

	note, exists := s.notes[id]
	if !exists {
		return Note{}, errors.ErrNoteNotFound
	}

	return note, nil
}

// Add adds a new note
func (s *MockStore) Add(_ context.Context, note Note) error {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.notes[note.ID]; exists {
		return errors.ErrNoteExists
	}

	s.notes[note.ID] = note
	return nil
}

// Update updates an existing note
func (s *MockStore) Update(_ context.Context, note Note) error {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.notes[note.ID]; !exists {
		return errors.ErrNoteNotFound
	}

	s.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID
func (s *MockStore) Delete(_ context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.notes[id]; !exists {
		return errors.ErrNoteNotFound
	}

	delete(s.notes, id)
	return nil
}

// BeginTx starts a new transaction
func (s *MockStore) BeginTx(_ context.Context) (Transaction, error) {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return nil, s.err
	}

	// Create a copy of the notes map for the transaction
	notes := make(map[string]Note)
	for k, v := range s.notes {
		notes[k] = v
	}

	return &Transaction{
		store:     s,
		notes:     notes,
		committed: false,
	}, nil
}
