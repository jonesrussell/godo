// Package config handles application configuration management
package config

import (
	"fmt"
	"strings"
)

// ValidateHotkeyBinding validates a hotkey binding configuration
func ValidateHotkeyBinding(binding HotkeyBinding) error {
	if binding.Key == "" {
		return fmt.Errorf("hotkey key is required")
	}

	// Validate modifiers
	for _, mod := range binding.Modifiers {
		if !isValidModifier(mod) {
			return fmt.Errorf("invalid modifier: %s", mod)
		}
	}

	// Validate key
	if !isValidKey(binding.Key) {
		return fmt.Errorf("invalid key: %s", binding.Key)
	}

	return nil
}

// isValidModifier checks if a modifier is valid
func isValidModifier(mod string) bool {
	validModifiers := map[string]bool{
		"Ctrl":  true,
		"Shift": true,
		"Alt":   true,
		"Win":   true,
	}
	return validModifiers[mod]
}

// isValidKey checks if a key is valid
func isValidKey(key string) bool {
	// Basic validation - can be expanded as needed
	if key == "" {
		return false
	}

	// Single character keys (A-Z, 0-9)
	if len(key) == 1 {
		return (key >= "A" && key <= "Z") || (key >= "0" && key <= "9")
	}

	// Function keys
	if strings.HasPrefix(key, "F") && len(key) <= 3 {
		return true
	}

	// Special keys
	specialKeys := map[string]bool{
		"Space":  true,
		"Return": true,
		"Escape": true,
		"Delete": true,
		"Tab":    true,
		"Left":   true,
		"Right":  true,
		"Up":     true,
		"Down":   true,
	}

	return specialKeys[key]
}

// ValidateHTTPConfig validates HTTP server configuration
func ValidateHTTPConfig(http HTTPConfig) error {
	if http.Port < 1 || http.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", http.Port)
	}

	if http.ReadTimeout < 0 {
		return fmt.Errorf("invalid read timeout: %d (must be non-negative)", http.ReadTimeout)
	}

	if http.WriteTimeout < 0 {
		return fmt.Errorf("invalid write timeout: %d (must be non-negative)", http.WriteTimeout)
	}

	if http.ReadHeaderTimeout < 0 {
		return fmt.Errorf("invalid read header timeout: %d (must be non-negative)", http.ReadHeaderTimeout)
	}

	if http.IdleTimeout < 0 {
		return fmt.Errorf("invalid idle timeout: %d (must be non-negative)", http.IdleTimeout)
	}

	return nil
}

// ValidateUIConfig validates UI configuration
func ValidateUIConfig(ui UIConfig) error {
	if ui.MainWindow.Width < 100 {
		return fmt.Errorf("main window width too small: %d (minimum 100)", ui.MainWindow.Width)
	}

	if ui.MainWindow.Height < 100 {
		return fmt.Errorf("main window height too small: %d (minimum 100)", ui.MainWindow.Height)
	}

	if ui.QuickNote.Width < 50 {
		return fmt.Errorf("quick note width too small: %d (minimum 50)", ui.QuickNote.Width)
	}

	if ui.QuickNote.Height < 50 {
		return fmt.Errorf("quick note height too small: %d (minimum 50)", ui.QuickNote.Height)
	}

	return nil
}
