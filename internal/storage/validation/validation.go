// Package validation provides validation functions for storage operations
package validation

import (
	"errors"
	"fmt"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

var (
	// ErrInvalidID indicates that the note ID is invalid
	ErrInvalidID = errors.New("invalid note ID")
	// ErrInvalidContent indicates that the note content is invalid
	ErrInvalidContent = errors.New("invalid note content")
)

// NoteValidator provides validation for note operations
type NoteValidator struct {
	store storage.Store // For uniqueness checks
}

// NewNoteValidator creates a new NoteValidator
func NewNoteValidator(store storage.Store) *NoteValidator {
	return &NoteValidator{store: store}
}

// ValidateNote validates a note for creation or update
func (v *NoteValidator) ValidateNote(note storage.Note) error {
	if err := v.ValidateID(note.ID); err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	if err := v.validateContent(note.Content); err != nil {
		return fmt.Errorf("invalid content: %w", err)
	}

	if err := v.validateTimestamps(note.CreatedAt, note.UpdatedAt); err != nil {
		return fmt.Errorf("invalid timestamps: %w", err)
	}

	return nil
}

// ValidateID validates the note ID
func (v *NoteValidator) ValidateID(id string) error {
	if id == "" {
		return fmt.Errorf("note ID cannot be empty")
	}

	if len(id) > MaxIDLength {
		return fmt.Errorf("note ID too long (max %d characters)", MaxIDLength)
	}

	return nil
}

// validateContent validates the note content
func (v *NoteValidator) validateContent(content string) error {
	if content == "" {
		return fmt.Errorf("note content cannot be empty")
	}

	if len(content) > MaxContentLength {
		return fmt.Errorf("note content too long (max %d characters)", MaxContentLength)
	}

	return nil
}

// validateTimestamps validates note timestamps
func (v *NoteValidator) validateTimestamps(createdAt, updatedAt int64) error {
	now := time.Now().Unix()

	if createdAt <= 0 {
		return fmt.Errorf("created_at must be positive")
	}

	if updatedAt <= 0 {
		return fmt.Errorf("updated_at must be positive")
	}

	if createdAt > now {
		return fmt.Errorf("created_at cannot be in the future")
	}

	if updatedAt > now {
		return fmt.Errorf("updated_at cannot be in the future")
	}

	if updatedAt < createdAt {
		return fmt.Errorf("updated_at cannot be before created_at")
	}

	return nil
}

const (
	// MaxIDLength is the maximum allowed length for note IDs
	MaxIDLength = 36
	// MaxContentLength is the maximum allowed length for note content
	MaxContentLength = 1000
)

// ValidateID checks if a note ID is valid
func ValidateID(id string) error {
	if id == "" {
		return ErrInvalidID
	}

	if len(id) > MaxIDLength {
		return ErrInvalidID
	}

	return nil
}

// ValidateContent checks if note content is valid
func ValidateContent(content string) error {
	if content == "" {
		return ErrInvalidContent
	}

	if len(content) > MaxContentLength {
		return ErrInvalidContent
	}

	return nil
}

// ValidateNote validates a Note struct
func ValidateNote(note storage.Note) error {
	// Validate ID
	if note.ID == "" {
		return fmt.Errorf("note ID cannot be empty")
	}

	if len(note.ID) > MaxIDLength {
		return fmt.Errorf("note ID too long (max %d characters)", MaxIDLength)
	}

	// Validate Content
	if note.Content == "" {
		return fmt.Errorf("note content cannot be empty")
	}

	if len(note.Content) > MaxContentLength {
		return fmt.Errorf("note content too long (max %d characters)", MaxContentLength)
	}

	// Validate timestamps
	if note.CreatedAt == 0 {
		return fmt.Errorf("created_at cannot be zero")
	}

	if note.UpdatedAt == 0 {
		return fmt.Errorf("updated_at cannot be zero")
	}

	if note.UpdatedAt < note.CreatedAt {
		return fmt.Errorf("updated_at cannot be before created_at")
	}

	return nil
}
