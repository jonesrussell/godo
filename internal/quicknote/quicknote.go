package quicknote

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI interface {
	Show(ctx context.Context) error
	GetInput() <-chan string
}

type fyneQuickNote struct {
	window fyne.Window
	input  chan string
}

func New() (UI, error) {
	input := make(chan string)
	return &fyneQuickNote{
		input: input,
	}, nil
}

func (f *fyneQuickNote) Show(ctx context.Context) error {
	entry := widget.NewEntry()
	entry.OnSubmitted = func(text string) {
		f.input <- text
		f.window.Close()
	}

	content := container.NewVBox(
		widget.NewLabel("Quick Note:"),
		entry,
	)

	f.window.SetContent(content)
	f.window.Resize(fyne.NewSize(300, 100))
	f.window.CenterOnScreen()
	f.window.Show()

	return nil
}

func (f *fyneQuickNote) GetInput() <-chan string {
	return f.input
}
