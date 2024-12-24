package sqlite

import (
	"testing"

	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestSQLiteStore(t *testing.T) {
	store, err := NewTestStore(t)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, store.Close())
	})

	testutil.RunStoreTests(t, store)
}

func TestSQLiteStore_Close(t *testing.T) {
	store, err := NewTestStore(t)
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)

	// Verify operations fail after close
	err = store.SaveNote("test")
	require.Error(t, err)
}
