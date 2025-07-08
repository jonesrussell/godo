// Package validation provides validation utilities for storage operations
package validation

import (
	"errors"
	"time"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
)

// TaskValidator validates task data
type TaskValidator struct {
	store storage.TaskReader // For uniqueness checks
}

// NewTaskValidator creates a new task validator
func NewTaskValidator(store storage.TaskReader) *TaskValidator {
	return &TaskValidator{
		store: store,
	}
}

// ValidateTask validates a task
func (v *TaskValidator) ValidateTask(task *model.Task) error {
	if err := v.validateContent(task.Content); err != nil {
		return err
	}

	if err := v.validateTimestamps(task); err != nil {
		return err
	}

	return nil
}

// validateContent validates task content
func (v *TaskValidator) validateContent(content string) error {
	if content == "" {
		return &model.ValidationError{
			Field:   "content",
			Message: "task content cannot be empty",
		}
	}

	if len(content) > 1000 {
		return &model.ValidationError{
			Field:   "content",
			Message: "task content cannot exceed 1000 characters",
		}
	}

	return nil
}

// validateTimestamps validates task timestamps
func (v *TaskValidator) validateTimestamps(task *model.Task) error {
	now := time.Now()

	// Check if created_at is in the future
	if task.CreatedAt.After(now) {
		return &model.ValidationError{
			Field:   "created_at",
			Message: "created_at cannot be in the future",
		}
	}

	// Check if updated_at is in the future
	if task.UpdatedAt.After(now) {
		return &model.ValidationError{
			Field:   "updated_at",
			Message: "updated_at cannot be in the future",
		}
	}

	// Check if updated_at is before created_at
	if task.UpdatedAt.Before(task.CreatedAt) {
		return &model.ValidationError{
			Field:   "updated_at",
			Message: "updated_at cannot be before created_at",
		}
	}

	return nil
}

// ValidateTaskUpdate validates task update operations
func (v *TaskValidator) ValidateTaskUpdate(original, updated *model.Task) error {
	// Validate the updated task
	if err := v.ValidateTask(updated); err != nil {
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
