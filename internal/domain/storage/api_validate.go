package storage

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidateAPIBaseURL ensures the API storage base URL is safe to use as an HTTP client base:
// scheme must be http or https, and user credentials must not be embedded in the URL.
func ValidateAPIBaseURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("API base URL is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid API base URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("API base URL scheme must be http or https, got %q", u.Scheme)
	}
	if u.User != nil {
		return fmt.Errorf("API base URL must not contain user credentials")
	}
	if u.Host == "" {
		return fmt.Errorf("API base URL must include a host")
	}
	return nil
}
