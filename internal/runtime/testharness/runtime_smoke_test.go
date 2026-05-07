package testharness

import (
	"context"
	gruntime "runtime"
	"sync/atomic"
	"testing"
	"time"

	runtimelayer "github.com/jonesrussell/godo/internal/runtime"
)

type mockApp struct {
	started  chan struct{}
	stop     chan struct{}
	shutdown atomic.Int32
}

func newMockApp() *mockApp {
	return &mockApp{
		started: make(chan struct{}),
		stop:    make(chan struct{}),
	}
}

func (m *mockApp) Run() {
	close(m.started)
	<-m.stop
}

func (m *mockApp) Shutdown(ctx context.Context) error {
	m.shutdown.Add(1)
	select {
	case <-m.stop:
		return nil
	default:
		close(m.stop)
		return nil
	}
}

func (m *mockApp) ForceKillTimeout() time.Duration {
	return 100 * time.Millisecond
}

func TestRuntimeRun_SmokeStartShutdownNormalize(t *testing.T) {
	before := gruntime.NumGoroutine()

	app := newMockApp()
	var cleanupCount atomic.Int32

	rootCtx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() {
		done <- runtimelayer.Run(rootCtx, app, func() {
			cleanupCount.Add(1)
		}, nil)
	}()

	select {
	case <-app.started:
	case <-time.After(1 * time.Second):
		t.Fatal("app did not start")
	}

	cancel()

	select {
	case code := <-done:
		if code != runtimelayer.ExitOK {
			t.Fatalf("unexpected exit code %d", code)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runtime.Run did not return")
	}

	if app.shutdown.Load() != 1 {
		t.Fatalf("shutdown count = %d, want 1", app.shutdown.Load())
	}
	if cleanupCount.Load() != 1 {
		t.Fatalf("cleanup count = %d, want 1", cleanupCount.Load())
	}

	time.Sleep(100 * time.Millisecond)
	gruntime.GC()
	time.Sleep(100 * time.Millisecond)
	after := gruntime.NumGoroutine()
	if after > before+8 {
		t.Fatalf("possible goroutine leak: before=%d after=%d", before, after)
	}
}
