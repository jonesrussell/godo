//go:build windows

package ui

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/lxn/walk"
	"github.com/lxn/win"
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
	if w.window, err = walk.NewMainWindowWithName("Quick Note"); err != nil {
		return err
	}

	// Set window properties
	w.window.SetMinMaxSize(walk.Size{Width: 400, Height: 60}, walk.Size{Width: 400, Height: 60})
	w.window.SetLayout(walk.NewVBoxLayout())

	// Create input field
	if w.input, err = walk.NewLineEdit(w.window); err != nil {
		return err
	}

	// Handle key events
	w.input.KeyPress().Attach(func(key walk.Key) {
		if key == walk.KeyReturn {
			text := w.input.Text()
			w.inputChan <- text
			w.window.Close()
		} else if key == walk.KeyEscape {
			w.window.Close()
		}
	})

	// Get primary monitor work area
	var mi win.MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	monitor := win.MonitorFromWindow(win.HWND(w.window.Handle()), win.MONITOR_DEFAULTTOPRIMARY)
	if !win.GetMonitorInfo(monitor, &mi) {
		return fmt.Errorf("failed to get monitor info")
	}

	// Calculate center position
	x := (int(mi.RcWork.Right-mi.RcWork.Left) - 400) / 2
	y := (int(mi.RcWork.Bottom-mi.RcWork.Top) - 60) / 2

	w.window.SetBounds(walk.Rectangle{X: x, Y: y, Width: 400, Height: 60})
	w.window.SetVisible(true)
	w.input.SetFocus()

	return nil
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
