#!/usr/bin/env bash
# Regenerate Wire injectors for the current host OS (linux or windows tags).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}/internal/application/container"
TAG=linux
case "$(go env GOOS)" in
windows) TAG=windows ;;
esac
go run github.com/google/wire/cmd/wire@v0.7.0 gen -tags "${TAG}"
cd "${ROOT}"
echo "--- wire_gen.go diff (none if unchanged) ---"
git diff -- internal/application/container/wire_gen.go || true
