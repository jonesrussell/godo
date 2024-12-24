//go:build !docker
// +build !docker

package quicknote

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
)

// customEntry extends widget.Entry to handle custom key events
type customEntry struct {
	widget.Entry
	onCtrlEnter func()
	onEscape    func()
	log         logger.Logger
}

// newCustomEntry creates a new customEntry instance
func newCustomEntry(log logger.Logger) *customEntry {
	entry := &customEntry{
		log: log,
	}
	entry.ExtendBaseWidget(entry)
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapWord
	return entry
}

// KeyDown handles keyboard events
func (e *customEntry) KeyDown(key *fyne.KeyEvent) {
	if e.log != nil {
		e.log.Debug("KeyDown event", "key", key.Name)
	}

	if key.Name == fyne.KeyEscape {
		if e.onEscape != nil {
			e.onEscape()
			return
		}
	}
	e.Entry.KeyDown(key)
}

// TypedKey handles typed key events
func (e *customEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyEscape {
		// Let KeyDown handle it
		return
	}
	e.Entry.TypedKey(key)
}

// TypedShortcut handles keyboard shortcuts
func (e *customEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if cs, ok := shortcut.(*desktop.CustomShortcut); ok {
		if cs.KeyName == fyne.KeyReturn || cs.KeyName == fyne.KeyEnter {
			if cs.Modifier == desktop.ControlModifier {
				if e.onCtrlEnter != nil {
					e.onCtrlEnter()
					return
				}
			}
		}
	}
	e.Entry.TypedShortcut(shortcut)
}
