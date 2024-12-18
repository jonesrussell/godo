package service

import (
	"context"
	"errors"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/repository"
)

var (
	ErrEmptyTitle = errors.New("todo title cannot be empty")
	ErrNotFound   = errors.New("todo not found")
)

type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(ctx context.Context, title, description string) (*model.Todo, error) {
	logger.Debug("Creating new todo",
		"title", title,
		"description", description)

	todo := model.NewTodo(title, description)

	if err := s.repo.Create(ctx, todo); err != nil {
		logger.Error("Failed to create todo",
			"title", title,
			"error", err)
		return nil, err
	}

	logger.Info("Successfully created todo",
		"id", todo.ID,
		"title", todo.Title)
	return todo, nil
}

func (s *TodoService) GetTodo(ctx context.Context, id int64) (*model.Todo, error) {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, ErrNotFound
	}
	return todo, nil
}

func (s *TodoService) ListTodos(ctx context.Context) ([]model.Todo, error) {
	return s.repo.List(ctx)
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo *model.Todo) error {
	if todo.Title == "" {
		return ErrEmptyTitle
	}

	existing, err := s.repo.GetByID(ctx, todo.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	return s.repo.Update(ctx, todo)
}

func (s *TodoService) ToggleTodoStatus(ctx context.Context, id int64) error {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if todo == nil {
		return ErrNotFound
	}

	todo.Completed = !todo.Completed
	return s.repo.Update(ctx, todo)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
