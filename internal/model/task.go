// Package model defines the core domain models for the application
package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Task represents a todo item - the core domain model
type Task struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTask creates a new Task item
func NewTask(content string) *Task {
	now := time.Now()
	return &Task{
		ID:        uuid.New().String(),
		Content:   content,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToggleDone toggles the done status of the task
func (t *Task) ToggleDone() {
	t.Done = !t.Done
	t.UpdatedAt = time.Now()
}

// UpdateContent updates the content of the task
func (t *Task) UpdateContent(content string) {
	t.Content = content
	t.UpdatedAt = time.Now()
}

// MarkDone marks the task as done
func (t *Task) MarkDone() {
	t.Done = true
	t.UpdatedAt = time.Now()
}

// MarkUndone marks the task as not done
func (t *Task) MarkUndone() {
	t.Done = false
	t.UpdatedAt = time.Now()
}

// IsValid validates the task content
func (t *Task) IsValid() error {
	if t.Content == "" {
		return &ValidationError{
			Field:   "content",
			Message: "task content cannot be empty",
		}
	}
	if len(t.Content) > 1000 {
		return &ValidationError{
			Field:   "content",
			Message: "task content cannot exceed 1000 characters",
		}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "validation error: " + e.Field + ": " + e.Message
}

var ErrTaskNotFound = errors.New("task not found")
