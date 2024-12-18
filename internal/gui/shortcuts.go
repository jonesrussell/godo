package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/jonesrussell/godo/internal/logger"
)

func (w *MainWindow) setupShortcuts() error {
	shortcut := &desktop.CustomShortcut{
		KeyName:  "g",
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt,
	}

	w.window.Canvas().AddShortcut(shortcut, func(shortcut fyne.Shortcut) {
		logger.Debug("Quick note hotkey triggered")
		w.handleQuickNote()
	})

	logger.Info("Registered quick note hotkey (Ctrl+Alt+G)")
	return nil
}
