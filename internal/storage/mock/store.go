// Package mock provides a mock implementation of the storage interface for testing
package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements the storage.Store interface for testing
type Store struct {
	mu    sync.RWMutex
	tasks map[string]storage.Task
	err   error
}

// Transaction represents a mock transaction
type Transaction struct {
	store     *Store
	tasks     map[string]storage.Task
	committed bool
}

// New creates a new mock store instance
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

// BeginTx starts a new transaction
func (s *Store) BeginTx(_ context.Context) (storage.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.err != nil {
		return nil, s.err
	}

	// Create a copy of the tasks map for the transaction
	tasks := make(map[string]storage.Task)
	for k, v := range s.tasks {
		tasks[k] = v
	}

	return &Transaction{
		store:     s,
		tasks:     tasks,
		committed: false,
	}, nil
}

// Add adds a new task
func (s *Store) Add(_ context.Context, task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	s.tasks[task.ID] = task
	return nil
}

// Get retrieves a task by ID
func (s *Store) Get(_ context.Context, id string) (storage.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.err != nil {
		return storage.Task{}, s.err
	}

	task, exists := s.tasks[id]
	if !exists {
		return storage.Task{}, fmt.Errorf("task with ID %s not found", id)
	}

	return task, nil
}

// List retrieves all tasks
func (s *Store) List(_ context.Context) ([]storage.Task, error) {
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

// Update updates an existing task
func (s *Store) Update(_ context.Context, task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.tasks[task.ID]; !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	task.UpdatedAt = time.Now().Unix()
	s.tasks[task.ID] = task
	return nil
}

// Delete removes a task by ID
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	delete(s.tasks, id)
	return nil
}

// Close cleans up any resources
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

	s.tasks = make(map[string]storage.Task)
	return nil
}

// Transaction methods

// Add adds a new task in the transaction
func (tx *Transaction) Add(_ context.Context, task storage.Task) error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	if _, exists := tx.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	tx.tasks[task.ID] = task
	return nil
}

// Get retrieves a task by ID in the transaction
func (tx *Transaction) Get(_ context.Context, id string) (storage.Task, error) {
	if tx.committed {
		return storage.Task{}, fmt.Errorf("transaction already committed")
	}

	task, exists := tx.tasks[id]
	if !exists {
		return storage.Task{}, fmt.Errorf("task with ID %s not found", id)
	}

	return task, nil
}

// List retrieves all tasks in the transaction
func (tx *Transaction) List(_ context.Context) ([]storage.Task, error) {
	if tx.committed {
		return nil, fmt.Errorf("transaction already committed")
	}

	tasks := make([]storage.Task, 0, len(tx.tasks))
	for _, task := range tx.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update updates an existing task in the transaction
func (tx *Transaction) Update(_ context.Context, task storage.Task) error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	if _, exists := tx.tasks[task.ID]; !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	task.UpdatedAt = time.Now().Unix()
	tx.tasks[task.ID] = task
	return nil
}

// Delete removes a task by ID in the transaction
func (tx *Transaction) Delete(_ context.Context, id string) error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	if _, exists := tx.tasks[id]; !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	delete(tx.tasks, id)
	return nil
}

// Commit commits the transaction
func (tx *Transaction) Commit() error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	tx.store.mu.Lock()
	defer tx.store.mu.Unlock()

	// Replace the store's tasks with the transaction's tasks
	tx.store.tasks = tx.tasks
	tx.committed = true
	return nil
}

// Rollback aborts the transaction
func (tx *Transaction) Rollback() error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	tx.committed = true
	return nil
}
