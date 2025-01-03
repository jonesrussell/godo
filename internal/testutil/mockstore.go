package testutil

import (
	"context"
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage/types"
)

// MockStore provides a mock implementation of the storage.Store interface
type MockStore struct {
	mu     sync.RWMutex
	notes  map[string]types.Note
	closed bool
	Error  error
}

// NewMockStore creates a new MockStore instance
func NewMockStore() *MockStore {
	return &MockStore{
		notes: make(map[string]types.Note),
	}
}

// Add implements Store.Add
func (m *MockStore) Add(ctx context.Context, note types.Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if note.ID == "" {
		return types.ErrEmptyID
	}

	if _, exists := m.notes[note.ID]; exists {
		return types.ErrDuplicateID
	}

	m.notes[note.ID] = note
	return nil
}

// Get implements Store.Get
func (m *MockStore) Get(ctx context.Context, id string) (types.Note, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Error != nil {
		return types.Note{}, m.Error
	}

	if id == "" {
		return types.Note{}, types.ErrEmptyID
	}

	note, exists := m.notes[id]
	if !exists {
		return types.Note{}, types.ErrNoteNotFound
	}

	return note, nil
}

// List implements Store.List
func (m *MockStore) List(ctx context.Context) ([]types.Note, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Error != nil {
		return nil, m.Error
	}

	notes := make([]types.Note, 0, len(m.notes))
	for _, note := range m.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

// Update implements Store.Update
func (m *MockStore) Update(ctx context.Context, note types.Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if note.ID == "" {
		return types.ErrEmptyID
	}

	if _, exists := m.notes[note.ID]; !exists {
		return types.ErrNoteNotFound
	}

	note.UpdatedAt = time.Now().Unix()
	m.notes[note.ID] = note
	return nil
}

// Delete implements Store.Delete
func (m *MockStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if id == "" {
		return types.ErrEmptyID
	}

	if _, exists := m.notes[id]; !exists {
		return types.ErrNoteNotFound
	}

	delete(m.notes, id)
	return nil
}

// Close implements Store.Close
func (m *MockStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if m.closed {
		return types.ErrStoreClosed
	}

	m.closed = true
	return nil
}

// BeginTx implements Store.BeginTx
func (m *MockStore) BeginTx(ctx context.Context) (types.Transaction, error) {
	return nil, types.ErrTransactionNotSupported
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
	m.notes = make(map[string]types.Note)
	m.closed = false
	m.Error = nil
}
