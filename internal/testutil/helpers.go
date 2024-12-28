// Package testutil provides testing utilities
package testutil

import (
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// TimePtr returns a pointer to the given time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// MockStore provides a mock implementation of storage.Store for testing
type MockStore struct {
	tasks []storage.Task
	mu    sync.RWMutex
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make([]storage.Task, 0),
	}
}

// Add implements storage.Store
func (s *MockStore) Add(task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
	return nil
}

// List implements storage.Store
func (s *MockStore) List() ([]storage.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks, nil
}

// Update implements storage.Store
func (s *MockStore) Update(task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == task.ID {
			s.tasks[i] = task
			return nil
		}
	}
	return storage.ErrTaskNotFound
}

// Delete implements storage.Store
func (s *MockStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return storage.ErrTaskNotFound
}