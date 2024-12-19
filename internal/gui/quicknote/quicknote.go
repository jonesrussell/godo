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
	input  *widget.Entry
	form   dialog.Dialog
}

// New creates a new QuickNote instance
func New(window fyne.Window, store storage.Store) *QuickNote {
	qn := &QuickNote{
		window: window,
		store:  store,
		input:  widget.NewMultiLineEntry(),
	}

	// Set up the input
	qn.input.SetPlaceHolder("Enter your note here...")

	// Create the form once
	qn.form = dialog.NewForm("Quick Note", "Save", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Note", qn.input),
		},
		func(save bool) {
			if save && qn.input.Text != "" {
				todo := model.NewTodo(qn.input.Text)
				if err := qn.store.Add(todo); err != nil {
					logger.Error("Failed to save todo", "error", err)
					dialog.ShowError(err, qn.window)
				} else {
					logger.Debug("Saved note as todo", "id", todo.ID, "content", todo.Content)
				}
			}
			qn.input.SetText("") // Clear the input
			qn.window.Hide()
		},
		qn.window)

	qn.form.Resize(fyne.NewSize(400, 200))
	return qn
}

// Show displays the quick note dialog
func (qn *QuickNote) Show() {
	qn.window.Show()
	qn.window.CenterOnScreen()
	qn.form.Show()
	qn.window.Canvas().Focus(qn.input)
}
