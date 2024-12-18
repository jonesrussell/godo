package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
)

type QuickNoteEntry struct {
	widget.Entry
	window fyne.Window
}

func NewQuickNoteEntry(win fyne.Window) *QuickNoteEntry {
	entry := &QuickNoteEntry{window: win}
	entry.ExtendBaseWidget(entry)
	entry.SetPlaceHolder("Enter your quick note...")
	return entry
}

// FocusGained implements fyne.Focusable
func (e *QuickNoteEntry) FocusGained() {
	e.Entry.FocusGained()
}

func (e *QuickNoteEntry) KeyDown(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyEscape {
		if e.window != nil {
			e.window.Close()
		}
		return
	}
	e.Entry.KeyDown(key)
}

func ShowQuickNote(ctx context.Context, gui *GUI) {
	logger.Debug("Opening quick note window")

	win := gui.fyneApp.NewWindow("Quick Note")
	win.Resize(fyne.NewSize(400, 100))
	win.SetFixedSize(true)

	input := NewQuickNoteEntry(win)
	win.SetContent(input)

	// Center before showing
	win.CenterOnScreen()

	// Show and focus
	win.Show()
	input.FocusGained()
	win.Canvas().Focus(input)
}
