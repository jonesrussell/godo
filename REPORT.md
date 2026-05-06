# Go 1.26 upgrade and audit report

This branch bumps `go.mod` to **Go 1.26**, applies **safe `go fix` / `gofmt` modernization** (`interface{}` → `any`, legacy `// +build` line removal where applicable, `min()` for slice bounds), improves **Wire invocation** via `go run …/wire@v0.7.0` in `Taskfile.build.yml`, adds **`native:darwin` / `wire:darwin`** tasks, and updates **CI** to Go 1.26 with Wire generation, `go vet`, `go fix` drift check, `go tool golangci-lint`, and tests.

Machine-readable findings: [`docs/audit/go126-audit-report.json`](docs/audit/go126-audit-report.json).  
Initial `go fix -diff` output (before `wire_gen` existed, including the early `InitializeApp` failure): [`docs/audit/go126-gofix-captured.diff`](docs/audit/go126-gofix-captured.diff).

---

## Executive summary — top five pre-release risks (paste as PR comment)

1. **CGO + Fyne + global hotkey portability** — Release binaries need CGO, platform linkers (X11/GL on Linux, MinGW for Windows cross-compile, Xcode on macOS). Audit host failed Linux link (`-lXxf86vm`) until dev packages are installed; Windows cross-build failed without `x86_64-w64-mingw32-gcc`. Treat CI image package lists as part of the release contract.

2. **Zero automated tests** — The tree contains no `*_test.go` files; `go test ./... -tags=wireinject` succeeds vacuously with **0% coverage**. Shipping without domain/storage/API tests is the largest functional risk.

3. **Static analysis debt (30 golangci-lint issues)** — Notable items: **duplicate** PATCH handlers in `internal/infrastructure/storage/api/store.go`, **gosec** on TLS `InsecureSkipVerify`, JWT field naming, SSRF-style client usage, **govet/shadow** in the same file, and **high cognitive complexity** in `HotkeyManager.Start`. CI runs the linter with `continue-on-error: true` until this backlog is cleared.

4. **`wire_gen.go` is gitignored** — Pattern `*_gen.go` excludes generated Wire output. Every clone and CI job **must** run Wire before `go build`. Drift or stale injectors are easy if someone commits partial state or skips the CI Wire step.

5. **HTTP API storage surface** — API-backed storage combines permissive TLS options, dynamic URLs, and repetitive HTTP code. Harden **base URL validation**, **TLS defaults**, and **redirect** handling before treating API mode as production-safe.

---

## Commands run (audit)

| Step | Result |
|------|--------|
| `go mod tidy` | OK after `go 1.26` |
| `go fix ./...` | Applied safe fixes; `go fix -diff` now empty |
| `gofmt -w .` | OK |
| `go vet ./...` | OK |
| `golangci-lint run ./...` | **30 issues** (see JSON `lint-issues`) |
| `go test ./... -tags=wireinject` | OK; no tests |
| `go test -coverprofile=cover.out ./... -tags=wireinject` | 0% total; `*.out` gitignored |
| `task build:native:linux` | Requires `wire` on PATH → addressed by `go run` in Taskfile; native link still needs system libs on the host |
| Cross Windows / Darwin | Captured failures in JSON `compile-errors` |

---

## Wire

- Regenerated with `go run github.com/google/wire/cmd/wire@v0.7.0 gen -tags linux`.
- Linux vs Windows `wire gen` outputs differ only in the `//go:generate` tag comment line when compared.

---

## Follow-ups (non-blocking for this draft PR)

- Clear or tune golangci-lint rules / fix findings, then flip CI lint step to hard-fail.
- Add `internal/...` unit and integration tests; re-enable meaningful coverage gates.
- Optionally stop ignoring `wire_gen.go` if the team wants generated code reviewed in PRs.
