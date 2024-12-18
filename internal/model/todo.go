package model

import (
	"time"

	"github.com/jonesrussell/godo/internal/logger"
)

type Todo struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTodo creates a new Todo with default values
func NewTodo(title, description string) *Todo {
	logger.Debug("Creating new todo",
		"title", title,
		"description", description)
	now := time.Now()
	return &Todo{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Toggle switches the completed status
func (t *Todo) Toggle() {
	logger.Debug("Toggling todo completion status", "id", t.ID)
	t.Completed = !t.Completed
	t.UpdatedAt = time.Now()
}

// MarkCompleted sets the todo as completed
func (t *Todo) MarkCompleted() {
	logger.Debug("Marking todo as completed", "id", t.ID)
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// MarkIncomplete sets the todo as incomplete
func (t *Todo) MarkIncomplete() {
	logger.Debug("Marking todo as incomplete", "id", t.ID)
	t.Completed = false
	t.UpdatedAt = time.Now()
}

// String implements the Stringer interface
func (t *Todo) String() string {
	status := "☐"
	if t.Completed {
		status = "☑"
	}
	return status + " " + t.Title
}
