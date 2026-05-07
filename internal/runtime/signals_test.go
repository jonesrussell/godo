package runtime

import (
	"context"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"
)

func TestWithSignals_CancelsOnSignal(t *testing.T) {
	ctx, cancel := WithSignals(context.Background(), syscall.SIGUSR1)
	t.Cleanup(cancel)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("find process: %v", err)
	}

	if err := proc.Signal(syscall.SIGUSR1); err != nil {
		t.Fatalf("send signal: %v", err)
	}

	select {
	case <-ctx.Done():
		// Expected.
	case <-time.After(2 * time.Second):
		t.Fatal("context was not canceled after signal")
	}
}

func TestWithSignals_CancelFuncWorksWithoutSignal(t *testing.T) {
	ctx, cancel := WithSignals(context.Background(), syscall.SIGUSR1)
	cancel()

	select {
	case <-ctx.Done():
		// Expected.
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context was not canceled by cancel func")
	}
}

func TestWithSignals_NoGoroutineLeakSmoke(t *testing.T) {
	before := runtime.NumGoroutine()

	for i := 0; i < 25; i++ {
		_, cancel := WithSignals(context.Background(), syscall.SIGUSR1)
		cancel()
	}

	// Give the runtime time to run defers/cleanup from signal.NotifyContext.
	time.Sleep(100 * time.Millisecond)
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	after := runtime.NumGoroutine()
	if after > before+5 {
		t.Fatalf("possible goroutine leak: before=%d after=%d", before, after)
	}
}

