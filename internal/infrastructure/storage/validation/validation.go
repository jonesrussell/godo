// Package validation provides validation utilities for storage operations
package validation

import (
	"errors"
	"time"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

// NoteValidator validates note data
type NoteValidator struct {
	store storage.NoteReader // For uniqueness checks
}

// NewNoteValidator creates a new note validator
func NewNoteValidator(store storage.NoteReader) *NoteValidator {
	return &NoteValidator{
		store: store,
	}
}

// ValidateNote validates a note
func (v *NoteValidator) ValidateNote(note *model.Note) error {
	if err := v.validateContent(note.Content); err != nil {
		return err
	}

	if err := v.validateTimestamps(note); err != nil {
		return err
	}

	return nil
}

// validateContent validates note content
func (v *NoteValidator) validateContent(content string) error {
	if content == "" {
		return &model.ValidationError{
			Field:   "content",
			Message: "note content cannot be empty",
		}
	}

	if len(content) > 1000 {
		return &model.ValidationError{
			Field:   "content",
			Message: "note content cannot exceed 1000 characters",
		}
	}

	return nil
}

// validateTimestamps validates note timestamps
func (v *NoteValidator) validateTimestamps(note *model.Note) error {
	now := time.Now()

	// Check if created_at is in the future
	if note.CreatedAt.After(now) {
		return &model.ValidationError{
			Field:   "created_at",
			Message: "created_at cannot be in the future",
		}
	}

	// Check if updated_at is in the future
	if note.UpdatedAt.After(now) {
		return &model.ValidationError{
			Field:   "updated_at",
			Message: "updated_at cannot be in the future",
		}
	}

	// Check if updated_at is before created_at
	if note.UpdatedAt.Before(note.CreatedAt) {
		return &model.ValidationError{
			Field:   "updated_at",
			Message: "updated_at cannot be before created_at",
		}
	}

	return nil
}

// ValidateNoteUpdate validates note update operations
func (v *NoteValidator) ValidateNoteUpdate(original, updated *model.Note) error {
	// Validate the updated note
	if err := v.ValidateNote(updated); err != nil {
		return err
	}

	// Check if ID changed
	if original.ID != updated.ID {
		return &model.ValidationError{
			Field:   "id",
			Message: "task ID cannot be changed",
		}
	}

	// Check if created_at changed
	if !original.CreatedAt.Equal(updated.CreatedAt) {
		return &model.ValidationError{
			Field:   "created_at",
			Message: "created_at cannot be modified",
		}
	}

	return nil
}

// Common validation errors
var (
	ErrEmptyContent     = errors.New("task content cannot be empty")
	ErrContentTooLong   = errors.New("task content cannot exceed 1000 characters")
	ErrInvalidTimestamp = errors.New("invalid timestamp")
)
