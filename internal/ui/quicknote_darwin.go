//go:build darwin

package ui

import (
	"context"
	"fmt"
)

type DarwinQuickNote struct {
	inputChan chan string
}

func NewQuickNoteUI() (QuickNoteUI, error) {
	return &DarwinQuickNote{
		inputChan: make(chan string, 1),
	}, nil
}

func (d *DarwinQuickNote) Show(ctx context.Context) error {
	return fmt.Errorf("macOS implementation not yet available")
}

func (d *DarwinQuickNote) Hide() error {
	return nil
}

func (d *DarwinQuickNote) GetInput() <-chan string {
	return d.inputChan
}
