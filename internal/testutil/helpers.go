// Package testutil provides testing utilities and mock implementations
package testutil

import (
	"time"
)

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// TimePtr returns a pointer to the given time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}
