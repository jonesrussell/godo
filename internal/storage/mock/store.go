// Package mock provides a mock implementation of the storage interface for testing
package mock

import (
	"context"
	"sync"

	"github.com/jonesrussell/godo/internal/storage"
)

// Store provides a mock implementation of storage.TaskStore
type Store struct {
	mu    sync.RWMutex
	tasks map[string]storage.Task
	Error error
}

// New creates a new mock store
func New() *Store {
	return &Store{
		tasks: make(map[string]storage.Task),
	}
}

// SetError sets the error to be returned by store operations
func (s *Store) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Error = err
}

// Reset clears all tasks and resets error state
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = make(map[string]storage.Task)
	s.Error = nil
}

// List returns all tasks
func (s *Store) List(_ context.Context) ([]storage.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return nil, s.Error
	}

	tasks := make([]storage.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByID retrieves a task by its ID
func (s *Store) GetByID(_ context.Context, id string) (storage.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Error != nil {
		return storage.Task{}, s.Error
	}

	task, exists := s.tasks[id]
	if !exists {
		return storage.Task{}, storage.ErrTaskNotFound
	}

	return task, nil
}

// Add creates a new task
func (s *Store) Add(_ context.Context, task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if _, exists := s.tasks[task.ID]; exists {
		return storage.ErrDuplicateID
	}

	s.tasks[task.ID] = task
	return nil
}

// Update modifies an existing task
func (s *Store) Update(_ context.Context, task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if _, exists := s.tasks[task.ID]; !exists {
		return storage.ErrTaskNotFound
	}

	s.tasks[task.ID] = task
	return nil
}

// Delete removes a task
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return s.Error
	}

	if _, exists := s.tasks[id]; !exists {
		return storage.ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

// Close is a no-op for the mock store
func (s *Store) Close() error {
	if s.Error != nil {
		return s.Error
	}
	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(_ context.Context) (storage.TaskTx, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Error != nil {
		return nil, s.Error
	}

	tx := &Tx{
		store: s,
		tasks: make(map[string]storage.Task),
	}

	// Copy current tasks
	for k, v := range s.tasks {
		tx.tasks[k] = v
	}

	return tx, nil
}

// Tx represents a mock transaction
type Tx struct {
	store *Store
	tasks map[string]storage.Task
}

// List returns all tasks in the transaction
func (t *Tx) List(_ context.Context) ([]storage.Task, error) {
	tasks := make([]storage.Task, 0, len(t.tasks))
	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByID retrieves a task by its ID in the transaction
func (t *Tx) GetByID(_ context.Context, id string) (storage.Task, error) {
	task, exists := t.tasks[id]
	if !exists {
		return storage.Task{}, storage.ErrTaskNotFound
	}
	return task, nil
}

// Add creates a new task in the transaction
func (t *Tx) Add(_ context.Context, task storage.Task) error {
	if _, exists := t.tasks[task.ID]; exists {
		return storage.ErrDuplicateID
	}
	t.tasks[task.ID] = task
	return nil
}

// Update modifies an existing task in the transaction
func (t *Tx) Update(_ context.Context, task storage.Task) error {
	if _, exists := t.tasks[task.ID]; !exists {
		return storage.ErrTaskNotFound
	}
	t.tasks[task.ID] = task
	return nil
}

// Delete removes a task in the transaction
func (t *Tx) Delete(_ context.Context, id string) error {
	if _, exists := t.tasks[id]; !exists {
		return storage.ErrTaskNotFound
	}
	delete(t.tasks, id)
	return nil
}

// Commit applies transaction changes
func (t *Tx) Commit() error {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	if t.store.Error != nil {
		return t.store.Error
	}

	t.store.tasks = make(map[string]storage.Task)
	for k, v := range t.tasks {
		t.store.tasks[k] = v
	}
	return nil
}

// Rollback discards transaction changes
func (t *Tx) Rollback() error {
	if t.store.Error != nil {
		return t.store.Error
	}
	return nil
}
