package testutil

import (
	"context"
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// MockStore provides a mock implementation of the storage.Store interface
type MockStore struct {
	mu     sync.RWMutex
	notes  map[string]storage.Note
	closed bool
	Error  error
}

// NewMockStore creates a new MockStore instance
func NewMockStore() *MockStore {
	return &MockStore{
		notes: make(map[string]storage.Note),
	}
}

// Add implements storage.Store.Add
func (m *MockStore) Add(ctx context.Context, note storage.Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if note.ID == "" {
		return storage.ErrEmptyID
	}

	if _, exists := m.notes[note.ID]; exists {
		return storage.ErrDuplicateID
	}

	m.notes[note.ID] = note
	return nil
}

// Get implements storage.Store.Get
func (m *MockStore) Get(ctx context.Context, id string) (storage.Note, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Error != nil {
		return storage.Note{}, m.Error
	}

	if id == "" {
		return storage.Note{}, storage.ErrEmptyID
	}

	note, exists := m.notes[id]
	if !exists {
		return storage.Note{}, storage.ErrNoteNotFound
	}

	return note, nil
}

// List implements storage.Store.List
func (m *MockStore) List(ctx context.Context) ([]storage.Note, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Error != nil {
		return nil, m.Error
	}

	notes := make([]storage.Note, 0, len(m.notes))
	for _, note := range m.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

// Update implements storage.Store.Update
func (m *MockStore) Update(ctx context.Context, note storage.Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if note.ID == "" {
		return storage.ErrEmptyID
	}

	if _, exists := m.notes[note.ID]; !exists {
		return storage.ErrNoteNotFound
	}

	note.UpdatedAt = time.Now().Unix()
	m.notes[note.ID] = note
	return nil
}

// Delete implements storage.Store.Delete
func (m *MockStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if id == "" {
		return storage.ErrEmptyID
	}

	if _, exists := m.notes[id]; !exists {
		return storage.ErrNoteNotFound
	}

	delete(m.notes, id)
	return nil
}

// Close implements storage.Store.Close
func (m *MockStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if m.closed {
		return storage.ErrStoreClosed
	}

	m.closed = true
	return nil
}

// BeginTx implements storage.Store.BeginTx
func (m *MockStore) BeginTx(ctx context.Context) (storage.Transaction, error) {
	return nil, storage.ErrTransactionNotSupported
}

// SetError sets the error to be returned by store operations
func (m *MockStore) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Error = err
}

// Reset clears all notes and resets the closed state (useful for testing)
func (m *MockStore) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notes = make(map[string]storage.Note)
	m.closed = false
	m.Error = nil
}
