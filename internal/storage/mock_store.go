package storage

import (
	"context"
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage/errors"
)

// MockStore provides a mock implementation of Store for testing
type MockStore struct {
	notes  map[string]Note
	mu     sync.RWMutex
	Error  error
	closed bool
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		notes: make(map[string]Note),
	}
}

// Add stores a new note
func (s *MockStore) Add(_ context.Context, note Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if note.ID == "" {
		return errors.ErrEmptyID
	}

	if _, exists := s.notes[note.ID]; exists {
		return errors.ErrDuplicateID
	}

	s.notes[note.ID] = note
	return nil
}

// Get retrieves a note by its ID
func (s *MockStore) Get(_ context.Context, id string) (Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return Note{}, s.Error
	}

	if s.closed {
		return Note{}, errors.ErrStoreClosed
	}

	note, exists := s.notes[id]
	if !exists {
		return Note{}, &errors.NotFoundError{ID: id}
	}

	return note, nil
}

// Update modifies an existing note
func (s *MockStore) Update(_ context.Context, note Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if _, exists := s.notes[note.ID]; !exists {
		return &errors.NotFoundError{ID: note.ID}
	}

	note.UpdatedAt = time.Now().Unix()
	s.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID
func (s *MockStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if _, exists := s.notes[id]; !exists {
		return &errors.NotFoundError{ID: id}
	}

	delete(s.notes, id)
	return nil
}

// List returns all notes
func (s *MockStore) List(_ context.Context) ([]Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return nil, s.Error
	}

	if s.closed {
		return nil, errors.ErrStoreClosed
	}

	notes := make([]Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

// Close marks the store as closed
func (s *MockStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	s.closed = true
	return nil
}

// BeginTx starts a new transaction
func (s *MockStore) BeginTx(_ context.Context) (Transaction, error) {
	return nil, ErrTransactionNotSupported
}

// Reset clears all notes and resets error state
func (s *MockStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notes = make(map[string]Note)
	s.Error = nil
	s.closed = false
}

// SetError sets the error to be returned by all operations
func (s *MockStore) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Error = err
}
