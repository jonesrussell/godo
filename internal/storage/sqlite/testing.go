package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// NewTestStore creates a new SQLite store for testing
func NewTestStore(t *testing.T) (*Store, func(), error) {
	t.Helper()

	logger, _ := zap.NewDevelopment()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	store, err := New(dbPath, logger)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		store.Close()
		os.Remove(dbPath)
	}

	return store, cleanup, nil
}

// NewTestDB creates a new in-memory SQLite database for testing
func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
