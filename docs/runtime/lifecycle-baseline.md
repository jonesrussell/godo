# Runtime Lifecycle Baseline

This document tracks lifecycle responsibility migration from legacy `main.go`
logic into `internal/runtime`.

## Pre-runtime baseline (legacy `main.go`)
- Manual `signal.Notify` with `sigChan`.
- Manual goroutine-based force-kill timeout and direct `os.Exit(1)`.
- Two independent `defer recover()` blocks with direct process termination.
- Direct type assertion (`myapp.(*core.App)`) in lifecycle path.

## WP03: Root Context + Signal Handling Integration
- Runtime now exposes `NewRootContext()` and `WithSignals(...)` to create
  signal-bound contexts (`SIGINT`, `SIGTERM`).
- Signal ownership moves from channel management to context cancellation.

## WP04: Panic Recovery + Exit Normalization
- Runtime now exposes panic recovery helpers that convert panics to errors
  (`WithPanicRecovery`, `RecoverFromPanic`) without stdout prints.
- Runtime normalizes outcome errors to exit codes via `NormalizeExit(...)`
  (`ExitOK`, `ExitError`, `ExitPanic`, `ExitForced`).

## WP05: Coordinated Shutdown + Force-Kill Timeout
- Runtime now exposes `CoordinatedShutdown(...)` with deterministic sequencing:
  timeout derivation -> app shutdown -> cleanup callback.
- Force-kill timeout is policy-based (`ForceKillTimeout()` or 3s default) and
  returns `ErrForcedShutdown` instead of calling `os.Exit`.

## WP06: main.go Thinning + Runtime Integration
Before WP06, `main.go` owned panic recovery, signal wiring, force-kill timeout,
and direct lifecycle orchestration. After WP06:

- `main.go` initializes DI and delegates lifecycle control to runtime.
- Root context is created with `runtime.NewRootContext()`.
- App execution is wrapped by runtime panic handling via `runtime.Run(...)`.
- Shutdown is coordinated through `runtime.CoordinatedShutdown(...)`.
- Exit behavior is mapped by `runtime.NormalizeExit(...)`.
- `main.go` calls `os.Exit(code)` exactly once at the end.

Responsibility shift:
- `main.go`: bootstrap + final process exit only.
- `internal/runtime`: signal/context, panic/error conversion, shutdown
  orchestration, and exit-code policy.
