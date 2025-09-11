package storage

import (
	"context"

	"github.com/jonesrussell/godo/internal/domain/model"
)

// StoreAdapter adapts the new NoteStore interface to the old Store interface
type StoreAdapter struct {
	store NoteStore
}

// NewStoreAdapter creates a new adapter
func NewStoreAdapter(store NoteStore) *StoreAdapter {
	return &StoreAdapter{store: store}
}

// List returns all notes
func (a *StoreAdapter) List() ([]model.Note, error) {
	return a.store.List(context.Background())
}

// Add stores a new note
func (a *StoreAdapter) Add(note *model.Note) error {
	return a.store.Add(context.Background(), note)
}

// Update modifies an existing note
func (a *StoreAdapter) Update(note *model.Note) error {
	return a.store.Update(context.Background(), note)
}

// Delete removes a note by ID
func (a *StoreAdapter) Delete(id string) error {
	return a.store.Delete(context.Background(), id)
}

// GetByID retrieves a note by its ID
func (a *StoreAdapter) GetByID(id string) (*model.Note, error) {
	note, err := a.store.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// Close releases any resources held by the store
func (a *StoreAdapter) Close() error {
	return a.store.Close()
}
