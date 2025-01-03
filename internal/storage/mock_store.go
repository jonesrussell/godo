package storage

import (
	"context"
	"sync"

	"github.com/jonesrussell/godo/internal/storage/errors"
)

// MockStore provides a mock implementation of Store for testing
type MockStore struct {
	sync.RWMutex
	notes     map[string]Note
	closed    bool
	Error     error
	AddCalled bool
}

// New creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		notes: make(map[string]Note),
	}
}

// List returns all notes
func (s *MockStore) List(ctx context.Context) ([]Note, error) {
	s.RLock()
	defer s.RUnlock()

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

// Get returns a note by ID
func (s *MockStore) Get(ctx context.Context, id string) (Note, error) {
	s.RLock()
	defer s.RUnlock()

	if s.Error != nil {
		return Note{}, s.Error
	}

	if s.closed {
		return Note{}, errors.ErrStoreClosed
	}

	if id == "" {
		return Note{}, errors.ErrEmptyID
	}

	note, ok := s.notes[id]
	if !ok {
		return Note{}, errors.ErrNoteNotFound
	}

	return note, nil
}

// Add creates a new note
func (s *MockStore) Add(ctx context.Context, note Note) error {
	s.Lock()
	defer s.Unlock()

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
	s.AddCalled = true
	return nil
}

// Update modifies an existing note
func (s *MockStore) Update(ctx context.Context, note Note) error {
	s.Lock()
	defer s.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if note.ID == "" {
		return errors.ErrEmptyID
	}

	if _, exists := s.notes[note.ID]; !exists {
		return errors.ErrNoteNotFound
	}

	s.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID
func (s *MockStore) Delete(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if id == "" {
		return errors.ErrEmptyID
	}

	if _, exists := s.notes[id]; !exists {
		return errors.ErrNoteNotFound
	}

	delete(s.notes, id)
	return nil
}

// BeginTx starts a new transaction
func (s *MockStore) BeginTx(ctx context.Context) (Transaction, error) {
	s.Lock()
	defer s.Unlock()

	if s.Error != nil {
		return nil, s.Error
	}

	if s.closed {
		return nil, errors.ErrStoreClosed
	}

	return &MockTransaction{store: s}, nil
}

// Close cleans up any resources
func (s *MockStore) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.Error != nil {
		return s.Error
	}

	s.closed = true
	return nil
}

// MockTransaction provides a mock implementation of Transaction for testing
type MockTransaction struct {
	store *MockStore
}

// List returns all notes in the transaction
func (t *MockTransaction) List(ctx context.Context) ([]Note, error) {
	return t.store.List(ctx)
}

// Get returns a note by ID in the transaction
func (t *MockTransaction) Get(ctx context.Context, id string) (Note, error) {
	return t.store.Get(ctx, id)
}

// Add creates a new note in the transaction
func (t *MockTransaction) Add(ctx context.Context, note Note) error {
	return t.store.Add(ctx, note)
}

// Update modifies an existing note in the transaction
func (t *MockTransaction) Update(ctx context.Context, note Note) error {
	return t.store.Update(ctx, note)
}

// Delete removes a note by ID in the transaction
func (t *MockTransaction) Delete(ctx context.Context, id string) error {
	return t.store.Delete(ctx, id)
}

// Commit commits the transaction
func (t *MockTransaction) Commit() error {
	return nil
}

// Rollback rolls back the transaction
func (t *MockTransaction) Rollback() error {
	return nil
}
