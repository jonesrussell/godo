// Package service provides business logic services for the application
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/domain/repository"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

//go:generate mockgen -destination=../../test/mocks/mock_taskservice.go -package=mocks github.com/jonesrussell/godo/internal/domain/service TaskService

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
	CreateTask(ctx context.Context, content string) (*model.Task, error)
	GetTask(ctx context.Context, id string) (*model.Task, error)
	UpdateTask(ctx context.Context, id string, updates TaskUpdateRequest) (*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
	ListTasks(ctx context.Context, filter *TaskFilter) ([]*model.Task, error)
}

// taskService implements TaskService
type taskService struct {
	repo   repository.TaskRepository
	logger logger.Logger
}

// NewTaskService creates a new TaskService instance
func NewTaskService(repo repository.TaskRepository, log logger.Logger) TaskService {
	return &taskService{
		repo:   repo,
		logger: log,
	}
}

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

func (s *taskService) validateTaskID(id string) error {
	if strings.TrimSpace(id) == "" {
		return &model.ValidationError{
			Field:   "id",
			Message: "task ID cannot be empty",
		}
	}
	if _, err := uuid.Parse(id); err != nil {
		return &model.ValidationError{
			Field:   "id",
			Message: "invalid task ID format",
		}
	}
	return nil
}

func (s *taskService) CreateTask(ctx context.Context, content string) (*model.Task, error) {
	s.logger.Info("Creating new task", "content_length", len(content))
	if err := s.validateTaskContent(content); err != nil {
		s.logger.Error("Task content validation failed", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	task := model.Task{
		ID:        uuid.New().String(),
		Content:   strings.TrimSpace(content),
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Add(ctx, &task); err != nil {
		s.logger.Error("Failed to store task", "task_id", task.ID, "error", err)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	s.logger.Info("Task created successfully", "task_id", task.ID)
	return &task, nil
}

func (s *taskService) GetTask(ctx context.Context, id string) (*model.Task, error) {
	s.logger.Info("Retrieving task", "task_id", id)
	if err := s.validateTaskID(id); err != nil {
		s.logger.Error("Task ID validation failed", "task_id", id, "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve task", "task_id", id, "error", err)
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}
	s.logger.Info("Task retrieved successfully", "task_id", id)
	return task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, id string, updates TaskUpdateRequest) (*model.Task, error) {
	s.logger.Info("Updating task", "task_id", id, "updates", updates)
	if err := s.validateTaskID(id); err != nil {
		s.logger.Error("Task ID validation failed", "task_id", id, "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	existingTask, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve existing task", "task_id", id, "error", err)
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}
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
	existingTask.UpdatedAt = time.Now()
	if updateErr := s.repo.Update(ctx, existingTask); updateErr != nil {
		s.logger.Error("Failed to update task", "task_id", id, "error", updateErr)
		return nil, fmt.Errorf("failed to update task: %w", updateErr)
	}
	s.logger.Info("Task updated successfully", "task_id", id)
	return existingTask, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	s.logger.Info("Deleting task", "task_id", id)
	if err := s.validateTaskID(id); err != nil {
		s.logger.Error("Task ID validation failed", "task_id", id, "error", err)
		return fmt.Errorf("validation failed: %w", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete task", "task_id", id, "error", err)
		return fmt.Errorf("failed to delete task: %w", err)
	}
	s.logger.Info("Task deleted successfully", "task_id", id)
	return nil
}

func (s *taskService) ListTasks(ctx context.Context, filter *TaskFilter) ([]*model.Task, error) {
	s.logger.Info("Retrieving tasks", "filter", filter)
	tasks, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("Failed to retrieve tasks", "error", err)
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	if filter != nil {
		tasks = s.applyFilters(tasks, filter)
	}
	s.logger.Info("Tasks retrieved successfully", "count", len(tasks))
	return tasks, nil
}

// applyFilters applies the given filters to the task list
func (s *taskService) applyFilters(tasks []*model.Task, filter *TaskFilter) []*model.Task {
	if filter == nil {
		return tasks
	}
	filtered := s.filterByCriteria(tasks, filter)
	return s.applyPagination(filtered, filter)
}

// filterByCriteria applies content and date filters
func (s *taskService) filterByCriteria(tasks []*model.Task, filter *TaskFilter) []*model.Task {
	var filtered []*model.Task
	for _, task := range tasks {
		if s.matchesFilter(task, filter) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// applyPagination applies limit and offset to the filtered results
func (s *taskService) applyPagination(tasks []*model.Task, filter *TaskFilter) []*model.Task {
	if filter.Limit == nil || *filter.Limit <= 0 {
		return tasks
	}
	limit := *filter.Limit
	offset := 0
	if filter.Offset != nil && *filter.Offset > 0 {
		offset = *filter.Offset
	}
	return s.sliceWithBounds(tasks, offset, limit)
}

// sliceWithBounds safely slices the tasks array with offset and limit
func (s *taskService) sliceWithBounds(tasks []*model.Task, offset, limit int) []*model.Task {
	if offset >= len(tasks) {
		return []*model.Task{}
	}
	end := offset + limit
	if end > len(tasks) {
		end = len(tasks)
	}
	return tasks[offset:end]
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
