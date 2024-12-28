package storage

import (
	"context"
	"sync"
)

// MockStore provides a mock implementation of TaskStore for testing
type MockStore struct {
	tasks map[string]Task
	mu    sync.RWMutex
	Error error // Error to return on next operation
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make(map[string]Task),
	}
}

// Reset clears all tasks and resets error
func (s *MockStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = make(map[string]Task)
	s.Error = nil
}

// List returns all tasks
func (s *MockStore) List(ctx context.Context) ([]Task, error) {
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

// Add creates a new task
func (s *MockStore) Add(ctx context.Context, task Task) error {
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

// Update modifies an existing task
func (s *MockStore) Update(ctx context.Context, task Task) error {
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

// Delete removes a task
func (s *MockStore) Delete(ctx context.Context, id string) error {
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

// Close is a no-op for the mock store
func (s *MockStore) Close() error {
	return nil
}

// GetByID retrieves a task by its ID
func (s *MockStore) GetByID(ctx context.Context, id string) (*Task, error) {
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
