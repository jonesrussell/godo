//go:build linux || windows

package hotkey

import "testing"

type stubQuick struct {
	called bool
}

func (s *stubQuick) Show() { s.called = true }
func (s *stubQuick) Hide() {}

type stubMain struct {
	called bool
}

func (s *stubMain) Show() { s.called = true }
func (s *stubMain) Hide() {}

func TestDispatchHotkeyForEntry_QuickNote(t *testing.T) {
	t.Parallel()
	q := &stubQuick{}
	m := &HotkeyManager{
		hotkeys: []HotkeyEntry{{quickNote: q}},
	}
	m.dispatchHotkeyForEntry(0)
	if !q.called {
		t.Fatal("expected quick note Show")
	}
}

func TestDispatchHotkeyForEntry_MainWindow(t *testing.T) {
	t.Parallel()
	w := &stubMain{}
	m := &HotkeyManager{
		hotkeys: []HotkeyEntry{{mainWindow: w}},
	}
	m.dispatchHotkeyForEntry(0)
	if !w.called {
		t.Fatal("expected main window Show")
	}
}
