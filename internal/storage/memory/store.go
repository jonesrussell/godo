// Package memory provides an in-memory implementation of the storage.TaskStore interface
package memory

import (
	"context"
	"sync"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements storage.TaskStore interface using in-memory storage
type Store struct {
	tasks map[string]model.Task
	mu    sync.RWMutex
}

// New creates a new memory store
func New() *Store {
	return &Store{
		tasks: make(map[string]model.Task),
	}
}

// Add adds a new task to the store
func (s *Store) Add(_ context.Context, task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; exists {
		return storage.ErrDuplicateID
	}

	s.tasks[task.ID] = *task
	return nil
}

// GetByID retrieves a task by its ID
func (s *Store) GetByID(_ context.Context, id string) (model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return model.Task{}, storage.ErrTaskNotFound
	}

	return task, nil
}

// List returns all tasks in the store
func (s *Store) List(_ context.Context) ([]model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update updates an existing task
func (s *Store) Update(_ context.Context, task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; !exists {
		return storage.ErrTaskNotFound
	}

	s.tasks[task.ID] = *task
	return nil
}

// Delete removes a task from the store
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return storage.ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

// Close implements storage.TaskStore interface
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks = make(map[string]model.Task)
	return nil
}
