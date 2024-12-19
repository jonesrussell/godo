package memory

import (
	"testing"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	store := New()
	storage.TestStore(t, store)
}

func TestMemoryStore_Concurrency(t *testing.T) {
	store := New()
	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(i int) {
			todo := model.NewTodo("Concurrent todo")
			err := store.Add(todo)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify count
	todos := store.List()
	assert.Len(t, todos, 10)
}
