# Tasks: Runtime Layer Modernization & Unified App Lifecycle

## Tasking Principles
- Keep each PR narrowly scoped and reviewable.
- Preserve behavior while moving lifecycle ownership.
- Avoid business logic, UI feature, and storage feature changes.
- Keep final `os.Exit` usage only in `main.go`.

## Work Packages

- [ ] WP01 - Baseline lifecycle mapping and test harness
  - [ ] Document current startup/shutdown sequence in `main.go`.
  - [ ] Capture current signal, panic, and timeout behavior expectations.
  - [ ] Add focused lifecycle regression test scaffolding (no behavior changes).
  - [ ] Define acceptance checks for post-refactor parity.

- [ ] WP02 - Introduce runtime package skeleton and lifecycle interface
  - [ ] Add runtime package and `Run(...)` API surface.
  - [ ] Define interface to replace concrete type assertions.
  - [ ] Wire app dependencies to interface contract without feature changes.
  - [ ] Keep runtime hooks for future audit/harness shutdown participation.

- [ ] WP03 - Move root context and signal handling into runtime
  - [ ] Replace manual signal channels with `signal.NotifyContext`.
  - [ ] Establish root cancellation propagation from runtime.
  - [ ] Remove duplicate signal handling from `main.go`.
  - [ ] Validate parity with lifecycle baseline tests.

- [ ] WP04 - Consolidate panic recovery and exit-code normalization
  - [ ] Move panic recovery into runtime control loop.
  - [ ] Remove duplicated panic wrappers from `main.go`.
  - [ ] Normalize runtime error and panic outcomes to exit codes.
  - [ ] Ensure `main.go` only performs final process exit.

- [ ] WP05 - Move force-kill timeout and coordinated shutdown orchestration
  - [ ] Move force-kill timeout policy into runtime.
  - [ ] Implement deterministic shutdown order:
    - UI
    - hotkeys
    - storage
    - audit flush hook (future)
    - harness hook (future)
  - [ ] Ensure shutdown is idempotent and tolerant of partial failures.
  - [ ] Verify timeout and cleanup behavior with tests.

- [ ] WP06 - Thin `main.go` and final lifecycle parity review
  - [ ] Reduce `main.go` to DI initialize -> runtime call -> final exit.
  - [ ] Remove remaining non-final `os.Exit` calls outside runtime boundary.
  - [ ] Remove obsolete lifecycle code paths from `main.go`.
  - [ ] Run full regression checks and prepare implementation PR summary.

## Suggested PR Sequence
1. PR-A: WP01
2. PR-B: WP02 + WP03
3. PR-C: WP04
4. PR-D: WP05
5. PR-E: WP06

## Definition of Done (Mission Execution)
- Runtime layer is the single lifecycle owner.
- `main.go` is a thin wrapper with one final exit point.
- No concrete runtime type assertions required.
- Signal handling, panic recovery, and timeout shutdown behavior remain stable.
- Unified shutdown sequencing is explicit and future-hook ready.
