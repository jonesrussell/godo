//go:build darwin

package ui

import (
	"context"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
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

	// Create input field
	d.input = cocoa.NSTextField_New()
	d.input.SetFrame(cocoa.NSRect{
		Origin: cocoa.NSPoint{X: 20, Y: 20},
		Size:   cocoa.NSSize{Width: 360, Height: 24},
	})

	// Handle key events
	d.input.SetTarget(core.Target(func(sender objc.Object) {
		text := d.input.StringValue()
		if text != "" {
			d.inputChan <- text
		}
		d.window.Close()
	}))
	d.input.SetAction(objc.Sel("sendAction:"))

	// Center the window
	screenRect := cocoa.NSScreen_Main().Frame()
	windowRect := cocoa.NSRect{
		Origin: cocoa.NSPoint{
			X: (screenRect.Size.Width - 400) / 2,
			Y: (screenRect.Size.Height - 60) / 2,
		},
		Size: cocoa.NSSize{Width: 400, Height: 60},
	}
	d.window.SetFrame(windowRect, false)

	d.window.SetContentView(d.input)
	d.window.MakeKeyAndOrderFront(nil)
	d.input.BecomeFirstResponder()

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
