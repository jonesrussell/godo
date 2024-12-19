package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
)

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		logger.Info("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		logger.Info("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		logger.Info("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		logger.Info("Lifecycle: Exited Foreground")
	})
}

func main() {
	// Initialize logger
	_, err := logger.Initialize()
	if err != nil {
		panic(err)
	}

	myApp := app.New()
	logLifecycle(myApp)
	mainWindow := myApp.NewWindow("Main Window")

	// Function to show quick note
	showQuickNote := func() {
		entry := widget.NewMultiLineEntry()
		entry.SetPlaceHolder("Enter your note here...")

		form := dialog.NewForm("Quick Note", "Save", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Note", entry),
		}, func(b bool) {
			if b {
				logger.Debug("Saving note", "content", entry.Text)
			}
		}, mainWindow)

		form.Resize(fyne.NewSize(400, 200))
		form.Show()
	}

	btn := widget.NewButton("Open Quick Note", func() {
		showQuickNote()
	})

	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to the simplified Fyne app!"),
		btn,
	))
	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.CenterOnScreen()

	mainWindow.ShowAndRun()
}
