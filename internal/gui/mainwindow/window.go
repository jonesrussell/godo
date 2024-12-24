package mainwindow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window represents the main application window
type Window struct {
	app   fyne.App
	win   fyne.Window
	store storage.Store
}

// New creates a new main window
func New(store storage.Store) *Window {
	fyneApp := app.New()
	win := fyneApp.NewWindow("Godo")

	return &Window{
		app:   fyneApp,
		win:   win,
		store: store,
	}
}

// Show displays the main window
func (w *Window) Show() {
	w.win.Show()
	w.app.Run()
}

// Setup initializes the window
func (w *Window) Setup() {
	content := container.NewVBox(
		widget.NewLabel("Welcome to Godo!"),
		widget.NewButton("Add Task", func() {
			// TODO: Implement add task
		}),
	)

	w.win.SetContent(content)
	w.win.Resize(fyne.NewSize(800, 600))
	w.win.CenterOnScreen()
}
