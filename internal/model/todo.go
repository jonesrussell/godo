package model

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a single todo item
type Todo struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTodo creates a new Todo item
func NewTodo(content string) *Todo {
	now := time.Now()
	return &Todo{
		ID:        uuid.New().String(),
		Content:   content,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
