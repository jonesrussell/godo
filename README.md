# Godo

Desktop **Go** app: **Fyne** UI, **global hotkeys**, and a **REST API** for notes.
**Clean Architecture** layers; **Google Wire** composes dependencies.

**Local-first companion:** hotkeys, windows, and state work offline. **SQLite** is the canonical store.
The charter targets durable data and an **append-only audit trail**. Caches are projections only.
Harnesses stay behind interfaces; **fakes and fixtures** cover tests until integration is chartered.

Today Godo still delivers quick-note and note-management workflows. The paragraphs above match
[.kittify/charter/charter.md](.kittify/charter/charter.md).

## What it does today

1. **Quick note** — Global hotkey opens a minimal window; capture and dismiss with low friction.
2. **Main window** — List and manage notes (CRUD-style).
3. **REST API** — Same operations over HTTP for scripts and integrations.

Default hotkeys (override in `config.yaml`):

| Action       | Default        |
| ------------ | -------------- |
| Quick note   | Ctrl+Shift+1   |
| Main window  | Ctrl+Shift+2   |

## Spec-driven development

Governance and structured feature work use
**[Spec Kitty](https://github.com/spec-kitty/spec-kitty)**.

- **Charter (human source):** [Project charter](.kittify/charter/charter.md)
  — local-first intent, SQLite, harness boundary.
- **CLI:** `spec-kitty verify-setup`, `spec-kitty charter sync` after edits.
  See Spec Kitty docs for other commands.
- **Agents:** `.kittify/config.yaml` and `.agents/skills/` (Cursor, Claude Code, Codex).

## Go 1.26 audit (machine-readable)

Results from the Go 1.26 upgrade / audit pass are recorded for tooling and CI reference:

- [docs/audit/go126-audit-report.json](docs/audit/go126-audit-report.json) — structured audit output.
- [go126-gofix-captured.diff](docs/audit/go126-gofix-captured.diff) — optional `gofix` diff capture.

## Features

- **Platforms:** Windows (system tray), Linux (no tray), macOS planned.
- **Storage:** SQLite; optional API-backed storage via `config.yaml`.
- **Hotkeys:** OS-level registration where supported (WSL2 has known limitations without extra setup).
- **Logging:** Structured (Zap); tune with config and `LOG_LEVEL`.
- **Quality:** `task fmt`, `task lint`, `go test ./... -tags=wireinject`.
  Charter defines coverage on core layers.
- **Builds:** cross-platform and Docker via Task.
  [Taskfile.yml](Taskfile.yml), [Taskfile.build.yml](Taskfile.build.yml).

## API

Default base URL: **`http://localhost:8008`** (`http.port` in [config.yaml](config.yaml)).

| Method | Path                  | Description   |
| ------ | --------------------- | ------------- |
| GET    | `/health`             | Health check  |
| GET    | `/api/v1/notes`      | List notes    |
| POST   | `/api/v1/notes`      | Create note   |
| PUT    | `/api/v1/notes/{id}` | Update note   |
| DELETE | `/api/v1/notes/{id}` | Delete note   |

```bash
curl -s http://localhost:8008/health
curl -s http://localhost:8008/api/v1/notes
```

## Prerequisites

- **Go 1.26+** ([go.mod](go.mod))
- **SQLite** (system `sqlite3` optional; app may use pure Go SQLite per build tags / config)
- **[Task](https://taskfile.dev/installation/)**
- **CGO** — global hotkeys need `golang.design/x/hotkey`. **Do not** `go mod vendor` (breaks CGO). Use
  `go mod download` and `go mod tidy`.

**Windows:** MinGW-w64 GCC for CGO. **Linting:** GNU diffutils may be required on Windows.

## Quick start

```bash
git clone https://github.com/jonesrussell/godo.git
cd godo
task install-tools   # optional dev tools from go.mod tool directive
task deps
task wire            # regenerate Wire when providers change
task run
```

Common checks:

```bash
task fmt && task lint && task test
task check           # fast format + lint
task dev             # fmt, lint, test
task build           # current platform → dist/
```

More detail: [CLAUDE.md](CLAUDE.md).

## Architecture

```text
main.go → container.InitializeApp() → Wire → domain / application / infrastructure
```

- **`internal/domain/`** — Models, repository interfaces, services (no outward infrastructure imports).
- **`internal/application/`** — Orchestration; Wire [container](internal/application/container/).
- **`internal/infrastructure/`** — Fyne UI, HTTP API, storage, hotkeys, logging, platform helpers.

## Contributing

1. Fork and branch from the integration branch agreed for your change.
2. Run `task fmt`, `task lint`, and `task test` (with `-tags=wireinject` where the tree requires Wire).
3. Open a PR with a clear description; conventional commits are appreciated.

## License

[MIT](LICENSE)

## Acknowledgments

- [Fyne](https://fyne.io/) for the GUI toolkit
- Contributors and everyone dogfooding capture workflows
