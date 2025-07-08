package repository

import (
	"context"
	"errors"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
)

type TaskRepository interface {
	Add(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Task, error)
}

type taskRepository struct {
	store storage.TaskStore
}

func NewTaskRepository(store storage.TaskStore) TaskRepository {
	return &taskRepository{store: store}
}

func (r *taskRepository) Add(ctx context.Context, task *model.Task) error {
	if err := task.IsValid(); err != nil {
		return err
	}
	return r.store.Add(ctx, task)
}

func (r *taskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := r.store.GetByID(ctx, id)
	if err != nil {
		return nil, mapStorageError(err)
	}
	return &task, nil
}

func (r *taskRepository) Update(ctx context.Context, task *model.Task) error {
	if err := task.IsValid(); err != nil {
		return err
	}
	return r.store.Update(ctx, task)
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	return r.store.Delete(ctx, id)
}

func (r *taskRepository) List(ctx context.Context) ([]*model.Task, error) {
	tasks, err := r.store.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*model.Task, len(tasks))
	for i := range tasks {
		result[i] = &tasks[i]
	}
	return result, nil
}

// mapStorageError maps storage errors to domain errors (expand as needed)
func mapStorageError(err error) error {
	if errors.Is(err, storage.ErrTaskNotFound) {
		return model.ErrTaskNotFound
	}
	return err
}
