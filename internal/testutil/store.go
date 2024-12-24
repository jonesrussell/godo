package testutil

import (
	"github.com/jonesrussell/godo/internal/storage"
)

// MockStore provides a mock implementation of storage.Store for testing
type MockStore struct {
	tasks []storage.Task
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		tasks: make([]storage.Task, 0),
	}
}

func (s *MockStore) Add(task storage.Task) error {
	s.tasks = append(s.tasks, task)
	return nil
}

func (s *MockStore) List() ([]storage.Task, error) {
	return s.tasks, nil
}

func (s *MockStore) Update(task storage.Task) error {
	for i, t := range s.tasks {
		if t.ID == task.ID {
			s.tasks[i] = task
			return nil
		}
	}
	return storage.ErrTaskNotFound
}

func (s *MockStore) Delete(id string) error {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return storage.ErrTaskNotFound
}
