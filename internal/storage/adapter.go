// Package storage provides interfaces and implementations for note persistence
package storage

import (
	"context"

	"github.com/jonesrussell/godo/internal/storage/types"
)

// Adapter provides a way to adapt between different store implementations
type Adapter struct {
	store types.Store
}

// NewAdapter creates a new store adapter
func NewAdapter(store types.Store) *Adapter {
	return &Adapter{store: store}
}

// Add adds a note to the store
func (a *Adapter) Add(ctx context.Context, note types.Note) error {
	return a.store.Add(ctx, note)
}

// Get retrieves a note from the store
func (a *Adapter) Get(ctx context.Context, id string) (types.Note, error) {
	return a.store.Get(ctx, id)
}

// List retrieves all notes from the store
func (a *Adapter) List(ctx context.Context) ([]types.Note, error) {
	return a.store.List(ctx)
}

// Update updates a note in the store
func (a *Adapter) Update(ctx context.Context, note types.Note) error {
	return a.store.Update(ctx, note)
}

// Delete removes a note from the store
func (a *Adapter) Delete(ctx context.Context, id string) error {
	return a.store.Delete(ctx, id)
}

// Close closes the store
func (a *Adapter) Close() error {
	return a.store.Close()
}

// BeginTx begins a new transaction
func (a *Adapter) BeginTx(ctx context.Context) (types.Transaction, error) {
	return a.store.BeginTx(ctx)
}
