// Package validation provides validation functions for storage operations
package validation

import (
	"fmt"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// TaskValidator provides validation for task operations
type TaskValidator struct {
	store storage.TaskReader // For uniqueness checks
}

// NewTaskValidator creates a new TaskValidator
func NewTaskValidator(store storage.TaskReader) *TaskValidator {
	return &TaskValidator{store: store}
}

// ValidateTask validates a task for creation or update
func (v *TaskValidator) ValidateTask(task storage.Task) error {
	if err := v.validateID(task.ID); err != nil {
		return &storage.ValidationError{
			Field:   "id",
			Message: err.Error(),
		}
	}

	if err := v.validateContent(task.Content); err != nil {
		return &storage.ValidationError{
			Field:   "content",
			Message: err.Error(),
		}
	}

	if err := v.validateTimestamps(task.CreatedAt, task.UpdatedAt); err != nil {
		return &storage.ValidationError{
			Field:   "timestamps",
			Message: err.Error(),
		}
	}

	return nil
}

// validateID validates the task ID
func (v *TaskValidator) validateID(id string) error {
	if id == "" {
		return fmt.Errorf("task ID cannot be empty")
	}

	if len(id) > 100 {
		return fmt.Errorf("task ID too long (max 100 characters)")
	}

	return nil
}

// validateContent validates the task content
func (v *TaskValidator) validateContent(content string) error {
	if content == "" {
		return fmt.Errorf("task content cannot be empty")
	}

	if len(content) > 1000 {
		return fmt.Errorf("task content too long (max 1000 characters)")
	}

	return nil
}

// validateTimestamps validates task timestamps
func (v *TaskValidator) validateTimestamps(createdAt, updatedAt time.Time) error {
	if createdAt.IsZero() {
		return fmt.Errorf("created_at timestamp cannot be zero")
	}

	if updatedAt.IsZero() {
		return fmt.Errorf("updated_at timestamp cannot be zero")
	}

	if updatedAt.Before(createdAt) {
		return fmt.Errorf("updated_at cannot be before created_at")
	}

	return nil
}

// ValidateConnection validates the database connection state
func ValidateConnection(err error) error {
	if err != nil {
		return &storage.ConnectionError{
			Operation: "validate connection",
			Err:       err,
		}
	}
	return nil
}

// ValidateTransaction validates transaction state
func ValidateTransaction(err error) error {
	if err != nil {
		return &storage.TransactionError{
			Operation: "validate transaction",
			Err:       err,
		}
	}
	return nil
}
