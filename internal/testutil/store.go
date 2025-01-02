package testutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// MockStore provides a thread-safe mock implementation of storage.Store
type MockStore struct {
	mu     sync.RWMutex
	tasks  map[string]storage.Task
	closed bool
	Error  error
}

// NewMockStore creates a new MockStore instance
func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make(map[string]storage.Task),
	}
}

// Add implements storage.Store.Add
func (m *MockStore) Add(task storage.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return storage.ErrStoreClosed
	}

	if task.ID == "" {
		return storage.ErrEmptyID
	}

	if _, exists := m.tasks[task.ID]; exists {
		return storage.ErrDuplicateID
	}

	m.tasks[task.ID] = task
	return nil
}

// GetByID implements storage.Store.GetByID
func (m *MockStore) GetByID(id string) (*storage.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, storage.ErrStoreClosed
	}

	if id == "" {
		return nil, storage.ErrEmptyID
	}

	task, exists := m.tasks[id]
	if !exists {
		return nil, storage.ErrTaskNotFound
	}

	return &task, nil
}

// Update implements storage.Store.Update
func (m *MockStore) Update(task storage.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return storage.ErrStoreClosed
	}

	if task.ID == "" {
		return storage.ErrEmptyID
	}

	if _, exists := m.tasks[task.ID]; !exists {
		return storage.ErrTaskNotFound
	}

	task.UpdatedAt = time.Now().Unix()
	m.tasks[task.ID] = task
	return nil
}

// Delete implements storage.Store.Delete
func (m *MockStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return storage.ErrStoreClosed
	}

	if id == "" {
		return storage.ErrEmptyID
	}

	if _, exists := m.tasks[id]; !exists {
		return storage.ErrTaskNotFound
	}

	delete(m.tasks, id)
	return nil
}

// List implements storage.Store.List
func (m *MockStore) List() ([]storage.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, storage.ErrStoreClosed
	}

	tasks := make([]storage.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Close implements storage.Store.Close
func (m *MockStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("store already closed")
	}

	m.closed = true
	return nil
}

// Reset clears all tasks and resets the closed state (useful for testing)
func (m *MockStore) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tasks = make(map[string]storage.Task)
	m.closed = false
}

// Get implements storage.Store.Get
func (m *MockStore) Get(ctx context.Context, id string) (storage.Task, error) {
	if m.Error != nil {
		return storage.Task{}, m.Error
	}
	task, exists := m.tasks[id]
	if !exists {
		return storage.Task{}, storage.ErrTaskNotFound
	}
	return task, nil
}
