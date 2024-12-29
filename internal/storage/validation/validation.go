// Package validation provides validation functions for storage operations
package validation

import (
	"errors"
	"fmt"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

var (
	// ErrInvalidID indicates that the task ID is invalid
	ErrInvalidID = errors.New("invalid task ID")
	// ErrInvalidContent indicates that the task content is invalid
	ErrInvalidContent = errors.New("invalid task content")
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
	if err := v.ValidateID(task.ID); err != nil {
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

// ValidateID validates the task ID
func (v *TaskValidator) ValidateID(id string) error {
	if id == "" {
		return fmt.Errorf("task ID cannot be empty")
	}

	if len(id) > MaxIDLength {
		return fmt.Errorf("task ID too long (max %d characters)", MaxIDLength)
	}

	return nil
}

// validateContent validates the task content
func (v *TaskValidator) validateContent(content string) error {
	if content == "" {
		return fmt.Errorf("task content cannot be empty")
	}

	if len(content) > MaxContentLength {
		return fmt.Errorf("task content too long (max %d characters)", MaxContentLength)
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

const (
	// MaxIDLength is the maximum allowed length for task IDs
	MaxIDLength = 100
	// MaxContentLength is the maximum allowed length for task content
	MaxContentLength = 1000
)

// ValidateID checks if a task ID is valid
func ValidateID(id string) error {
	if len(id) > MaxIDLength {
		return ErrInvalidID
	}
	return nil
}

// ValidateContent checks if task content is valid
func ValidateContent(content string) error {
	if len(content) > MaxContentLength {
		return ErrInvalidContent
	}
	return nil
}
