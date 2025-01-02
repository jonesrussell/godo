package api

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreation(t *testing.T) {
	// Using correct field names and time handling
	task := storage.Task{
		Title:     "Buy groceries",
		Completed: false,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Using correct field names in assertions
	assert.Equal(t, "Buy groceries", task.Title)

	// Using correct field name for completion status
	assert.False(t, task.Completed)

	// Using proper time comparison with Unix timestamps
	assert.WithinDuration(t, time.Now(), time.Unix(task.CreatedAt, 0), time.Second)
}

func TestTaskUpdate(t *testing.T) {
	originalTask := &storage.Task{
		Title:     "Original task",
		Completed: true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	updatedTask := &storage.Task{
		Title:     "Updated task",
		Completed: false,
		UpdatedAt: time.Now().Unix(),
	}

	assert.NotEqual(t, originalTask.Title, updatedTask.Title)
	assert.NotEqual(t, originalTask.Completed, updatedTask.Completed)
}
