// Package mock provides a mock implementation of the storage interface for testing
package mock

import (
	"sync"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"
)

// Store provides a mock implementation of storage.Store for testing
type Store struct {
	mu     sync.RWMutex
	tasks  map[string]storage.Task
	closed bool

	// Operation tracking for tests
	AddCalled     bool
	UpdateCalled  bool
	DeleteCalled  bool
	ListCalled    bool
	GetByIDCalled bool

	// Error simulation
	Error error
}

// New creates a new mock store
func New() *Store {
	return &Store{
		tasks: make(map[string]storage.Task),
	}
}

// Add simulates adding a task
func (m *Store) Add(task storage.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.AddCalled = true
	if m.Error != nil {
		return m.Error
	}

	if m.closed {
		return errors.ErrStoreClosed
	}

	if task.ID == "" {
		return errors.ErrEmptyID
	}

	if _, exists := m.tasks[task.ID]; exists {
		return errors.ErrDuplicateID
	}

	m.tasks[task.ID] = task
	return nil
}

// Update simulates updating a task
func (m *Store) Update(task storage.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.UpdateCalled = true
	if m.Error != nil {
		return m.Error
	}

	if m.closed {
		return errors.ErrStoreClosed
	}

	if task.ID == "" {
		return errors.ErrEmptyID
	}

	if _, exists := m.tasks[task.ID]; !exists {
		return errors.ErrTaskNotFound
	}

	m.tasks[task.ID] = task
	return nil
}

// Delete simulates deleting a task
func (m *Store) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DeleteCalled = true
	if m.Error != nil {
		return m.Error
	}

	if m.closed {
		return errors.ErrStoreClosed
	}

	if id == "" {
		return errors.ErrEmptyID
	}

	if _, exists := m.tasks[id]; !exists {
		return errors.ErrTaskNotFound
	}

	delete(m.tasks, id)
	return nil
}

// List simulates listing all tasks
func (m *Store) List() ([]storage.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.ListCalled = true
	if m.Error != nil {
		return nil, m.Error
	}

	if m.closed {
		return nil, errors.ErrStoreClosed
	}

	tasks := make([]storage.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByID simulates retrieving a task by ID
func (m *Store) GetByID(id string) (*storage.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.GetByIDCalled = true
	if m.Error != nil {
		return nil, m.Error
	}

	if m.closed {
		return nil, errors.ErrStoreClosed
	}

	if id == "" {
		return nil, errors.ErrEmptyID
	}

	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.ErrTaskNotFound
	}

	return &task, nil
}

// Close simulates closing the store
func (m *Store) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Error != nil {
		return m.Error
	}

	if m.closed {
		return errors.ErrStoreClosed
	}

	m.closed = true
	return nil
}

// Reset resets the mock store to its initial state
func (m *Store) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tasks = make(map[string]storage.Task)
	m.closed = false
	m.AddCalled = false
	m.UpdateCalled = false
	m.DeleteCalled = false
	m.ListCalled = false
	m.GetByIDCalled = false
	m.Error = nil
}
