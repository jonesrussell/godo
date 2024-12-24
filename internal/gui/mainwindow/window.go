package mainwindow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/gui/mainwindow/systray"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window represents the main application window
type Window struct {
	app     fyne.App
	win     fyne.Window
	store   storage.Store
	systray systray.Interface
	log     logger.Logger
}

// New creates a new main window
func New(store storage.Store, log logger.Logger) *Window {
	fyneApp := app.New()
	win := fyneApp.NewWindow("Godo")

	w := &Window{
		app:   fyneApp,
		win:   win,
		store: store,
		log:   log,
	}

	// Set close handler to hide window instead of closing app
	win.SetCloseIntercept(func() {
		w.Hide()
	})

	// Setup system tray
	w.systray = systray.New(fyneApp, log)
	w.systray.Setup(fyne.NewMenu("Godo",
		fyne.NewMenuItem("Show", func() {
			w.Show()
		}),
		fyne.NewMenuItem("Quit", func() {
			w.app.Quit()
		}),
	))
	w.systray.SetIcon(theme.HomeIcon())

	return w
}

// Show displays the main window
func (w *Window) Show() {
	w.win.Show()
	w.win.RequestFocus()
}

// Hide hides the main window
func (w *Window) Hide() {
	w.win.Hide()
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

	// Show window initially
	w.Show()
}

// Run starts the application main loop
func (w *Window) Run() {
	w.app.Run()
}
