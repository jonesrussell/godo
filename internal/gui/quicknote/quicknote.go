package quicknote

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
)

// QuickNote represents a quick note dialog
type QuickNote struct {
	window fyne.Window
	store  storage.Store
}

// New creates a new QuickNote instance
func New(window fyne.Window, store storage.Store) *QuickNote {
	return &QuickNote{
		window: window,
		store:  store,
	}
}

// Show displays the quick note dialog
func (qn *QuickNote) Show() {
	// Ensure the window is shown when needed
	qn.window.Show()

	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Enter your note here...")

	form := dialog.NewForm("Quick Note", "Save", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Note", entry),
	}, func(save bool) {
		if save && entry.Text != "" {
			todo := model.NewTodo(entry.Text)
			if err := qn.store.Add(todo); err != nil {
				logger.Error("Failed to save todo", "error", err)
				dialog.ShowError(err, qn.window)
			} else {
				logger.Debug("Saved note as todo", "id", todo.ID, "content", todo.Content)
			}
		}
		// Hide the window after dialog is closed
		qn.window.Hide()
	}, qn.window)

	form.Resize(fyne.NewSize(400, 200))
	form.Show()
}
