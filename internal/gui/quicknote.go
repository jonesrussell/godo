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

	qnCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	quickNote := gui.app.GetQuickNote()
	if quickNote == nil {
		logger.Error("Failed to get quick note instance")
		return
	}

	// Create a new Fyne window for quick note
	win := gui.fyneApp.NewWindow("Quick Note")
	win.Resize(fyne.NewSize(300, 100))
	win.CenterOnScreen()

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter your quick note...")

	submit := widget.NewButton("Add", func() {
		logger.Debug("Quick note submitted", "text", input.Text)
		if input.Text != "" {
			todoService := gui.app.GetTodoService()
			if _, err := todoService.CreateTodo(qnCtx, input.Text, ""); err != nil {
				logger.Error("Failed to create todo from quick note", "error", err)
			} else {
				logger.Debug("Successfully created todo from quick note")
			}
		}
		win.Close()
	})

	content := container.NewVBox(
		input,
		submit,
	)

	win.SetContent(content)
	win.Canvas().Focus(input)

	// Handle window close
	win.SetOnClosed(func() {
		logger.Debug("Quick note window closed")
		cancel()
	})

	logger.Debug("Showing quick note window")
	win.Show()
}
