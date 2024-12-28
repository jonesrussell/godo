// Package mock provides a mock implementation of the storage interface for testing
package mock

import (
	"context"
	"sync"

	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements storage.TaskStore for testing
type Store struct {
	tasks map[string]storage.Task
	mu    sync.RWMutex
	err   error
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
	s.err = err
}

// Reset clears all tasks and errors
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = make(map[string]storage.Task)
	s.err = nil
}

// List returns all tasks
func (s *Store) List(ctx context.Context) ([]storage.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.err != nil {
		return nil, s.err
	}

	tasks := make([]storage.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByID returns a task by its ID
func (s *Store) GetByID(ctx context.Context, id string) (*storage.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.err != nil {
		return nil, s.err
	}

	task, ok := s.tasks[id]
	if !ok {
		return nil, storage.ErrTaskNotFound
	}
	return &task, nil
}

// Add creates a new task
func (s *Store) Add(ctx context.Context, task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.tasks[task.ID]; exists {
		return storage.ErrDuplicateID
	}

	s.tasks[task.ID] = task
	return nil
}

// Update replaces an existing task
func (s *Store) Update(ctx context.Context, task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.tasks[task.ID]; !exists {
		return storage.ErrTaskNotFound
	}

	s.tasks[task.ID] = task
	return nil
}

// Delete removes a task
func (s *Store) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.tasks[id]; !exists {
		return storage.ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

// Close implements io.Closer
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	s.tasks = nil
	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(ctx context.Context) (storage.TaskTx, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return nil, s.err
	}

	// Create a snapshot of current tasks
	tasks := make(map[string]storage.Task)
	for k, v := range s.tasks {
		tasks[k] = v
	}

	return &Tx{
		store:    s,
		tasks:    tasks,
		original: s.tasks,
	}, nil
}

// Tx implements storage.TaskTx for testing
type Tx struct {
	store    *Store
	tasks    map[string]storage.Task
	original map[string]storage.Task
	mu       sync.RWMutex
}

// List returns all tasks within the transaction
func (t *Tx) List(ctx context.Context) ([]storage.Task, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tasks := make([]storage.Task, 0, len(t.tasks))
	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByID returns a task by its ID within the transaction
func (t *Tx) GetByID(ctx context.Context, id string) (*storage.Task, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	task, ok := t.tasks[id]
	if !ok {
		return nil, storage.ErrTaskNotFound
	}
	return &task, nil
}

// Add creates a new task within the transaction
func (t *Tx) Add(ctx context.Context, task storage.Task) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.tasks[task.ID]; exists {
		return storage.ErrDuplicateID
	}

	t.tasks[task.ID] = task
	return nil
}

// Update replaces an existing task within the transaction
func (t *Tx) Update(ctx context.Context, task storage.Task) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.tasks[task.ID]; !exists {
		return storage.ErrTaskNotFound
	}

	t.tasks[task.ID] = task
	return nil
}

// Delete removes a task within the transaction
func (t *Tx) Delete(ctx context.Context, id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.tasks[id]; !exists {
		return storage.ErrTaskNotFound
	}

	delete(t.tasks, id)
	return nil
}

// Commit commits the transaction
func (t *Tx) Commit() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	// Update store's tasks with transaction's tasks
	t.store.tasks = make(map[string]storage.Task)
	for k, v := range t.tasks {
		t.store.tasks[k] = v
	}

	return nil
}

// Rollback rolls back the transaction
func (t *Tx) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Restore original tasks
	t.tasks = make(map[string]storage.Task)
	for k, v := range t.original {
		t.tasks[k] = v
	}

	return nil
}
