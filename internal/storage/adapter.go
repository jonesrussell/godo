package storage

import (
	"context"
)

// StoreAdapter adapts the new TaskStore interface to the old Store interface
type StoreAdapter struct {
	store TaskStore
}

// NewStoreAdapter creates a new adapter
func NewStoreAdapter(store TaskStore) *StoreAdapter {
	return &StoreAdapter{store: store}
}

// List returns all tasks
func (a *StoreAdapter) List() ([]Task, error) {
	return a.store.List(context.Background())
}

// Add stores a new task
func (a *StoreAdapter) Add(task Task) error {
	return a.store.Add(context.Background(), task)
}

// Update modifies an existing task
func (a *StoreAdapter) Update(task Task) error {
	return a.store.Update(context.Background(), task)
}

// Delete removes a task by ID
func (a *StoreAdapter) Delete(id string) error {
	return a.store.Delete(context.Background(), id)
}

// GetByID retrieves a task by its ID
func (a *StoreAdapter) GetByID(id string) (*Task, error) {
	task, err := a.store.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Close releases any resources held by the store
func (a *StoreAdapter) Close() error {
	return a.store.Close()
}
