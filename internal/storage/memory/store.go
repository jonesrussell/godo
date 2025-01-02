// Package memory provides an in-memory implementation of the storage.TaskStore interface
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements the storage.Store interface using in-memory storage
type Store struct {
	mu    sync.RWMutex
	tasks map[string]storage.Task
}

// New creates a new memory store instance
func New() *Store {
	return &Store{
		tasks: make(map[string]storage.Task),
	}
}

// Add adds a new task
func (s *Store) Add(_ context.Context, task storage.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

	if _, exists := s.tasks[task.ID]; !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	s.tasks[task.ID] = task
	return nil
}

// Delete removes a task by ID
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

	s.tasks = make(map[string]storage.Task)
	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(_ context.Context) (storage.TaskTx, error) {
	return nil, fmt.Errorf("transactions not supported in memory store")
}
