#!/usr/bin/env bash
# Fail if wire_gen.go content does not match internal/application/container/wire_gen.sha256.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}/internal/application/container"
TAG=linux
case "$(go env GOOS)" in
windows) TAG=windows ;;
esac
go run github.com/google/wire/cmd/wire@v0.7.0 gen -tags "${TAG}"
cd "${ROOT}"
GEN="${ROOT}/internal/application/container/wire_gen.go"
SUM_FILE="${ROOT}/internal/application/container/wire_gen.sha256"
if ! command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "${GEN}" | awk '{print $1}')
else
  ACTUAL=$(sha256sum "${GEN}" | awk '{print $1}')
fi
EXPECTED=$(tr -d '[:space:]' <"${SUM_FILE}")
if [ "${ACTUAL}" != "${EXPECTED}" ]; then
  echo "wire_gen drift — run ./scripts/regen-wire.sh, refresh internal/application/container/wire_gen.sha256, and commit." >&2
  echo "expected=${EXPECTED}" >&2
  echo "actual=${ACTUAL}" >&2
  git diff -- internal/application/container/wire_gen.go >&2 || true
  exit 1
fi
