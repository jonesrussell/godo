// Package platform provides platform-specific detection and utilities
package platform

import (
	"os"
	"strings"
	"sync"
)

var (
	isWSL2     bool
	isWSL2Once sync.Once
)

// IsWSL2 returns true if the application is running in a WSL2 environment
// The result is cached after the first call since platform detection won't change during runtime
func IsWSL2() bool {
	isWSL2Once.Do(func() {
		isWSL2 = checkWSLEnvironment() || checkKernelVersion()
	})
	return isWSL2
}

// checkWSLEnvironment checks for WSL-specific environment variables
func checkWSLEnvironment() bool {
	envVars := []string{"WSL_DISTRO_NAME", "WSLENV", "WSL_INTEROP"}
	for _, env := range envVars {
		if strings.Contains(strings.ToLower(os.Getenv(env)), "wsl") {
			return true
		}
	}
	return false
}

// checkKernelVersion checks if the kernel version contains "microsoft" indicating WSL2
func checkKernelVersion() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		// Silently return false - this is expected in non-Linux environments
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

// IsHeadless returns true if the environment doesn't support GUI features
func IsHeadless() bool {
	display := os.Getenv("DISPLAY")
	return display == ""
}

// SupportsGUI returns true if the environment supports GUI features
func SupportsGUI() bool {
	return !IsHeadless() && !IsWSL2()
}
