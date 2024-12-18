//go:build windows
// +build windows

package systray

func init() {
	newManager = func() Manager {
		return &windowsSystray{}
	}
}

type windowsSystray struct {
	defaultManager
}
