#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

VERSION="${VERSION:-v0.1.0}"
OUT_DIR="release/manifest"
OUT_FILE="${OUT_FILE:-$OUT_DIR/$VERSION.json}"
mkdir -p "$OUT_DIR"

module="$(GOWORK=off go list -m)"
commit="$(git rev-parse --short HEAD 2>/dev/null || printf 'unknown')"
status="$(git status --short 2>/dev/null | sed 's/"/\\"/g' | paste -sd ';' -)"
created_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

cat > "$OUT_FILE" <<JSON
{
  "version": "$VERSION",
  "module": "$module",
  "core_package": "$module/pkg/postgresx",
  "commit": "$commit",
  "created_at": "$created_at",
  "dirty": $([[ -n "$status" ]] && printf true || printf false),
  "evidence": [
    "GOWORK=off make ci-extended",
    "GOWORK=off make release-preflight VERSION=$VERSION",
    "GOWORK=off make release-evidence-check",
    "GOWORK=off make release-final-check"
  ]
}
JSON

echo "$OUT_FILE"
