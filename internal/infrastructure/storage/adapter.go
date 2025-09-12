// Package storage provides adapters for backward compatibility
package storage

import (
	"context"

	"github.com/jonesrussell/godo/internal/domain/model"
	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
)

// NoteStoreAdapter adapts UnifiedNoteStorage to the old NoteStore interface
type NoteStoreAdapter struct {
	store domainstorage.UnifiedNoteStorage
}

// NewNoteStoreAdapter creates a new adapter
func NewNoteStoreAdapter(unifiedStore domainstorage.UnifiedNoteStorage) *NoteStoreAdapter {
	return &NoteStoreAdapter{
		store: unifiedStore,
	}
}

// Add creates a new note
func (a *NoteStoreAdapter) Add(ctx context.Context, note *model.Note) error {
	createdNote, err := a.store.CreateNote(ctx, note.Content)
	if err != nil {
		return err
	}
	// Copy the generated ID and timestamps back to the original note
	note.ID = createdNote.ID
	note.CreatedAt = createdNote.CreatedAt
	note.UpdatedAt = createdNote.UpdatedAt
	return nil
}

// GetByID retrieves a note by ID
func (a *NoteStoreAdapter) GetByID(ctx context.Context, id string) (model.Note, error) {
	note, err := a.store.GetNote(ctx, id)
	if err != nil {
		return model.Note{}, err
	}
	return *note, nil
}

// Update modifies an existing note
func (a *NoteStoreAdapter) Update(ctx context.Context, note *model.Note) error {
	updatedNote, err := a.store.UpdateNote(ctx, note.ID, note.Content, note.Done)
	if err != nil {
		return err
	}
	// Copy the updated timestamp back to the original note
	note.UpdatedAt = updatedNote.UpdatedAt
	return nil
}

// Delete removes a note by ID
func (a *NoteStoreAdapter) Delete(ctx context.Context, id string) error {
	return a.store.DeleteNote(ctx, id)
}

// List returns all notes
func (a *NoteStoreAdapter) List(ctx context.Context) ([]model.Note, error) {
	notes, err := a.store.GetAllNotes(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.Note, len(notes))
	for i, note := range notes {
		result[i] = *note
	}

	return result, nil
}

// Close closes the storage
func (a *NoteStoreAdapter) Close() error {
	return a.store.Close()
}
