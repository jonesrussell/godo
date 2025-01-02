// Package validation provides validation functions for storage operations
package validation

import (
	"errors"
	"fmt"

	"github.com/jonesrussell/godo/internal/storage"
)

var (
	// ErrInvalidID indicates that the task ID is invalid
	ErrInvalidID = errors.New("invalid task ID")
	// ErrInvalidTitle indicates that the task title is invalid
	ErrInvalidTitle = errors.New("invalid task title")
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

	if err := v.validateTitle(task.Title); err != nil {
		return &storage.ValidationError{
			Field:   "title",
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

// validateTitle validates the task title
func (v *TaskValidator) validateTitle(title string) error {
	if title == "" {
		return fmt.Errorf("task title cannot be empty")
	}

	if len(title) > MaxTitleLength {
		return fmt.Errorf("task title too long (max %d characters)", MaxTitleLength)
	}

	return nil
}

// validateTimestamps validates task timestamps
func (v *TaskValidator) validateTimestamps(createdAt, updatedAt int64) error {
	if createdAt == 0 {
		return fmt.Errorf("created_at timestamp cannot be zero")
	}

	if updatedAt == 0 {
		return fmt.Errorf("updated_at timestamp cannot be zero")
	}

	if updatedAt < createdAt {
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
	// MaxTitleLength is the maximum allowed length for task titles
	MaxTitleLength = 1000
)

// ValidateID checks if a task ID is valid
func ValidateID(id string) error {
	if len(id) > MaxIDLength {
		return ErrInvalidID
	}
	return nil
}

// ValidateTitle checks if task title is valid
func ValidateTitle(title string) error {
	if len(title) > MaxTitleLength {
		return ErrInvalidTitle
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field %s: %s", e.Field, e.Message)
}

// ValidateTask validates a Task struct
func ValidateTask(task storage.Task) error {
	// Validate ID
	if task.ID == "" {
		return &ValidationError{
			Field:   "id",
			Message: "ID cannot be empty",
		}
	}

	// Validate Title
	if task.Title == "" {
		return &ValidationError{
			Field:   "title",
			Message: "Title cannot be empty",
		}
	}

	// Validate CreatedAt
	if task.CreatedAt == 0 {
		return &ValidationError{
			Field:   "created_at",
			Message: "CreatedAt must be set",
		}
	}

	// Validate UpdatedAt
	if task.UpdatedAt == 0 {
		return &ValidationError{
			Field:   "updated_at",
			Message: "UpdatedAt must be set",
		}
	}

	// Validate timestamps order
	if task.UpdatedAt < task.CreatedAt {
		return &ValidationError{
			Field:   "timestamps",
			Message: "UpdatedAt cannot be before CreatedAt",
		}
	}

	return nil
}
