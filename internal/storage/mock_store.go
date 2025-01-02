package storage

import (
	"context"
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage/errors"
)

// MockStore provides a mock implementation of TaskStore for testing
type MockStore struct {
	tasks  map[string]Task
	mu     sync.RWMutex
	Error  error
	closed bool
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make(map[string]Task),
	}
}

// Add stores a new task
func (s *MockStore) Add(_ context.Context, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if task.ID == "" {
		return errors.ErrEmptyID
	}

	if _, exists := s.tasks[task.ID]; exists {
		return errors.ErrDuplicateID
	}

	s.tasks[task.ID] = task
	return nil
}

// Get retrieves a task by its ID
func (s *MockStore) Get(_ context.Context, id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return Task{}, s.Error
	}

	if s.closed {
		return Task{}, errors.ErrStoreClosed
	}

	task, exists := s.tasks[id]
	if !exists {
		return Task{}, &errors.NotFoundError{ID: id}
	}

	return task, nil
}

// Update modifies an existing task
func (s *MockStore) Update(_ context.Context, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if _, exists := s.tasks[task.ID]; !exists {
		return &errors.NotFoundError{ID: task.ID}
	}

	task.UpdatedAt = time.Now().Unix()
	s.tasks[task.ID] = task
	return nil
}

// Delete removes a task by ID
func (s *MockStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if s.closed {
		return errors.ErrStoreClosed
	}

	if _, exists := s.tasks[id]; !exists {
		return &errors.NotFoundError{ID: id}
	}

	delete(s.tasks, id)
	return nil
}

// List returns all tasks
func (s *MockStore) List(_ context.Context) ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return nil, s.Error
	}

	if s.closed {
		return nil, errors.ErrStoreClosed
	}

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
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

// Reset clears all tasks and resets error state
func (s *MockStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks = make(map[string]Task)
	s.Error = nil
	s.closed = false
}

// SetError sets the error to be returned by all operations
func (s *MockStore) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Error = err
}
