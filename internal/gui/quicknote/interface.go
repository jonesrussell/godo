package quicknote

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
)

// Interface defines the behavior of a quick note window
type Interface interface {
	// Initialize sets up the window with the given app and logger
	Initialize(app fyne.App, log logger.Logger)
	// Show displays the quick note window
	Show()
	// Hide hides the quick note window
	Hide()
}
