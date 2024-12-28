package storage

import (
	"sync"
)

// MockStore provides a mock implementation of Store for testing
type MockStore struct {
	tasks map[string]Task
	mu    sync.RWMutex
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make(map[string]Task),
	}
}

// List returns all tasks
func (s *MockStore) List() ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Add creates a new task
func (s *MockStore) Add(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task
	return nil
}

// Update modifies an existing task
func (s *MockStore) Update(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	s.tasks[task.ID] = task
	return nil
}

// Delete removes a task
func (s *MockStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

// Close is a no-op for the mock store
func (s *MockStore) Close() error {
	return nil
}
