package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Initialize the Fyne application
	myApp := app.New()
	mainWindow := myApp.NewWindow("Main Window")

	// Function to show quick note
	showQuickNote := func() {
		// Create the entry widget for the note
		entry := widget.NewMultiLineEntry()
		entry.SetPlaceHolder("Enter your note here...")

		// Use a form dialog to display the quick note entry
		form := dialog.NewForm("Quick Note", "Save", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Note", entry),
		}, func(b bool) {
			if b {
				// Log the entered text (or handle save)
				println("Saving note:", entry.Text)
				// Implement save functionality here if needed
			}
		}, mainWindow)

		form.Resize(fyne.NewSize(400, 200))
		form.Show()
	}

	// Button to open the quick note window
	btn := widget.NewButton("Open Quick Note", func() {
		showQuickNote()
	})

	// Set up the main window content
	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to the simplified Fyne app!"),
		btn,
	))
	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.CenterOnScreen()

	// Show the main window
	mainWindow.ShowAndRun()
}
