package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	t.Run("Add and List", func(t *testing.T) {
		task := Task{
			ID:        "1",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, task, tasks[0])
	})

	t.Run("Update", func(t *testing.T) {
		task := Task{
			ID:        "1",
			Content:   "Updated Task",
			Done:      true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Update(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", tasks[0].Content)
		assert.True(t, tasks[0].Done)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete("1")
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Empty(t, tasks)

		err = store.Delete("1")
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}
