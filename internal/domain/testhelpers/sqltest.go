// Package testhelpers provides test-only utilities for domain integration tests.
// It imports infrastructure storage (SQLite) and is intended for use from *_test.go
// packages only — not from production code paths.
package testhelpers

import (
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/domain/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	sqlitestore "github.com/jonesrussell/godo/internal/infrastructure/storage/sqlite"
)

// NewTempSQLiteUnified opens a modernc-backed SQLite file under t.TempDir(), runs
// migrations, and returns a UnifiedNoteStorage adapter plus a cleanup function.
func NewTempSQLiteUnified(t *testing.T, log logger.Logger) (storage.UnifiedNoteStorage, func()) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "notes.db")
	raw, err := sqlitestore.New(dbPath, log)
	if err != nil {
		t.Fatalf("sqlite.New: %v", err)
	}
	adapter := sqlitestore.NewUnifiedAdapter(raw)
	cleanup := func() {
		if err := adapter.Close(); err != nil {
			t.Logf("close adapter: %v", err)
		}
	}
	return adapter, cleanup
}
