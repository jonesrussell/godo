// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"context"
)

// LegacyStoreAdapter adapts the new TaskStore interface to the old Store interface
type LegacyStoreAdapter struct {
	store TaskStore
}

// NewLegacyStoreAdapter creates a new adapter
func NewLegacyStoreAdapter(store TaskStore) *LegacyStoreAdapter {
	return &LegacyStoreAdapter{store: store}
}

// List returns all tasks
func (a *LegacyStoreAdapter) List() ([]Task, error) {
	return a.store.List(context.Background())
}

// Add stores a new task
func (a *LegacyStoreAdapter) Add(task Task) error {
	return a.store.Add(context.Background(), task)
}

// Update modifies an existing task
func (a *LegacyStoreAdapter) Update(task Task) error {
	return a.store.Update(context.Background(), task)
}

// Delete removes a task by ID
func (a *LegacyStoreAdapter) Delete(id string) error {
	return a.store.Delete(context.Background(), id)
}

// Get retrieves a task by its ID
func (a *LegacyStoreAdapter) Get(id string) (Task, error) {
	return a.store.Get(context.Background(), id)
}

// Close releases any resources held by the store
func (a *LegacyStoreAdapter) Close() error {
	return a.store.Close()
}
