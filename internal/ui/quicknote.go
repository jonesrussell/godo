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
}

func NewQuickNote(service service.TodoServicer) *QuickNote {
	return &QuickNote{
		service: service,
	}
}

func (qn *QuickNote) Update() {
	logger.Debug("Updating quick note window...")
}

func (qn *QuickNote) Show() {
	logger.Debug("Opening quick note window...")

	a := app.New()
	w := a.NewWindow("Quick Note")
	w.Resize(fyne.NewSize(300, 100))
	w.CenterOnScreen()

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
		w.Close()
	}

	content := container.NewVBox(
		widget.NewLabel("🗒️ Quick Note"),
		input,
		widget.NewLabel("Press Enter to save • Esc to cancel"),
	)

	w.SetContent(content)

	// Handle Escape key
	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			logger.Debug("Quick note cancelled")
			w.Close()
		}
	})

	// Focus the input field
	w.Canvas().Focus(input)

	w.Show()
}
