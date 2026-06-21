#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=common.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh" "${1:?staging required}"
branch_gate

GO_IMAGE="${GO_TEST_IMAGE:-golang:1.25.11-alpine}"

if command -v go >/dev/null 2>&1 && command -v make >/dev/null 2>&1; then
  make test
elif command -v docker >/dev/null 2>&1; then
  echo "[+] unit test (docker) — go not on host"
  docker run --rm \
    -v "${PWD}/src:/src" \
    -w /src \
    "$GO_IMAGE" \
    go build -o /dev/null .
else
  echo "[!] go+make or docker required for unit tests" >&2
  exit 1
fi
