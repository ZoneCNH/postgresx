#!/usr/bin/env bash
set -euo pipefail

if ! command -v rg >/dev/null 2>&1; then
  echo "rg is required for secret scan" >&2
  exit 1
fi

status=0
scan() {
  local pattern="$1"
  if rg -n --hidden \
    --glob '!.git/**' \
    --glob '!go.sum' \
    --glob '!docs/goal.md' \
    --glob '!scripts/ci/secret_scan.sh' \
    "$pattern" .; then
    status=1
  fi
}

scan "postgres(?:ql)?://[^[:space:]'\"\`]+:[^[:space:]'\"\`]+@"
scan "PGPASSWORD=[^[:space:]'\"\`]+"
scan "password=[^[:space:]'\"\`]+"

if [[ "$status" -ne 0 ]]; then
  echo "secret-like PostgreSQL credential pattern found" >&2
  exit "$status"
fi

echo "secret scan passed"
