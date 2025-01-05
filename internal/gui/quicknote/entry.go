package quicknote

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Entry is a custom entry widget for quick notes
type Entry struct {
	widget.Entry
	onCtrlEnter func()
}

// NewEntry creates a new quick note entry widget
func NewEntry() *Entry {
	entry := &Entry{}
	entry.MultiLine = true
	entry.ExtendBaseWidget(entry)
	return entry
}

// SetOnCtrlEnter sets the callback for when Ctrl+Enter is pressed
func (e *Entry) SetOnCtrlEnter(callback func()) {
	e.onCtrlEnter = callback
}

// TypedShortcut handles keyboard shortcuts
func (e *Entry) TypedShortcut(shortcut fyne.Shortcut) {
	if shortcut, ok := shortcut.(*desktop.CustomShortcut); ok {
		if shortcut.KeyName == fyne.KeyReturn && shortcut.Modifier == fyne.KeyModifierControl && e.onCtrlEnter != nil {
			e.onCtrlEnter()
			return
		}
	}
	e.Entry.TypedShortcut(shortcut)
}
