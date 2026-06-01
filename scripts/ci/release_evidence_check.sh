#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

required_files=(
  "docs/EVIDENCE-20260601.md"
  "docs/RELEASE_MANIFEST-v0.1.0.md"
  "docs/RETROSPECTIVE-GOAL-20260601-001.md"
  "docs/VERSION_MATRIX.md"
  "release/manifest/v0.1.0.json"
  "docs/evidence/20260601/dependencies.txt"
  "docs/evidence/20260601/go-test.txt"
  "docs/evidence/20260601/go-test-race.txt"
  "docs/evidence/20260601/go-vet.txt"
  "docs/evidence/20260601/gofmt.txt"
  "docs/evidence/20260601/migration-up-down-up.txt"
  "docs/evidence/20260601/no-xgo-deps.txt"
  "docs/evidence/20260601/secret-scan.txt"
)

for file in "${required_files[@]}"; do
  if [[ ! -s "$file" ]]; then
    echo "missing release evidence artifact: $file" >&2
    exit 1
  fi
done

if rg -n 'github.com/bytechainx|github.com/ZoneCNH/postgresx/pkg/postgresx/(examples|contracts)|go get github.com/ZoneCNH/postgresx/pkg/postgresx' \
  README.md docs contracts release scripts .github --glob '!docs/goal.md' --glob '!docs/evidence/20260601/*'; then
  echo "release evidence contains stale module/package references" >&2
  exit 1
fi

echo "release evidence check passed"
