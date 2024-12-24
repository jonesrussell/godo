//go:build docker
// +build docker

package app

func (a *App) setupHotkey() error {
	// No-op for Docker environment
	return nil
}
