package storage

import (
	"context"
	"sync"
)

// MockStore is a mock implementation of TaskStore for testing
type MockStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
	Error error
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make(map[string]Task),
	}
}

// Reset clears all tasks and resets error state
func (s *MockStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = make(map[string]Task)
	s.Error = nil
}

// List returns all tasks in the store
func (s *MockStore) List(_ context.Context) ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return nil, s.Error
	}

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Add adds a new task to the store
func (s *MockStore) Add(_ context.Context, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if _, exists := s.tasks[task.ID]; exists {
		return ErrDuplicateID
	}

	s.tasks[task.ID] = task
	return nil
}

// Update updates an existing task
func (s *MockStore) Update(_ context.Context, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if _, exists := s.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	s.tasks[task.ID] = task
	return nil
}

// Delete removes a task from the store
func (s *MockStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

// GetByID retrieves a task by its ID
func (s *MockStore) GetByID(_ context.Context, id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return nil, s.Error
	}

	task, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return &task, nil
}

// Close implements TaskStore interface
func (s *MockStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	s.tasks = make(map[string]Task)
	return nil
}
