package repository

import (
	"context"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/domain/storage"
)

type NoteRepository interface {
	Add(ctx context.Context, note *model.Note) error
	GetByID(ctx context.Context, id string) (*model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Note, error)
}

type noteRepository struct {
	store storage.UnifiedNoteStorage
}

func NewNoteRepository(store storage.UnifiedNoteStorage) NoteRepository {
	return &noteRepository{store: store}
}

func (r *noteRepository) Add(ctx context.Context, note *model.Note) error {
	if err := note.IsValid(); err != nil {
		return err
	}
	// Use the unified storage CreateNote method
	createdNote, err := r.store.CreateNote(ctx, note.Content)
	if err != nil {
		return err
	}
	// Copy the generated ID and timestamps back to the original note
	note.ID = createdNote.ID
	note.CreatedAt = createdNote.CreatedAt
	note.UpdatedAt = createdNote.UpdatedAt
	return nil
}

func (r *noteRepository) GetByID(ctx context.Context, id string) (*model.Note, error) {
	note, err := r.store.GetNote(ctx, id)
	if err != nil {
		return nil, mapStorageError(err)
	}
	return note, nil
}

func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
	if err := note.IsValid(); err != nil {
		return err
	}
	// Use the unified storage UpdateNote method
	updatedNote, err := r.store.UpdateNote(ctx, note.ID, note.Content, note.Done)
	if err != nil {
		return err
	}
	// Copy the updated timestamp back to the original note
	note.UpdatedAt = updatedNote.UpdatedAt
	return nil
}

func (r *noteRepository) Delete(ctx context.Context, id string) error {
	return r.store.DeleteNote(ctx, id)
}

func (r *noteRepository) List(ctx context.Context) ([]*model.Note, error) {
	return r.store.GetAllNotes(ctx)
}

// mapStorageError maps storage errors to domain errors (expand as needed)
func mapStorageError(err error) error {
	// Check for NotFoundError from infrastructure layer
	if err != nil && err.Error() != "" {
		// Simple string matching for now - could be improved with error types
		if err.Error() == "note not found: "+err.Error() {
			return model.ErrNoteNotFound
		}
	}
	return err
}
