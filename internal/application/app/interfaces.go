package app

import (
	"time"

	"fyne.io/fyne/v2"
)

//go:generate mockgen -destination=../../test/mocks/mock_app.go -package=mocks github.com/jonesrussell/godo/internal/application/app UI,Application

// UI defines the user interface operations
type UI interface {
	Show()
	Hide()
	SetContent(content fyne.CanvasObject)
	Resize(size fyne.Size)
	CenterOnScreen()
}

// Application defines the core application behavior
type Application interface {
	SetupUI() error
	Run()
	Cleanup()
	ForceKillTimeout() time.Duration
}
