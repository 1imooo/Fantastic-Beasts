#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=common.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh" "${1:?staging required}"
branch_gate

if command -v go >/dev/null 2>&1 && command -v make >/dev/null 2>&1; then
  make test
elif command -v docker >/dev/null 2>&1; then
  # Builder stage compile check — uses build context (COPY), not bind mounts.
  # Bind-mounting src/ fails on some self-hosted runners (job container vs host paths).
  echo "[+] unit test (docker build — builder stage)"
  docker build --target builder -q -f Dockerfile .
else
  echo "[!] go+make or docker required for unit tests" >&2
  exit 1
fi
