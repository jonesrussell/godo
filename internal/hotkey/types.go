package hotkey

import "syscall"

// MSG represents a Windows message structure
type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// HotkeyConfig defines the configuration for a hotkey
type HotkeyConfig struct {
	WindowHandle syscall.Handle
	ID           int
	Modifiers    uint
	Key          rune
}

// DefaultConfig provides a different hotkey combination
var DefaultConfig = HotkeyConfig{
	WindowHandle: 0,
	ID:           1,
	Modifiers:    MOD_CONTROL | MOD_ALT,
	Key:          'G',
}
