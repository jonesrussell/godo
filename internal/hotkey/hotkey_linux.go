//go:build linux
// +build linux

package hotkey

import (
	"context"

	hook "github.com/robotn/gohook"
)

type linuxHotkeyManager struct {
	eventChan chan struct{}
	stop      chan struct{}
}

func init() {
	newPlatformHotkeyManager = func() (HotkeyManager, error) {
		return &linuxHotkeyManager{
			eventChan: make(chan struct{}),
			stop:      make(chan struct{}),
		}, nil
	}
}

func (h *linuxHotkeyManager) Start(ctx context.Context) error {
	go func() {
		evChan := hook.Start()
		defer hook.End()

		for {
			select {
			case <-h.stop:
				return
			case ev := <-evChan:
				if ev.Kind == hook.KeyHold {
					h.eventChan <- struct{}{}
				}
			}
		}
	}()
	return nil
}

func (h *linuxHotkeyManager) Stop() error {
	close(h.stop)
	return nil
}

func (h *linuxHotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}

func (h *linuxHotkeyManager) RegisterHotkey(binding HotkeyBinding) error {
	// gohook handles the registration automatically
	// You might want to store the binding for reference
	return nil
}
