# Runtime Lifecycle Test Harness Scaffold (WP01)

## Purpose
This directory is a documentation scaffold for lifecycle-focused tests that will be added after the runtime layer is introduced.

WP01 intentionally adds **no Go implementation code**. This file captures planned harness structure and acceptance intent so later WPs can implement tests without re-discovering scope.

## Target Under Test (Future)
Planned runtime entrypoint:
- `runtime.Run(app)` (or final equivalent decided in runtime implementation WPs)

Planned lifecycle contract focus:
- startup orchestration
- root context and signal cancellation
- panic recovery normalization
- unified cleanup ordering
- force-kill timeout handling
- exit code mapping

## Proposed Test Layout (Future)
- `internal/runtime/`
  - `runtime.go` (future implementation)
  - `runtime_test.go` (table-driven lifecycle tests)
  - `fakes_test.go` (fake lifecycle participants)
  - `signals_test.go` (signal/cancel behavior tests)
  - `shutdown_order_test.go` (cleanup sequencing tests)
  - `panic_exit_test.go` (panic and exit-code behavior)

## Planned Fake Components
Fakes should model lifecycle participants without touching concrete GUI/storage internals:
- fake app lifecycle adapter
- fake UI shutdown participant
- fake hotkey participant
- fake storage participant
- fake audit flush participant (future hook)
- fake harness participant (future hook)

Each fake should track:
- invocation count
- invocation order
- cancellation visibility
- returned error (if configured)
- blocked/hanging mode for timeout tests

## Planned Test Scenarios
1. **Startup happy path**
   - runtime initializes and enters blocking run loop.
2. **Signal-driven shutdown**
   - root context cancellation triggers ordered cleanup.
3. **Cleanup order enforcement**
   - verifies UI -> hotkey -> storage -> audit -> harness order.
4. **Partial cleanup failure**
   - errors are collected/logged without skipping remaining participants.
5. **Panic normalization**
   - panic is recovered and mapped to deterministic exit code.
6. **Force-kill timeout**
   - hung participant triggers timeout/forced termination behavior.
7. **Idempotent shutdown**
   - duplicate shutdown signals do not double-run participant cleanup.

## Acceptance Criteria for Harness Implementation (Future WPs)
- Tests run without requiring real Fyne windows or global hotkey registration.
- No direct dependency on `*core.App` type assertions in runtime tests.
- Deterministic ordering assertions for lifecycle participants.
- Coverage for both normal and failure shutdown paths.

## WP01 Constraint Reminder
- This scaffold is documentation-only.
- Runtime code changes are deferred to subsequent work packages.
