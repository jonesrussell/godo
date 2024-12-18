//go:build linux
// +build linux

package systray

func init() {
	newManager = func() Manager {
		return &linuxSystray{}
	}
}

type linuxSystray struct {
	defaultManager
}
