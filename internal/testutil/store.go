package testutil

import (
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
		return fmt.Errorf("store already closed")
	}

	if task.ID == "" {
		return fmt.Errorf("empty ID")
	}

	if _, exists := m.tasks[task.ID]; exists {
		return fmt.Errorf("duplicate ID")
	}

	m.tasks[task.ID] = task
	return nil
}

// GetByID implements storage.Store.GetByID
func (m *MockStore) GetByID(id string) (*storage.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, fmt.Errorf("store already closed")
	}

	if id == "" {
		return nil, fmt.Errorf("empty ID")
	}

	task, exists := m.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	return &task, nil
}

// Update implements storage.Store.Update
func (m *MockStore) Update(task storage.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("store already closed")
	}

	if task.ID == "" {
		return fmt.Errorf("empty ID")
	}

	if _, exists := m.tasks[task.ID]; !exists {
		return fmt.Errorf("task not found")
	}

	task.UpdatedAt = time.Now()
	m.tasks[task.ID] = task
	return nil
}

// Delete implements storage.Store.Delete
func (m *MockStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("store already closed")
	}

	if id == "" {
		return fmt.Errorf("empty ID")
	}

	if _, exists := m.tasks[id]; !exists {
		return fmt.Errorf("task not found")
	}

	delete(m.tasks, id)
	return nil
}

// List implements storage.Store.List
func (m *MockStore) List() ([]storage.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, fmt.Errorf("store already closed")
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
