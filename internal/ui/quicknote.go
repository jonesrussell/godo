package ui

import (
	"context"

	"fyne.io/fyne/v2"
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

// NewQuickNote creates a new QuickNote instance using an existing Fyne app
func NewQuickNote(service service.TodoServicer, app fyne.App) *QuickNote {
	return &QuickNote{
		service: service,
		app:     app,
		window:  app.NewWindow("Quick Note"),
	}
}

func (qn *QuickNote) Update() {
	logger.Debug("Updating quick note window...")
}

func (qn *QuickNote) Show() {
	logger.Debug("Opening quick note window...")

	qn.window.Resize(fyne.NewSize(300, 100))
	qn.window.SetFixedSize(true)
	qn.window.CenterOnScreen()

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter quick note...")
	input.OnSubmitted = func(text string) {
		if text != "" {
			todo, err := qn.service.CreateTodo(context.Background(), text, "")
			if err != nil {
				logger.Error("Failed to create todo: %v", err)
				return
			}
			logger.Debug("Created quick note: %s (ID: %d)", text, todo.ID)

			// Verify the todo was created
			fetchedTodo, err := qn.service.GetTodo(context.Background(), todo.ID)
			if err != nil {
				logger.Error("Failed to verify todo creation: %v", err)
			} else {
				logger.Debug("Verified todo creation: %+v", fetchedTodo)
			}
		}
		qn.window.Hide()
	}

	title := widget.NewLabel("🗒️ Quick Note")
	title.TextStyle = fyne.TextStyle{Bold: true}

	hint := widget.NewLabel("Press Enter to save • Esc to cancel")
	hint.TextStyle = fyne.TextStyle{Italic: true}

	content := container.NewVBox(
		title,
		input,
		hint,
	)

	paddedContent := container.NewPadded(content)
	qn.window.SetContent(paddedContent)

	qn.window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			logger.Debug("Quick note cancelled")
			qn.window.Hide()
		}
	})

	qn.window.SetCloseIntercept(func() {
		qn.window.Hide()
	})

	qn.window.Show()
	qn.window.RequestFocus()
	qn.window.Canvas().Focus(input)

	logger.Debug("Quick note window should now be visible")
}

// QuickNoteUI defines the interface for platform-specific quick note implementations
type QuickNoteUI interface {
	// Show displays the quick note input window
	Show(ctx context.Context) error
	// Hide closes the quick note window
	Hide() error
	// GetInput returns the channel that receives user input
	GetInput() <-chan string
}
