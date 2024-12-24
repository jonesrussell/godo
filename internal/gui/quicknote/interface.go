package quicknote

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Interface defines the common behavior for quick note windows
type Interface interface {
	Initialize(app fyne.App, log logger.Logger)
	Show()
	Hide()
}

// New creates a new quick note window based on build tags
func New(store storage.Store) Interface {
	return newWindow(store)
}
