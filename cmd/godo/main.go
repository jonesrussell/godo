package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/assets"
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

func setupSystemTray(myApp fyne.App, showQuickNote func()) {
	if desk, ok := myApp.(desktop.App); ok {
		// Load the icon
		iconBytes, err := assets.GetIcon()
		if err != nil {
			logger.Error("Failed to load icon", "error", err)
			return
		}

		// Create a static resource for the system tray icon
		icon := fyne.NewStaticResource("icon", iconBytes)

		quickNote := fyne.NewMenuItem("Quick Note", nil)
		quickNote.Icon = icon
		menu := fyne.NewMenu("Godo", quickNote)

		quickNote.Action = func() {
			logger.Debug("Opening quick note from tray")
			myApp.SendNotification(fyne.NewNotification("Quick Note", "Opening quick note..."))
			showQuickNote()
		}

		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(icon)
		logger.Info("System tray initialized")
	} else {
		logger.Warn("System tray not supported on this platform")
	}
}

func main() {
	// Initialize logger
	_, err := logger.Initialize()
	if err != nil {
		panic(err)
	}

	// Create app with unique ID
	myApp := app.NewWithID("io.github.jonesrussell.godo")
	logLifecycle(myApp)

	mainWindow := myApp.NewWindow("Godo")

	// Function to show quick note
	showQuickNote := func() {
		// Ensure the window is shown when needed
		mainWindow.Show()

		entry := widget.NewMultiLineEntry()
		entry.SetPlaceHolder("Enter your note here...")

		form := dialog.NewForm("Quick Note", "Save", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Note", entry),
		}, func(b bool) {
			if b {
				logger.Debug("Saving note", "content", entry.Text)
			}
			// Hide the window after dialog is closed
			mainWindow.Hide()
		}, mainWindow)

		form.Resize(fyne.NewSize(400, 200))
		form.Show()
	}

	// Set up system tray with the showQuickNote function
	setupSystemTray(myApp, showQuickNote)

	// Hide main window by default
	mainWindow.SetCloseIntercept(func() {
		logger.Debug("Window close intercepted, hiding instead")
		mainWindow.Hide()
	})

	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to the simplified Fyne app!"),
		widget.NewButton("Open Quick Note", showQuickNote),
	))
	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.CenterOnScreen()

	// Start hidden
	mainWindow.Hide()
	myApp.Run()
}
