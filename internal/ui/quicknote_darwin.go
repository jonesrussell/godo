//go:build darwin

package ui

import (
	"context"

	"github.com/progrium/macdriver/cocoa"
)

type DarwinQuickNote struct {
	window    cocoa.NSWindow
	input     cocoa.NSTextField
	inputChan chan string
}

func NewQuickNoteUI() (QuickNoteUI, error) {
	return &DarwinQuickNote{
		inputChan: make(chan string, 1),
	}, nil
}

func (d *DarwinQuickNote) Show(ctx context.Context) error {
	cocoa.TerminateAfterWindowsClose = false

	d.window = cocoa.NSWindow_New()
	d.window.SetTitle("Quick Note")
	d.window.SetStyleMask(
		cocoa.NSWindowStyleMaskTitled |
			cocoa.NSWindowStyleMaskClosable |
			cocoa.NSWindowStyleMaskMiniaturizable,
	)
	d.window.SetFrame(cocoa.NSRect{
		Origin: cocoa.NSPoint{X: 200, Y: 200},
		Size:   cocoa.NSSize{Width: 400, Height: 60},
	}, false)

	d.input = cocoa.NSTextField_New()
	d.input.SetFrame(cocoa.NSRect{
		Origin: cocoa.NSPoint{X: 20, Y: 20},
		Size:   cocoa.NSSize{Width: 360, Height: 24},
	})

	d.window.SetContentView(d.input)
	d.window.MakeKeyAndOrderFront(nil)

	return nil
}

func (d *DarwinQuickNote) Hide() error {
	if d.window != nil {
		d.window.Close()
	}
	return nil
}

func (d *DarwinQuickNote) GetInput() <-chan string {
	return d.inputChan
}
