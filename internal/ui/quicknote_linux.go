//go:build linux

package ui

import (
	"context"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type LinuxQuickNote struct {
	window    *gtk.Window
	input     *gtk.Entry
	inputChan chan string
}

func NewQuickNoteUI() (QuickNoteUI, error) {
	gtk.Init(nil)
	return &LinuxQuickNote{
		inputChan: make(chan string, 1),
	}, nil
}

func (l *LinuxQuickNote) Show(ctx context.Context) error {
	var err error
	l.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return err
	}

	l.window.SetTitle("Quick Note")
	l.window.SetDefaultSize(400, 60)
	l.window.SetResizable(false)

	l.input, err = gtk.EntryNew()
	if err != nil {
		return err
	}

	// Handle Enter key
	l.input.Connect("activate", func() {
		text, _ := l.input.GetText()
		l.inputChan <- text
		l.window.Close()
	})

	// Handle Escape key
	l.window.Connect("key-press-event", func(_ interface{}, event *gdk.Event) bool {
		keyEvent := &gdk.EventKey{Event: event}
		if keyEvent.KeyVal() == gdk.KEY_Escape {
			l.window.Close()
			return true
		}
		return false
	})

	l.window.Add(l.input)

	// Center the window
	if err := l.centerWindow(); err != nil {
		return err
	}

	l.window.ShowAll()
	l.input.GrabFocus()

	return nil
}

func (l *LinuxQuickNote) Hide() error {
	if l.window != nil {
		l.window.Close()
	}
	return nil
}

func (l *LinuxQuickNote) GetInput() <-chan string {
	return l.inputChan
}

func (l *LinuxQuickNote) centerWindow() error {
	screen, err := l.window.GetScreen()
	if err != nil {
		return err
	}

	monitor := screen.GetMonitorAtWindow(l.window.GetWindow())
	geometry := screen.GetMonitorGeometry(monitor)

	x := (geometry.GetWidth() - 400) / 2
	y := (geometry.GetHeight() - 60) / 2

	l.window.Move(x, y)
	return nil
}
