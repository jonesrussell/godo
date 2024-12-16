package hotkey

import (
	"context"
	"testing"
	"time"
)

func TestHotkeyRegistration(t *testing.T) {
	manager := NewHotkeyManager()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start manager in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- manager.Start(ctx)
	}()

	// Wait briefly to ensure registration
	time.Sleep(500 * time.Millisecond)

	// Cleanup
	cancel()
	if err := <-errChan; err != nil && err != context.Canceled {
		t.Errorf("unexpected error: %v", err)
	}
}
