package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
)

type QuickNote struct {
	service service.TodoServicer
	app     fyne.App
	window  fyne.Window
}

func NewQuickNote(service service.TodoServicer) *QuickNote {
	a := app.New()
	return &QuickNote{
		service: service,
		app:     a,
		window:  a.NewWindow("Quick Note"),
	}
}

func (qn *QuickNote) Update() {
	logger.Debug("Updating quick note window...")
}

func (qn *QuickNote) Show() {
	logger.Debug("Opening quick note window...")

	qn.window.Resize(fyne.NewSize(300, 100))
	qn.window.CenterOnScreen()
	qn.window.RequestFocus()

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter quick note...")
	input.OnSubmitted = func(text string) {
		if text != "" {
			todo, err := qn.service.CreateTodo(context.Background(), "quick", text)
			if err != nil {
				logger.Error("Failed to create todo: %v", err)
				return
			}
			logger.Debug("Created quick note: %s (ID: %d)", text, todo.ID)
		}
		qn.window.Close()
	}

	content := container.NewVBox(
		widget.NewLabel("🗒️ Quick Note"),
		input,
		widget.NewLabel("Press Enter to save • Esc to cancel"),
	)

	qn.window.SetContent(content)

	// Handle Escape key
	qn.window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			logger.Debug("Quick note cancelled")
			qn.window.Close()
		}
	})

	// Focus the input field
	qn.window.Canvas().Focus(input)

	qn.window.Show()
	logger.Debug("Quick note window should now be visible")
}
