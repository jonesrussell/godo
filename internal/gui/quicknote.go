package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
)

func ShowQuickNote(ctx context.Context, gui *GUI) {
	logger.Debug("Opening quick note window")

	// Create window before other elements
	w := gui.fyneApp.NewWindow("Quick Note")
	w.SetFixedSize(true)

	// Create input field
	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Enter your quick note here...")
	// Make entry fill available space
	entry.Resize(fyne.NewSize(380, 150))
	entry.Move(fyne.NewPos(10, 10))

	submitBtn := widget.NewButton("Save", func() {
		if entry.Text != "" {
			logger.Debug("Saving note: " + entry.Text)
		}
		w.Close()
	})

	// Use a border container for better layout
	content := container.NewBorder(
		nil,       // top
		submitBtn, // bottom
		nil,       // left
		nil,       // right
		entry,     // center (fills remaining space)
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 200))
	w.CenterOnScreen()
	w.Canvas().Focus(entry)
	w.Show()
}
