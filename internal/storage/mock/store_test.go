package mock

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestStore_Add(t *testing.T) {
	store := New()
	ctx := context.Background()

	// Test adding a new task
	task := storage.Task{
		ID:          "test-1",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err := store.Add(ctx, task)
	assert.NoError(t, err, "Should add task without error")

	// Test adding a duplicate task
	err = store.Add(ctx, task)
	assert.Error(t, err, "Should fail to add duplicate task")
}

func TestStore_Get(t *testing.T) {
	store := New()
	ctx := context.Background()

	// Test getting a non-existent task
	_, err := store.Get(ctx, "nonexistent")
	assert.Error(t, err, "Should fail to get non-existent task")

	// Add a task and get it
	task := storage.Task{
		ID:          "test-1",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = store.Add(ctx, task)
	assert.NoError(t, err, "Should add task without error")

	retrieved, err := store.Get(ctx, task.ID)
	assert.NoError(t, err, "Should get task without error")
	assert.Equal(t, task.Title, retrieved.Title, "Task content should match")
	assert.Equal(t, task.Completed, retrieved.Completed, "Task completion status should match")
}

func TestStore_List(t *testing.T) {
	store := New()
	ctx := context.Background()

	// Test empty list
	tasks, err := store.List(ctx)
	assert.NoError(t, err, "Should list tasks without error")
	assert.Empty(t, tasks, "Task list should be empty")

	// Add a task and list again
	task := storage.Task{
		ID:          "test-1",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = store.Add(ctx, task)
	assert.NoError(t, err, "Should add task without error")

	tasks, err = store.List(ctx)
	assert.NoError(t, err, "Should list tasks without error")
	assert.Len(t, tasks, 1, "Task list should have one item")
	assert.Equal(t, task.Title, tasks[0].Title, "Task content should match")
}

func TestStore_Update(t *testing.T) {
	store := New()
	ctx := context.Background()

	// Test updating a non-existent task
	task := storage.Task{
		ID:          "nonexistent",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err := store.Update(ctx, task)
	assert.Error(t, err, "Should fail to update non-existent task")

	// Add a task and update it
	err = store.Add(ctx, task)
	assert.NoError(t, err, "Should add task without error")

	task.Title = "Updated Task"
	task.Completed = true

	err = store.Update(ctx, task)
	assert.NoError(t, err, "Should update task without error")

	updated, err := store.Get(ctx, task.ID)
	assert.NoError(t, err, "Should get updated task without error")
	assert.Equal(t, task.Title, updated.Title, "Task content should be updated")
	assert.Equal(t, task.Completed, updated.Completed, "Task completion status should be updated")
}

func TestStore_Delete(t *testing.T) {
	store := New()
	ctx := context.Background()

	// Test deleting a non-existent task
	err := store.Delete(ctx, "nonexistent")
	assert.Error(t, err, "Should fail to delete non-existent task")

	// Add a task and delete it
	task := storage.Task{
		ID:          "test-1",
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = store.Add(ctx, task)
	assert.NoError(t, err, "Should add task without error")

	err = store.Delete(ctx, task.ID)
	assert.NoError(t, err, "Should delete task without error")

	_, err = store.Get(ctx, task.ID)
	assert.Error(t, err, "Should fail to get deleted task")
}

func TestStore_BeginTx(t *testing.T) {
	store := New()
	ctx := context.Background()

	// Test beginning a transaction
	tx, err := store.BeginTx(ctx)
	assert.Error(t, err, "Should fail to begin transaction")
	assert.Nil(t, tx, "Transaction should be nil")
}
