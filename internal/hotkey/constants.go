package hotkey

// Common modifiers for all platforms
const (
	MOD_ALT = 1 << iota
	MOD_CONTROL
	MOD_SHIFT
	MOD_WIN
)

// Platform-specific modifiers
const (
	MOD_COMMAND = MOD_WIN // macOS command key
	MOD_OPTION  = MOD_ALT // macOS option key
)
