// Package utils provides utility functions for the application
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// OpenFile opens a file with the default system application
func OpenFile(filePath string) error {
	// Resolve relative paths to absolute paths
	if !filepath.IsAbs(filePath) {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}
		filePath = absPath
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// OpenLogFile opens the main log file
func OpenLogFile(logPath string) error {
	return OpenFile(logPath)
}

// OpenErrorLogFile opens the error log file
func OpenErrorLogFile(errorLogPath string) error {
	return OpenFile(errorLogPath)
}
