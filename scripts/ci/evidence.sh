#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

OUT_DIR="docs/evidence/20260601"
mkdir -p "$OUT_DIR"

GO="${GO:-go}"

run_capture() {
  local output="$1"
  shift
  {
    printf '$ %q' "$1"
    shift || true
    for arg in "$@"; do printf ' %q' "$arg"; done
    printf '\n'
    "$@"
  } >"$output" 2>&1
}

# Keep evidence generation side-effect-light and deterministic. Expensive live
# integration evidence is collected by the integration gate and may skip without
# a DSN/Docker unless POSTGRESX_REQUIRE_INTEGRATION=1 is set.
GOWORK=off "$GO" list -m all >"$OUT_DIR/dependencies.txt"
GOWORK=off "$GO" vet ./... >"$OUT_DIR/go-vet.txt" 2>&1
GOWORK=off "$GO" test ./... >"$OUT_DIR/go-test.txt" 2>&1
GOWORK=off "$GO" test -race ./... >"$OUT_DIR/go-test-race.txt" 2>&1
bash ./scripts/ci/secret_scan.sh >"$OUT_DIR/secret-scan.txt" 2>&1
bash ./scripts/ci/migration_up_down_up.sh >"$OUT_DIR/migration-up-down-up.txt" 2>&1
{
  echo '$ GOWORK=off go list -deps ./... | rg github.com/(bytechainx|ZoneCNH)/x\.go'
  if GOWORK=off "$GO" list -deps ./... | rg 'github.com/(bytechainx|ZoneCNH)/x\.go'; then
    exit 1
  fi
  echo 'no x.go dependency found'
} >"$OUT_DIR/no-xgo-deps.txt" 2>&1
{
  echo '$ gofmt -l $(find . -name *.go -not -path ./.git/*)'
  files="$(gofmt -l $(find . -name '*.go' -not -path './.git/*'))"
  if [[ -n "$files" ]]; then
    echo "$files"
    exit 1
  fi
  echo 'gofmt clean'
} >"$OUT_DIR/gofmt.txt" 2>&1

echo "evidence refreshed under $OUT_DIR"
