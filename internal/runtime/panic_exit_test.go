package runtime

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"
)

type captureLogger struct {
	msgs []string
}

func (c *captureLogger) Error(msg string, keysAndValues ...any) {
	c.msgs = append(c.msgs, msg)
}

func TestWithPanicRecovery_ReturnsErrorOnPanic(t *testing.T) {
	t.Parallel()

	log := &captureLogger{}
	err := WithPanicRecovery(log, func() error {
		panic("boom")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *RecoveredPanicError
	if !errors.As(err, &pe) {
		t.Fatalf("expected RecoveredPanicError, got %T", err)
	}
	if pe.Value.(string) != "boom" {
		t.Fatalf("unexpected panic value %#v", pe.Value)
	}
	if len(log.msgs) != 1 || log.msgs[0] != "panic recovered" {
		t.Fatalf("unexpected log: %+v", log.msgs)
	}
}

func TestWithPanicRecovery_NoPanicReturnsFnError(t *testing.T) {
	t.Parallel()

	want := errors.New("fn failed")
	err := WithPanicRecovery(nil, func() error {
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("expected fn error, got %v", err)
	}
}

func TestWithPanicRecovery_NoPanicNil(t *testing.T) {
	t.Parallel()

	err := WithPanicRecovery(nil, func() error { return nil })
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecoverFromPanic_NilLogger(t *testing.T) {
	t.Parallel()

	err := RecoverFromPanic(nil, "x")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNormalizeExit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, ExitOK},
		{"generic", errors.New("fail"), ExitError},
		{"recovered type", &RecoveredPanicError{Value: 1}, ExitPanic},
		{"wrapped sentinel", fmt.Errorf("wrap: %w", ErrRecoveredPanic), ExitPanic},
		{"forced", ErrForcedShutdown, ExitForced},
		{"wrapped forced", fmt.Errorf("wrap: %w", ErrForcedShutdown), ExitForced},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NormalizeExit(tt.err); got != tt.want {
				t.Fatalf("NormalizeExit() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestWithPanicRecovery_NoGoroutineLeakSmoke(t *testing.T) {
	before := runtime.NumGoroutine()

	for i := 0; i < 40; i++ {
		_ = WithPanicRecovery(nil, func() error { return nil })
		_ = WithPanicRecovery(nil, func() error {
			panic(i)
		})
	}

	time.Sleep(120 * time.Millisecond)
	runtime.GC()
	time.Sleep(120 * time.Millisecond)

	after := runtime.NumGoroutine()
	if after > before+8 {
		t.Fatalf("possible goroutine leak: before=%d after=%d", before, after)
	}
}
