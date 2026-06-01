#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v rg >/dev/null 2>&1; then
  echo "rg is required for contract checks" >&2
  exit 1
fi

required_files=(
  "contracts/config.schema.json"
  "contracts/error.schema.json"
  "contracts/health.schema.json"
  "contracts/metrics.md"
  "contracts/public_api.md"
  "docs/api.md"
)

for file in "${required_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "missing contract file: $file" >&2
    exit 1
  fi
done

for schema in contracts/*.schema.json; do
  rg -q '"\$schema"[[:space:]]*:' "$schema" || { echo "$schema missing \$schema" >&2; exit 1; }
  rg -q '"\$id"[[:space:]]*:[[:space:]]*"https://github.com/ZoneCNH/postgresx/contracts/' "$schema" || { echo "$schema has unexpected \$id" >&2; exit 1; }
  rg -q '"title"[[:space:]]*:' "$schema" || { echo "$schema missing title" >&2; exit 1; }
done

GO="${GO:-go}"
GOWORK=off "$GO" test ./contracts
