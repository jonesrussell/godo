package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
)

// NewTestStore creates a new SQLite store for testing
func NewTestStore(t *testing.T) (*Store, error) {
	t.Helper()

	// Create test database in temp directory
	dbPath := filepath.Join(t.TempDir(), "test.db")

	// Use test logger
	log := logger.NewTestLogger(t)

	return New(dbPath, log)
}
