// Package sqlite provides SQLite adapter for the unified storage interface
package sqlite

import (
	"context"
	"fmt"

	"github.com/jonesrussell/godo/internal/domain/model"
)

// UnifiedAdapter adapts the existing SQLite Store to the UnifiedNoteStorage interface
type UnifiedAdapter struct {
	store *Store
}

// NewUnifiedAdapter creates a new unified adapter for SQLite storage
func NewUnifiedAdapter(store *Store) *UnifiedAdapter {
	return &UnifiedAdapter{
		store: store,
	}
}

// CreateNote creates a new note
func (a *UnifiedAdapter) CreateNote(ctx context.Context, content string) (*model.Note, error) {
	note := model.NewNote(content)

	if err := note.IsValid(); err != nil {
		return nil, err
	}

	if err := a.store.Add(ctx, note); err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return note, nil
}

// GetNote retrieves a note by ID
func (a *UnifiedAdapter) GetNote(ctx context.Context, id string) (*model.Note, error) {
	note, err := a.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// GetAllNotes retrieves all notes
func (a *UnifiedAdapter) GetAllNotes(ctx context.Context) ([]*model.Note, error) {
	notes, err := a.store.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Note, len(notes))
	for i, note := range notes {
		result[i] = &note
	}

	return result, nil
}

// UpdateNote updates a note
func (a *UnifiedAdapter) UpdateNote(ctx context.Context, id string, content string, done bool) (*model.Note, error) {
	// First get the existing note
	note, err := a.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the note fields
	note.Content = content
	note.Done = done
	note.UpdateContent(content) // This also updates UpdatedAt

	// Save the updated note
	if err := a.store.Update(ctx, &note); err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return &note, nil
}

// DeleteNote deletes a note
func (a *UnifiedAdapter) DeleteNote(ctx context.Context, id string) error {
	if err := a.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	return nil
}

// ToggleDone toggles the done status of a note
func (a *UnifiedAdapter) ToggleDone(ctx context.Context, id string) (*model.Note, error) {
	note, err := a.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	note.ToggleDone()

	if err := a.store.Update(ctx, &note); err != nil {
		return nil, fmt.Errorf("failed to toggle note status: %w", err)
	}

	return &note, nil
}

// MarkDone marks a note as done
func (a *UnifiedAdapter) MarkDone(ctx context.Context, id string) (*model.Note, error) {
	note, err := a.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	note.MarkDone()

	if err := a.store.Update(ctx, &note); err != nil {
		return nil, fmt.Errorf("failed to mark note as done: %w", err)
	}

	return &note, nil
}

// MarkUndone marks a note as undone
func (a *UnifiedAdapter) MarkUndone(ctx context.Context, id string) (*model.Note, error) {
	note, err := a.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	note.MarkUndone()

	if err := a.store.Update(ctx, &note); err != nil {
		return nil, fmt.Errorf("failed to mark note as undone: %w", err)
	}

	return &note, nil
}

// Close closes the storage
func (a *UnifiedAdapter) Close() error {
	return a.store.Close()
}
