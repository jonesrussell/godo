//go:build windows
// +build windows

package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
)

func init() {
	newPlatformQuickNoteUI = func() (QuickNoteUI, error) {
		return NewFyneQuickNote(), nil
	}
}

type FyneQuickNote struct {
	app    fyne.App
	window fyne.Window
	input  chan string
}

func NewFyneQuickNote() *FyneQuickNote {
	return &FyneQuickNote{
		app:   app.New(),
		input: make(chan string),
	}
}

func (f *FyneQuickNote) Show(ctx context.Context) error {
	logger.Debug("Showing Fyne quick note window")

	// Create window
	f.window = f.app.NewWindow("Quick Note")
	f.window.Resize(fyne.NewSize(300, 50))
	f.window.CenterOnScreen()

	// Create input field
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Type your note and press Enter...")
	entry.OnSubmitted = func(text string) {
		f.input <- text
		f.window.Close()
	}

	// Create container
	content := container.NewVBox(entry)
	f.window.SetContent(content)

	// Focus the input
	f.window.Canvas().Focus(entry)

	// Handle escape key
	f.window.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyEscape {
			f.window.Close()
		}
	})

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		f.window.Close()
	}()

	// Show window
	f.window.Show()

	return nil
}

func (f *FyneQuickNote) GetInput() <-chan string {
	return f.input
}
