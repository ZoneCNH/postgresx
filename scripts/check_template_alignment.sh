#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

status=0
required_paths=(
  "go.mod"
  "pkg/postgresx"
  "contracts/config.schema.json"
  "contracts/error.schema.json"
  "contracts/health.schema.json"
  "contracts/metrics.md"
  "examples/basic/main.go"
  "testkit/postgres.go"
  "docs/EVIDENCE-20260601.md"
  "docs/RELEASE_MANIFEST-v0.1.0.md"
  "docs/RETROSPECTIVE-GOAL-20260601-001.md"
  ".agent"
  "release/manifest"
)

for path in "${required_paths[@]}"; do
  if [[ ! -e "$path" ]]; then
    echo "missing template-required path: $path" >&2
    status=1
  fi
done

module="$(GOWORK=off go list -m)"
if [[ "$module" != "github.com/ZoneCNH/postgresx" ]]; then
  echo "unexpected module path: $module" >&2
  status=1
fi

if find . -maxdepth 1 -name '*.go' -print | grep -q .; then
  echo "root package Go files are not allowed; core must live in pkg/postgresx" >&2
  status=1
fi

stale_module_path='github.com/bytechainx''/postgresx'
nested_non_core_path='github.com/ZoneCNH/postgresx/pkg/postgresx/''(examples|testkit|contracts|docs|internal)'
if rg -n "$stale_module_path|$nested_non_core_path" \
  --glob '!go.sum' \
  --glob '!docs/goal.md' \
  .; then
  echo "stale module/package path found" >&2
  status=1
fi

if [[ "$status" -ne 0 ]]; then
  exit "$status"
fi

echo "template alignment check passed"
