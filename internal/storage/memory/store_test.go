package memory

import (
	"testing"

	"github.com/jonesrussell/godo/internal/storage"
	storagetesting "github.com/jonesrussell/godo/internal/storage/testing"
)

func TestMemoryStore(t *testing.T) {
	suite := &storagetesting.StoreSuite{
		NewStore: func() storage.Store {
			return New()
		},
	}
	suite.Run(t)
}
