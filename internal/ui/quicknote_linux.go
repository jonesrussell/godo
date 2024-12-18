//go:build linux

package ui

import (
	"context"

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

	l.input, err = gtk.EntryNew()
	if err != nil {
		return err
	}

	l.input.Connect("activate", func() {
		text, _ := l.input.GetText()
		l.inputChan <- text
		l.window.Close()
	})

	l.window.Add(l.input)
	l.window.ShowAll()

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
