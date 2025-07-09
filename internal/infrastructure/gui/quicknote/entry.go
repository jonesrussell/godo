package quicknote

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Entry is a custom entry widget for quick notes
type Entry struct {
	*widget.Entry
	onCtrlEnter func()
	onEscape    func()
}

// NewEntry creates a new quick note entry widget
func NewEntry() *Entry {
	entry := &Entry{Entry: widget.NewEntry()}
	entry.MultiLine = true
	entry.ExtendBaseWidget(entry)
	return entry
}

// SetOnCtrlEnter sets the callback for when Ctrl+Enter is pressed
func (e *Entry) SetOnCtrlEnter(callback func()) {
	e.onCtrlEnter = callback
}

// SetOnEscape sets the callback for when Escape is pressed
func (e *Entry) SetOnEscape(callback func()) {
	e.onEscape = callback
}

// TypedShortcut handles keyboard shortcuts
func (e *Entry) TypedShortcut(shortcut fyne.Shortcut) {
	if shortcut, ok := shortcut.(*desktop.CustomShortcut); ok {
		if shortcut.KeyName == fyne.KeyReturn && shortcut.Modifier == fyne.KeyModifierControl && e.onCtrlEnter != nil {
			e.onCtrlEnter()
			return
		}
		if shortcut.KeyName == fyne.KeyEscape && e.onEscape != nil {
			e.onEscape()
			return
		}
	}
	e.Entry.TypedShortcut(shortcut)
}
