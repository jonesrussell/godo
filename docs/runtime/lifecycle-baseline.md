# Runtime Lifecycle Baseline

This document tracks lifecycle responsibilities while the runtime layer is being
introduced incrementally.

## Existing Signal Pattern (Current `main.go`)
- `main.go` creates a `sigChan` using `signal.Notify` for `SIGINT`/`SIGTERM`.
- A goroutine waits on the signal, runs cleanup, asserts `myapp.(*core.App)` to call
  `Quit()`, then starts another goroutine that sleeps `ForceKillTimeout()` (or 3s)
  and calls `os.Exit(1)` (“force kill”).

## WP03: Root Context + Signal Handling Integration
WP03 (when integrated) replaces the bespoke channel with `signal.NotifyContext`
helpers (`NewRootContext`, `WithSignals`) so cancellation flows through
`context.Context`. Wiring into `main.go` is a follow-up.

## WP04: Panic Recovery + Exit Normalization
`main.go` currently uses two `defer` + `recover()` blocks that print and exit.
WP04 centralizes panic-to-error conversion (`WithPanicRecovery`, `RecoverFromPanic`)
and stable exit codes (`NormalizeExit`, mapping `ErrForcedShutdown` → `ExitForced`).

## WP05: Coordinated Shutdown + Force-Kill Timeout
The force-kill path in `main.go` couples a sleeping goroutine to `os.Exit(1)`, which
bypasses normal error handling and repeats cleanup concerns.

WP05 adds `internal/runtime.CoordinatedShutdown` and the `ShutdownApp` interface:

- Shutdown uses a **`context.WithTimeout`** whose deadline comes from
  `app.ForceKillTimeout()`, defaulting to **3 seconds** when that duration is zero
  (`DefaultForceKillTimeout`).
- If shutdown hits **deadline exceeded**, runtime returns **`ErrForcedShutdown`** so
  **`NormalizeExit` → `ExitForced`** (`ExitForced`), with **no `os.Exit` and no stdout
  output** inside `internal/runtime`.
- A **`sync.Once` guard** ensures the optional resource **`cleanup`** callback runs
  **exactly once** per coordinated shutdown invocation.
- **Ordering** today is deterministic at the orchestration layer: derive timeout →
  `Shutdown(ctx)` → map errors; future work can extend `Shutdown` implementations
  to sequence UI, hotkeys, storage, audit, and harness without growing this skeleton.

Replacing `main.go`’s goroutine + `os.Exit` with `CoordinatedShutdown` remains a later
integration step.
