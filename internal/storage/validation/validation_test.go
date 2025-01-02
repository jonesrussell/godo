package validation

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestValidateTask(t *testing.T) {
	t.Run("valid task", func(t *testing.T) {
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := ValidateTask(task)
		assert.NoError(t, err)
	})

	t.Run("empty ID", func(t *testing.T) {
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := ValidateTask(task)
		assert.Error(t, err)
	})

	t.Run("empty title", func(t *testing.T) {
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := ValidateTask(task)
		assert.Error(t, err)
	})

	t.Run("zero created at", func(t *testing.T) {
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: 0,
			UpdatedAt: now,
		}

		err := ValidateTask(task)
		assert.Error(t, err)
	})

	t.Run("zero updated at", func(t *testing.T) {
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: 0,
		}

		err := ValidateTask(task)
		assert.Error(t, err)
	})

	t.Run("updated at before created at", func(t *testing.T) {
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now - 1,
		}

		err := ValidateTask(task)
		assert.Error(t, err)
	})
}
