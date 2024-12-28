// Package memory provides an in-memory implementation of the storage interface
package memory

import (
	"sync"

	"github.com/jonesrussell/godo/internal/storage"
)

// Store provides an in-memory task storage implementation
type Store struct {
	tasks []storage.Task
	mu    sync.RWMutex
}

// New creates a new in-memory store
func New() *Store {
	return &Store{
		tasks: make([]storage.Task, 0),
	}
}

// Add stores a new task
func (s *Store) Add(task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
	return nil
}

// List returns all stored tasks
func (s *Store) List() ([]storage.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks, nil
}

// Update modifies an existing task
func (s *Store) Update(task storage.Task) error {
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

// Delete removes a task by ID
func (s *Store) Delete(id string) error {
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
