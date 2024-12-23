package memory

import (
	"testing"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()

	t.Run("Add and List", func(t *testing.T) {
		task := storage.Task{
			ID:    "1",
			Title: "Test Task",
		}

		err := store.Add(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, task, tasks[0])
	})

	t.Run("Update", func(t *testing.T) {
		task := storage.Task{
			ID:    "1",
			Title: "Updated Task",
		}

		err := store.Update(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", tasks[0].Title)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete("1")
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Empty(t, tasks)

		err = store.Delete("1")
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)
	})
}
