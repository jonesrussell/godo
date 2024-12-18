//go:build windows

package ui

import (
	"context"

	"github.com/lxn/walk"
)

type WindowsQuickNote struct {
	window    *walk.MainWindow
	input     *walk.LineEdit
	inputChan chan string
}

func NewQuickNoteUI() (QuickNoteUI, error) {
	return &WindowsQuickNote{
		inputChan: make(chan string, 1),
	}, nil
}

func (w *WindowsQuickNote) Show(ctx context.Context) error {
	var err error
	w.window, err = walk.NewMainWindow()
	if err != nil {
		return err
	}

	w.window.SetTitle("Quick Note")
	w.window.SetSize(walk.Size{Width: 400, Height: 60})
	w.window.SetLayout(walk.NewVBoxLayout())

	w.input, err = walk.NewLineEdit(w.window)
	if err != nil {
		return err
	}

	w.input.KeyPress().Attach(func(key walk.Key) {
		if key == walk.KeyReturn {
			text := w.input.Text()
			w.inputChan <- text
			w.window.Close()
		}
	})

	return w.window.Show()
}

func (w *WindowsQuickNote) Hide() error {
	if w.window != nil {
		w.window.Close()
	}
	return nil
}

func (w *WindowsQuickNote) GetInput() <-chan string {
	return w.inputChan
}
