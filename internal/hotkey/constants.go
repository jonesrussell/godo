package hotkey

// Common modifiers for all platforms
const (
	MOD_ALT     = 1 << iota // Alt key (Option on macOS)
	MOD_CONTROL             // Ctrl key (Command on macOS)
	MOD_SHIFT               // Shift key
	MOD_WIN                 // Windows key (Command on macOS)
)

// Key codes that are common across platforms
const (
	KEY_A = "a"
	KEY_B = "b"
	KEY_C = "c"
	// ... add other common keys as needed
	KEY_G = "g"
	// ... continue with other keys
)

// Common modifier names for configuration
var ModifierNames = map[string]uint{
	"alt":     MOD_ALT,
	"ctrl":    MOD_CONTROL,
	"shift":   MOD_SHIFT,
	"win":     MOD_WIN,
	"command": MOD_WIN, // macOS alias
	"option":  MOD_ALT, // macOS alias
}
