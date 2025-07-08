// Package service provides business logic services for the application
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
)

//go:generate mockgen -destination=../../test/mocks/mock_taskservice.go -package=mocks github.com/jonesrussell/godo/internal/service TaskService

// TaskFilter represents filtering options for task queries
type TaskFilter struct {
	Done          *bool      `json:"done,omitempty"`
	Content       *string    `json:"content,omitempty"`
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	Limit         *int       `json:"limit,omitempty"`
	Offset        *int       `json:"offset,omitempty"`
}

// TaskUpdateRequest represents a request to update a task
type TaskUpdateRequest struct {
	Content *string `json:"content,omitempty"`
	Done    *bool   `json:"done,omitempty"`
}

// TaskService defines the interface for task business logic operations
type TaskService interface {
	// CreateTask creates a new task with validation and business rules
	CreateTask(ctx context.Context, content string) (*model.Task, error)

	// GetTask retrieves a task by ID with proper error handling
	GetTask(ctx context.Context, id string) (*model.Task, error)

	// UpdateTask updates a task with validation and business rules
	UpdateTask(ctx context.Context, id string, updates TaskUpdateRequest) (*model.Task, error)

	// DeleteTask deletes a task with proper cleanup
	DeleteTask(ctx context.Context, id string) error

	// ListTasks retrieves tasks with optional filtering
	ListTasks(ctx context.Context, filter *TaskFilter) ([]*model.Task, error)
}

// taskService implements TaskService
type taskService struct {
	store  storage.TaskStore
	logger logger.Logger
}

// NewTaskService creates a new TaskService instance
func NewTaskService(store storage.TaskStore, log logger.Logger) TaskService {
	return &taskService{
		store:  store,
		logger: log,
	}
}

// validateTaskContent validates task content according to business rules
func (s *taskService) validateTaskContent(content string) error {
	if strings.TrimSpace(content) == "" {
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

// validateTaskID validates task ID format
func (s *taskService) validateTaskID(id string) error {
	if strings.TrimSpace(id) == "" {
		return &model.ValidationError{
			Field:   "id",
			Message: "task ID cannot be empty",
		}
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return &model.ValidationError{
			Field:   "id",
			Message: "invalid task ID format",
		}
	}

	return nil
}

// CreateTask creates a new task with validation and business rules
func (s *taskService) CreateTask(ctx context.Context, content string) (*model.Task, error) {
	s.logger.Info("Creating new task", "content_length", len(content))

	// Validate content
	if err := s.validateTaskContent(content); err != nil {
		s.logger.Error("Task content validation failed", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create task with generated ID and timestamps
	task := model.Task{
		ID:        uuid.New().String(),
		Content:   strings.TrimSpace(content),
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the task
	if err := s.store.Add(ctx, &task); err != nil {
		s.logger.Error("Failed to store task", "task_id", task.ID, "error", err)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	s.logger.Info("Task created successfully", "task_id", task.ID)
	return &task, nil
}

// GetTask retrieves a task by ID with proper error handling
func (s *taskService) GetTask(ctx context.Context, id string) (*model.Task, error) {
	s.logger.Info("Retrieving task", "task_id", id)

	// Validate ID
	if err := s.validateTaskID(id); err != nil {
		s.logger.Error("Task ID validation failed", "task_id", id, "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Retrieve task from storage
	task, err := s.store.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve task", "task_id", id, "error", err)
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}

	s.logger.Info("Task retrieved successfully", "task_id", id)
	return &task, nil
}

// UpdateTask updates a task with validation and business rules
func (s *taskService) UpdateTask(ctx context.Context, id string, updates TaskUpdateRequest) (*model.Task, error) {
	s.logger.Info("Updating task", "task_id", id, "updates", updates)

	// Validate ID
	if err := s.validateTaskID(id); err != nil {
		s.logger.Error("Task ID validation failed", "task_id", id, "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing task
	existingTask, err := s.store.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve existing task", "task_id", id, "error", err)
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}

	// Apply updates
	if updates.Content != nil {
		if validErr := s.validateTaskContent(*updates.Content); validErr != nil {
			s.logger.Error("Task content validation failed", "task_id", id, "error", validErr)
			return nil, fmt.Errorf("validation failed: %w", validErr)
		}
		existingTask.Content = strings.TrimSpace(*updates.Content)
	}

	if updates.Done != nil {
		existingTask.Done = *updates.Done
	}

	// Update timestamp
	existingTask.UpdatedAt = time.Now()

	// Store updated task
	if updateErr := s.store.Update(ctx, &existingTask); updateErr != nil {
		s.logger.Error("Failed to update task", "task_id", id, "error", updateErr)
		return nil, fmt.Errorf("failed to update task: %w", updateErr)
	}

	s.logger.Info("Task updated successfully", "task_id", id)
	return &existingTask, nil
}

// DeleteTask deletes a task with proper cleanup
func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	s.logger.Info("Deleting task", "task_id", id)

	// Validate ID
	if err := s.validateTaskID(id); err != nil {
		s.logger.Error("Task ID validation failed", "task_id", id, "error", err)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Delete task from storage
	if err := s.store.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete task", "task_id", id, "error", err)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	s.logger.Info("Task deleted successfully", "task_id", id)
	return nil
}

// ListTasks retrieves tasks with optional filtering
func (s *taskService) ListTasks(ctx context.Context, filter *TaskFilter) ([]*model.Task, error) {
	s.logger.Info("Retrieving tasks", "filter", filter)

	// Get all tasks from storage
	tasks, err := s.store.List(ctx)
	if err != nil {
		s.logger.Error("Failed to retrieve tasks", "error", err)
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	// Convert to pointers for filtering
	taskPtrs := make([]*model.Task, len(tasks))
	for i := range tasks {
		taskPtrs[i] = &tasks[i]
	}

	// Apply filters if provided
	if filter != nil {
		taskPtrs = s.applyFilters(taskPtrs, filter)
	}

	s.logger.Info("Tasks retrieved successfully", "count", len(taskPtrs))
	return taskPtrs, nil
}

// applyFilters applies the given filters to the task list
func (s *taskService) applyFilters(tasks []*model.Task, filter *TaskFilter) []*model.Task {
	if filter == nil {
		return tasks
	}

	var filtered []*model.Task
	for _, task := range tasks {
		if s.matchesFilter(task, filter) {
			filtered = append(filtered, task)
		}
	}

	// Apply limit and offset
	if filter.Limit != nil && *filter.Limit > 0 {
		limit := *filter.Limit
		if filter.Offset != nil && *filter.Offset > 0 {
			offset := *filter.Offset
			if offset < len(filtered) {
				if offset+limit > len(filtered) {
					limit = len(filtered) - offset
				}
				filtered = filtered[offset : offset+limit]
			} else {
				filtered = filtered[:0]
			}
		} else {
			if limit > len(filtered) {
				limit = len(filtered)
			}
			filtered = filtered[:limit]
		}
	}

	return filtered
}

// matchesFilter checks if a task matches the given filter criteria
func (s *taskService) matchesFilter(task *model.Task, filter *TaskFilter) bool {
	if filter.Done != nil && task.Done != *filter.Done {
		return false
	}

	if filter.Content != nil && !strings.Contains(strings.ToLower(task.Content), strings.ToLower(*filter.Content)) {
		return false
	}

	if filter.CreatedAfter != nil && task.CreatedAt.Before(*filter.CreatedAfter) {
		return false
	}

	if filter.CreatedBefore != nil && task.CreatedAt.After(*filter.CreatedBefore) {
		return false
	}

	return true
}
