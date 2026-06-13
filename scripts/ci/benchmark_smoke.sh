#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
out="$(mktemp)"
trap 'rm -f "$out"' EXIT

GOWORK=off "$GO" test ./pkg/postgresx -run '^$' -bench '^Benchmark' -benchtime=50x -benchmem -count=1 | tee "$out"

if ! rg -q '^Benchmark' "$out"; then
  echo "benchmark smoke did not execute any Benchmark* functions" >&2
  exit 1
fi

echo "benchmark smoke passed"
