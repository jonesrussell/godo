// Package testfixtures provides shared test setup helpers (no production imports of test code).
package testfixtures

import (
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/sqlite"
)

// NewTempSQLiteStore creates a SQLite-backed store in a temporary directory using
// the same migrations as production. Pure Go driver (modernc.org/sqlite); no CGO.
func NewTempSQLiteStore(t *testing.T) *sqlite.Store {
	t.Helper()
	log := logger.NewNoopLogger()
	dbPath := filepath.Join(t.TempDir(), "notes.db")
	st, err := sqlite.New(dbPath, log)
	if err != nil {
		t.Fatalf("sqlite.New: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}
