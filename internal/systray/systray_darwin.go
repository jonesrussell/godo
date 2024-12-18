//go:build darwin
// +build darwin

package systray

func init() {
	newManager = func() Manager {
		return &darwinSystray{}
	}
}

type darwinSystray struct {
	defaultManager
}
