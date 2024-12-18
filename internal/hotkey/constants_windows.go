//go:build windows
// +build windows

package hotkey

// Windows API constants
const (
	WIN_MOD_ALT                     = 0x0001
	WIN_MOD_CONTROL                 = 0x0002
	WIN_MOD_SHIFT                   = 0x0004
	WIN_MOD_WIN                     = 0x0008
	WIN_MOD_NOREPEAT                = 0x4000
	WM_HOTKEY                       = 0x0312
	PM_REMOVE                       = 0x0001
	ERROR_HOTKEY_ALREADY_REGISTERED = 0x0402
	ERROR_SUCCESS                   = 0
)

// Map platform-independent modifiers to Windows-specific modifiers
var WindowsModifierMap = map[uint]uint{
	MOD_ALT:     WIN_MOD_ALT,
	MOD_CONTROL: WIN_MOD_CONTROL,
	MOD_SHIFT:   WIN_MOD_SHIFT,
	MOD_WIN:     WIN_MOD_WIN,
}

// Helper function to convert platform-independent modifiers to Windows modifiers
func ToWindowsModifiers(mods uint) uint {
	var result uint
	for platformMod, winMod := range WindowsModifierMap {
		if mods&platformMod != 0 {
			result |= winMod
		}
	}
	return result
}
