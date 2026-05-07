//go:build linux || windows
// +build linux windows

package hotkey

import (
	"testing"

	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

func TestHotkeyManager_listenChannels_nilHotkeys(t *testing.T) {
	t.Parallel()

	log := logger.NewTestLogger(t)
	m := &HotkeyManager{
		log: log,
		hotkeys: []HotkeyEntry{
			{},
			{},
		},
	}

	ch := m.listenChannels()
	if len(ch) != len(m.hotkeys) {
		t.Fatalf("len: got %d want %d", len(ch), len(m.hotkeys))
	}
	for i, c := range ch {
		if c != nil {
			t.Fatalf("expected nil channel at %d", i)
		}
	}
}
