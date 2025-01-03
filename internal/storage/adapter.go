// Package storage provides interfaces and implementations for note persistence
package storage

import (
	"context"
)

// LegacyStoreAdapter adapts the new Store interface to the old Store interface
type LegacyStoreAdapter struct {
	store Store
}

// NewLegacyStoreAdapter creates a new adapter
func NewLegacyStoreAdapter(store Store) *LegacyStoreAdapter {
	return &LegacyStoreAdapter{store: store}
}

// List returns all notes
func (a *LegacyStoreAdapter) List() ([]Note, error) {
	return a.store.List(context.Background())
}

// Add stores a new note
func (a *LegacyStoreAdapter) Add(note Note) error {
	return a.store.Add(context.Background(), note)
}

// Update modifies an existing note
func (a *LegacyStoreAdapter) Update(note Note) error {
	return a.store.Update(context.Background(), note)
}

// Delete removes a note by ID
func (a *LegacyStoreAdapter) Delete(id string) error {
	return a.store.Delete(context.Background(), id)
}

// Get retrieves a note by its ID
func (a *LegacyStoreAdapter) Get(id string) (Note, error) {
	return a.store.Get(context.Background(), id)
}

// Close releases any resources held by the store
func (a *LegacyStoreAdapter) Close() error {
	return a.store.Close()
}
