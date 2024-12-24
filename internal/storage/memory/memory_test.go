package memory

import (
	"testing"

	"github.com/jonesrussell/godo/internal/testutil"
)

func TestMemoryStore(t *testing.T) {
	store := New()
	testutil.RunStoreTests(t, store)
}
