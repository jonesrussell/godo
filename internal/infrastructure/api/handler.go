// Package api implements the HTTP server and API endpoints
package api

import (
	"net/http"
	"time"

	"github.com/jonesrussell/godo/internal/domain/model"
)

// TaskHandler defines the interface for task-related HTTP handlers
type TaskHandler interface {
	// List returns all tasks
	List(w http.ResponseWriter, r *http.Request)
	// Create creates a new task
	Create(w http.ResponseWriter, r *http.Request)
	// Get returns a specific task
	Get(w http.ResponseWriter, r *http.Request)
	// Update replaces an existing task
	Update(w http.ResponseWriter, r *http.Request)
	// Patch partially updates an existing task
	Patch(w http.ResponseWriter, r *http.Request)
	// Delete removes a task
	Delete(w http.ResponseWriter, r *http.Request)
}

// CreateTaskRequest represents a request to create a new task
type CreateTaskRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
}

// UpdateTaskRequest represents a request to update an existing task
type UpdateTaskRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
	Done    bool   `json:"done"`
}

// PatchTaskRequest represents a request to partially update a task
type PatchTaskRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,max=1000"`
	Done    *bool   `json:"done,omitempty"`
}

// TaskResponse represents a task in API responses
type TaskResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTaskResponse creates a TaskResponse from a model.Task
func NewTaskResponse(task *model.Task) TaskResponse {
	return TaskResponse{
		ID:        task.ID,
		Content:   task.Content,
		Done:      task.Done,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

// TaskListResponse represents a list of tasks in API responses
type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

// NewTaskListResponse creates a TaskListResponse from a slice of model.Tasks
func NewTaskListResponse(tasks []model.Task) TaskListResponse {
	response := TaskListResponse{
		Tasks: make([]TaskResponse, len(tasks)),
	}
	for i, task := range tasks {
		response.Tasks[i] = NewTaskResponse(&task)
	}
	return response
}

// ErrorResponse represents an error in API responses
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// ValidationErrorResponse represents a validation error in API responses
type ValidationErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}
