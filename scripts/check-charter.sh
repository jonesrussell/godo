#!/usr/bin/env bash
# Fail if the Spec Kitty charter is missing or empty (CI / local guard).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHARTER="${ROOT}/.kittify/charter/charter.md"
if [[ ! -f "${CHARTER}" ]]; then
  echo "check-charter: missing ${CHARTER}" >&2
  exit 1
fi
if [[ ! -s "${CHARTER}" ]]; then
  echo "check-charter: empty ${CHARTER}" >&2
  exit 1
fi
echo "check-charter: OK (${CHARTER})"
