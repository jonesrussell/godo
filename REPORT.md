# PR: CI hardening, tests, Wire enforcement, and targeted lint fixes

Branch: `fix/ci-tests-wire` (from `upgrade/go1.26-audit`).  
Audit reference: [`docs/audit/go126-audit-report.json`](docs/audit/go126-audit-report.json).

## Summary

| Area | Change |
|------|--------|
| **CI** | Go **1.26**; **Wire drift** gate via `scripts/check-wire-drift.sh` + committed `wire_gen.sha256`; **`go fix -diff`** must be empty; **`go vet`** + **`go test -tags=wireinject`**; **golangci-lint** runs with `continue-on-error: true` and prints the issue summary line; **`build-matrix`** job attempts linux/windows/darwin cross-compile with `continue-on-error: true` and uploads `/tmp/godo-*` artifacts. |
| **Wire** | `scripts/regen-wire.sh` prints `git diff` for `wire_gen.go`; `task wire:regen` / `task wire:drift-check`. **Do not commit** `wire_gen.go` (still gitignored); update **`wire_gen.sha256`** whenever Wire inputs change. |
| **Tests** | `internal/domain/testfixtures` (temp SQLite via `modernc.org/sqlite`), **`note_repository_sqlite_test`**, **`internal/infrastructure/http/api_integration_test`** (JWT + CRUD), TLS/patch/API URL unit tests. |
| **API / TLS** | Default remains **verified TLS**; opt-in only via **`storage.api.tls_insecure_skip_verify`** (deprecated alias **`insecure_skip_verify`** still merged in `wire.go`). **`ValidateAPIBaseURL`**: `http`/`https` only, no embedded credentials. |
| **Lint / quality** | Deduplicated PATCH JSON paths in **`api/store.go`**, fixed **shadowing** (decode / sqlite adapter), removed duplicate **404** branch, **`http.NoBody`**, extracted **hotkey** event loop helpers + dispatch unit tests. |
| **Docs** | [`docs/BUILDING.md`](docs/BUILDING.md) — system packages and Wire checksum workflow. |

## Acceptance criteria (PR description)

- [x] CI runs Wire generation and **fails on drift** vs `wire_gen.sha256`.
- [x] `go test ./... -tags=wireinject` passes with new tests.
- [x] `go vet` and **`go fix -diff`** drift checks pass (no diffs).
- [x] Critical audit items addressed: duplicate PATCH client code consolidated, TLS insecure opt-in + URL validation, shadowing reduced, hotkey `Start` split for lower complexity.

## Known limitations

- Cross-platform **CGO** builds still need host toolchains (X11/GL on Linux, MinGW for Windows, Xcode on macOS). **`build-matrix`** is best-effort and may fail for darwin/Windows cross; logs and artifacts are retained.

## Follow-ups

- Clear remaining golangci issues and remove `continue-on-error` on lint.
- Optionally track `wire_gen.go` in git instead of a checksum file if the team prefers `git diff`-only drift detection.

## PR body (short summary — paste into GitHub)

- **Wire drift gate:** `scripts/check-wire-drift.sh` + `wire_gen.sha256`; CI fails if injectors drift without updating the manifest.
- **Tests:** SQLite-backed repository CRUD + `httptest` API CRUD with JWT; pure Go SQLite driver (no CGO in tests).
- **TLS / API storage:** default secure TLS; opt-in `tls_insecure_skip_verify` only; base URL must be `http`/`https` with no embedded credentials.
- **Lint / safety:** consolidated PATCH helpers, fixed govet shadowing, hotkey `Start` refactored into smaller functions with unit tests.
- **Cross-builds:** optional `build-matrix` job (linux/windows/darwin) is `continue-on-error` with artifacts for logs.
