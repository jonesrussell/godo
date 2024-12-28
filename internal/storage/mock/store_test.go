package mock

import (
	"testing"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"
	storagetesting "github.com/jonesrussell/godo/internal/storage/testing"
	"github.com/stretchr/testify/assert"
)

func TestMockStore(t *testing.T) {
	// Run the standard storage test suite
	suite := &storagetesting.StoreSuite{
		NewStore: func() storage.Store {
			return New()
		},
	}
	suite.Run(t)

	// Additional tests specific to mock functionality
	t.Run("ErrorSimulation", func(t *testing.T) {
		store := New()
		defer store.Close()

		expectedErr := errors.ErrTaskNotFound
		store.Error = expectedErr

		_, err := store.List()
		assert.ErrorIs(t, err, expectedErr)
		assert.True(t, store.ListCalled)

		err = store.Add(storage.Task{ID: "test"})
		assert.ErrorIs(t, err, expectedErr)
		assert.True(t, store.AddCalled)

		err = store.Update(storage.Task{ID: "test"})
		assert.ErrorIs(t, err, expectedErr)
		assert.True(t, store.UpdateCalled)

		err = store.Delete("test")
		assert.ErrorIs(t, err, expectedErr)
		assert.True(t, store.DeleteCalled)

		_, err = store.GetByID("test")
		assert.ErrorIs(t, err, expectedErr)
		assert.True(t, store.GetByIDCalled)
	})

	t.Run("Reset", func(t *testing.T) {
		store := New()
		defer store.Close()

		// Add a task and simulate an error
		store.Add(storage.Task{ID: "test"})
		store.Error = errors.ErrTaskNotFound

		// Reset the store
		store.Reset()

		// Verify the store is in its initial state
		assert.Empty(t, store.tasks)
		assert.False(t, store.closed)
		assert.False(t, store.AddCalled)
		assert.False(t, store.UpdateCalled)
		assert.False(t, store.DeleteCalled)
		assert.False(t, store.ListCalled)
		assert.False(t, store.GetByIDCalled)
		assert.Nil(t, store.Error)

		// Verify operations work after reset
		err := store.Add(storage.Task{ID: "test"})
		assert.NoError(t, err)
	})
}
