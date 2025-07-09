// Package platform provides platform-specific detection and utilities
package platform

import (
	"os"
	"runtime"
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
	switch runtime.GOOS {
	case "windows":
		// On Windows, check for common headless indicators
		// These environment variables are typically set in CI/CD or server environments
		headlessIndicators := []string{
			"CI", "BUILD_ID", "JENKINS_URL", "GITHUB_ACTIONS", "TRAVIS", "CIRCLECI",
			"GITLAB_CI", "TEAMCITY_VERSION", "BAMBOO_BUILDKEY", "GO_SERVER_URL",
		}
		for _, indicator := range headlessIndicators {
			if os.Getenv(indicator) != "" {
				return true
			}
		}
		// Check if we're in a Windows service or console session
		if os.Getenv("SESSIONNAME") == "Console" && os.Getenv("USERNAME") == "" {
			return true
		}
		return false
	case "linux":
		// On Linux, check for DISPLAY environment variable
		display := os.Getenv("DISPLAY")
		return display == ""
	case "darwin":
		// On macOS, GUI is typically available
		return false
	default:
		// For other platforms, assume headless
		return true
	}
}

// SupportsGUI returns true if the environment supports GUI features
func SupportsGUI() bool {
	return !IsHeadless() && !IsWSL2()
}
