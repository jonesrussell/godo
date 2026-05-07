package runtime

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

type fakeShutdownApp struct {
	timeout     time.Duration
	shutdownFn  func(ctx context.Context) error
	timeoutRead atomic.Int32
}

func (f *fakeShutdownApp) ForceKillTimeout() time.Duration {
	f.timeoutRead.Add(1)
	return f.timeout
}

func (f *fakeShutdownApp) Shutdown(ctx context.Context) error {
	if f.shutdownFn != nil {
		return f.shutdownFn(ctx)
	}
	return nil
}

func TestCoordinatedShutdown_CompletesInTime(t *testing.T) {
	t.Parallel()

	var cleanups atomic.Int32
	app := &fakeShutdownApp{
		timeout: 500 * time.Millisecond,
		shutdownFn: func(ctx context.Context) error {
			select {
			case <-time.After(5 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	err := CoordinatedShutdown(context.Background(), app, func() {
		cleanups.Add(1)
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if cleanups.Load() != 1 {
		t.Fatalf("cleanup count = %d, want 1", cleanups.Load())
	}
	if NormalizeExit(err) != ExitOK {
		t.Fatalf("NormalizeExit(nil) chain: got %d", NormalizeExit(err))
	}
}

func TestCoordinatedShutdown_ErrForcedShutdownOnTimeout(t *testing.T) {
	t.Parallel()

	var cleanups atomic.Int32
	app := &fakeShutdownApp{
		timeout: 30 * time.Millisecond,
		shutdownFn: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}

	err := CoordinatedShutdown(context.Background(), app, func() {
		cleanups.Add(1)
	})
	if !errors.Is(err, ErrForcedShutdown) {
		t.Fatalf("want ErrForcedShutdown, got %v", err)
	}
	if NormalizeExit(err) != ExitForced {
		t.Fatalf("NormalizeExit = %d, want ExitForced", NormalizeExit(err))
	}
	if cleanups.Load() != 1 {
		t.Fatalf("cleanup count = %d, want 1", cleanups.Load())
	}
}

func TestCoordinatedShutdown_ZeroTimeoutUsesDefault(t *testing.T) {
	prev := defaultForceKillTimeout
	defaultForceKillTimeout = 25 * time.Millisecond
	t.Cleanup(func() { defaultForceKillTimeout = prev })

	var cleanups atomic.Int32
	app := &fakeShutdownApp{
		timeout: 0,
		shutdownFn: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}

	err := CoordinatedShutdown(context.Background(), app, func() {
		cleanups.Add(1)
	})
	if !errors.Is(err, ErrForcedShutdown) {
		t.Fatalf("want ErrForcedShutdown, got %v", err)
	}
	if cleanups.Load() != 1 {
		t.Fatalf("cleanup count = %d, want 1", cleanups.Load())
	}
}

func TestCoordinatedShutdown_CleanupOncePerCall(t *testing.T) {
	t.Parallel()

	var cleanups atomic.Int32
	app := &fakeShutdownApp{timeout: time.Second}

	_ = CoordinatedShutdown(context.Background(), app, func() {
		cleanups.Add(1)
	})
	_ = CoordinatedShutdown(context.Background(), app, func() {
		cleanups.Add(1)
	})

	if cleanups.Load() != 2 {
		t.Fatalf("cleanup count = %d, want 2 (two CoordinatedShutdown calls)", cleanups.Load())
	}
}

func TestCoordinatedShutdown_NoGoroutineLeakSmoke(t *testing.T) {
	before := runtime.NumGoroutine()

	for range 30 {
		app := &fakeShutdownApp{
			timeout: 20 * time.Millisecond,
			shutdownFn: func(ctx context.Context) error {
				return nil
			},
		}
		_ = CoordinatedShutdown(context.Background(), app, func() {})
	}

	time.Sleep(150 * time.Millisecond)
	runtime.GC()
	time.Sleep(150 * time.Millisecond)

	after := runtime.NumGoroutine()
	if after > before+10 {
		t.Fatalf("possible goroutine leak: before=%d after=%d", before, after)
	}
}

func TestCoordinatedShutdown_NilApp(t *testing.T) {
	t.Parallel()

	err := CoordinatedShutdown(context.Background(), nil, func() {})
	if err == nil {
		t.Fatal("expected error")
	}
}
