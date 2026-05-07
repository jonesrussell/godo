# PR: CI hardening, tests, Wire enforcement, and targeted lint fixes

Branch: `fix/ci-tests-wire` (from `upgrade/go1.26-audit`).  
Prior audit: [`docs/audit/go126-audit-report.json`](docs/audit/go126-audit-report.json).

## Summary

| Area | Change |
|------|--------|
| **CI** | Go **1.26**; **Wire drift** on the Linux quality job via `scripts/check-wire-drift.sh` + `wire_gen.sha256` (Windows `wire_gen` text differs, so the main build matrix does not re-run the checksum gate); **`go fix -diff`** must be empty; **`go vet`** + **`go test -tags=wireinject`**; **golangci-lint** with `continue-on-error: true` plus an issue-summary echo; optional **`build-matrix`** job (linux/windows/darwin) with `continue-on-error: true` and uploadable artifacts. **Wire / `go run`:** the build job no longer exports `GOOS`/`GOARCH` for the whole step (that forced `go run wire` to target Windows on `ubuntu-latest` and broke `wire:windows`); **`build-matrix`** runs Wire under `env -u GOOS -u GOARCH`. **Headless tests:** quality job starts **Xvfb** and sets `DISPLAY` so `csturiale/hotkey` init does not panic during `go test`. |
| **Wire** | `scripts/regen-wire.sh` prints `git diff` for `wire_gen.go`; `task wire:regen` / `task wire:drift-check`. **Do not commit** `wire_gen.go` (still gitignored); update **`wire_gen.sha256`** when Wire inputs change (Linux graph). |
| **Tests** | `internal/domain/testfixtures` (temp SQLite via `modernc.org/sqlite`), repository SQLite CRUD tests, `internal/infrastructure/http` API integration tests with JWT, plus small tests for TLS defaults, PATCH helper, and API URL validation. |
| **API / TLS** | Default **verified TLS**; opt-in via **`storage.api.tls_insecure_skip_verify`** (legacy **`insecure_skip_verify`** still honored in `wire.go`). **`ValidateAPIBaseURL`**: `http`/`https` only, no embedded credentials. |
| **Lint / quality** | Consolidated PATCH JSON paths in **`api/store.go`**, fixed **shadowing**, **`http.NoBody`**, removed duplicate 404 branch, refactored **hotkey** `Start` into helpers with unit tests, fixed **sqlite adapter** shadow warnings. |
| **Docs** | [`docs/BUILDING.md`](docs/BUILDING.md) — system packages, MinGW, Xcode/CLT, and Wire checksum workflow. |

## Acceptance criteria

- [x] CI runs Wire generation and **fails on drift** vs `wire_gen.sha256` (Linux quality job).
- [x] `go test ./... -tags=wireinject` passes with new tests.
- [x] `go vet` and **`go fix -diff`** drift checks pass (no diffs).
- [x] Critical audit items addressed: duplicate PATCH client code consolidated, TLS insecure explicit opt-in + URL validation, shadowing reduced, hotkey `Start` split for lower complexity.

## Known limitations

- Cross-platform **CGO** builds still need host toolchains. **`build-matrix`** is best-effort and may fail (especially darwin from `ubuntu-latest`); artifacts capture logs.

## PR body (Slack-style bullets)

- **Wire drift gate (Linux):** `scripts/check-wire-drift.sh` + `wire_gen.sha256`; CI fails if injectors drift without updating the manifest.
- **Tests:** SQLite-backed repository CRUD + `httptest` API CRUD with JWT; pure Go SQLite driver (no CGO in tests).
- **TLS / API storage:** default secure TLS; opt-in `tls_insecure_skip_verify` only; base URL must be `http`/`https` with no embedded credentials.
- **Lint / safety:** consolidated PATCH helpers, fixed govet shadowing, hotkey `Start` refactored into smaller functions with unit tests.
- **Cross-builds:** optional `build-matrix` job is `continue-on-error` with artifacts for diagnostics.

## Follow-ups

- Clear remaining golangci issues and remove `continue-on-error` on lint.
- Consider a second checksum for Windows `wire_gen` if drift should be enforced on Windows builds too.
- Optionally track `wire_gen.go` in git instead of a checksum file if the team prefers diff-only drift detection.
