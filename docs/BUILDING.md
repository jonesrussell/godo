# Building Godo

Cross-platform builds use **CGO** (Fyne GUI, global hotkeys, OpenGL). See `docs/audit/go126-audit-report.json` (`cgo-issues`, `compile-errors`) for CI notes from the Go 1.26 audit.

## Linux (native)

Install toolchain and X11/GL headers (Debian/Ubuntu example):

```bash
sudo apt-get install -y gcc pkg-config \
  libgl1-mesa-dev libx11-dev libxcursor-dev libxrandr-dev \
  libxinerama-dev libxi-dev libxxf86vm-dev libglx-dev
```

Then:

```bash
task wire
task build
```

## Windows (native or cross from Linux)

- **Native:** MSVC or MinGW-w64 with CGO enabled.
- **Cross-compile from Linux/WSL:** `mingw-w64` providing `x86_64-w64-mingw32-gcc` (see `Taskfile.build.yml` `cross-windows`).

```bash
task build:cross-windows
```

## macOS

Install **Xcode Command Line Tools** for the Apple `clang` toolchain, then build on a Mac host (`task build:native:darwin`). Darwin cross-compiles from Linux with CGO are not supported for this app.

## Wire

Regenerate injectors and inspect the diff:

```bash
task wire:regen
```

CI enforces that `internal/application/container/wire_gen.go` matches `internal/application/container/wire_gen.sha256` after generation (`scripts/check-wire-drift.sh`). When you change `wire.go` or providers, run `task wire:regen`, update the checksum file, and commit both in the same change.

## Tests (no CGO for SQLite tests)

```bash
go test ./... -tags=wireinject
```

SQLite-backed tests use `modernc.org/sqlite` (pure Go).
