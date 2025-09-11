package repository

import (
	"context"
	"errors"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

type NoteRepository interface {
	Add(ctx context.Context, note *model.Note) error
	GetByID(ctx context.Context, id string) (*model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Note, error)
}

type noteRepository struct {
	store storage.NoteStore
}

func NewNoteRepository(store storage.NoteStore) NoteRepository {
	return &noteRepository{store: store}
}

func (r *noteRepository) Add(ctx context.Context, note *model.Note) error {
	if err := note.IsValid(); err != nil {
		return err
	}
	return r.store.Add(ctx, note)
}

func (r *noteRepository) GetByID(ctx context.Context, id string) (*model.Note, error) {
	note, err := r.store.GetByID(ctx, id)
	if err != nil {
		return nil, mapStorageError(err)
	}
	return &note, nil
}

func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
	if err := note.IsValid(); err != nil {
		return err
	}
	return r.store.Update(ctx, note)
}

func (r *noteRepository) Delete(ctx context.Context, id string) error {
	return r.store.Delete(ctx, id)
}

func (r *noteRepository) List(ctx context.Context) ([]*model.Note, error) {
	notes, err := r.store.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*model.Note, len(notes))
	for i := range notes {
		result[i] = &notes[i]
	}
	return result, nil
}

// mapStorageError maps storage errors to domain errors (expand as needed)
func mapStorageError(err error) error {
	if errors.Is(err, storage.ErrNoteNotFound) {
		return model.ErrNoteNotFound
	}
	return err
}
