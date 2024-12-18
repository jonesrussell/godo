//go:build windows

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
