// internal/gui/mainwindow.go
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/logger"
)

type MainWindow struct {
	gui    *GUI
	window fyne.Window
}

func NewMainWindow(gui *GUI) *MainWindow {
	return &MainWindow{
		gui:    gui,
		window: gui.fyneApp.NewWindow("Godo"),
	}
}

func (w *MainWindow) Setup() error {
	w.window.Resize(fyne.NewSize(800, 600))
	w.window.CenterOnScreen()

	// Set up the content
	content := widget.NewLabel("Welcome to Godo")
	w.window.SetContent(content)

	// Load and set application icon
	if err := w.setupIcon(); err != nil {
		logger.Error("Failed to setup icon", "error", err)
	}

	// Setup system tray
	if err := w.setupSystemTray(); err != nil {
		logger.Error("Failed to setup system tray", "error", err)
	}

	// Setup shortcuts
	if err := w.setupShortcuts(); err != nil {
		logger.Error("Failed to setup shortcuts", "error", err)
	}

	// Show the quick note window after main window setup
	ShowQuickNote(w.gui.ctx, w.gui)

	// Show main window initially
	w.window.Show()
	return nil
}

func (w *MainWindow) setupIcon() error {
	iconBytes, err := assets.GetIcon()
	if err != nil {
		return err
	}

	icon := fyne.NewStaticResource("icon", iconBytes)
	w.gui.fyneApp.SetIcon(icon)
	return nil
}
