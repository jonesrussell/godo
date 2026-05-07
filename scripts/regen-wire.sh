#!/usr/bin/env bash
# Regenerate Google Wire DI for the application container (current platform).
# Run from repo root after changing providers in internal/application/container/wire.go.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}/internal/application/container"
exec go run github.com/google/wire/cmd/wire@v0.7.0
