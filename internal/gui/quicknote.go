package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
)

func ShowQuickNote(ctx context.Context, gui *GUI) {
	logger.Debug("Opening quick note window")

	w := gui.fyneApp.NewWindow("Quick Note")
	w.SetFixedSize(true)

	// Create the entry widget for the note
	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Enter your note here...")

	// Use a form dialog to display the quick note entry
	form := dialog.NewForm("Quick Note", "Save", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Note", entry),
	}, func(b bool) {
		if b {
			logger.Debug("Saving note", "content", entry.Text)
			// TODO: Implement save functionality
		}
	}, w)

	// Resize and center the form dialog
	form.Resize(fyne.NewSize(400, 200))
	w.CenterOnScreen()
	form.Show()
}
