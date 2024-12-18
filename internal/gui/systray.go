package gui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/jonesrussell/godo/internal/logger"
)

func (w *MainWindow) setupSystemTray() error {
	desk, ok := w.gui.fyneApp.(desktop.App)
	if !ok {
		logger.Error("System tray not supported on this platform")
		return errors.New("system tray not supported on this platform")
	}

	menu := fyne.NewMenu("Godo",
		fyne.NewMenuItem("Open", w.handleOpen),
		fyne.NewMenuItem("Quick Note", w.handleQuickNote),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", w.handleQuit),
	)

	desk.SetSystemTrayMenu(menu)
	desk.SetSystemTrayIcon(w.gui.fyneApp.Icon())
	return nil
}

func (w *MainWindow) handleOpen() {
	logger.Debug("Opening main window")
	w.window.Show()
	w.window.RequestFocus()
	w.window.CenterOnScreen()
}

func (w *MainWindow) handleQuickNote() {
	logger.Debug("Opening quick note from tray")
	ShowQuickNote(w.gui.ctx, w.gui)
}

func (w *MainWindow) handleQuit() {
	logger.Info("Quitting application")
	w.gui.Cleanup()
}
