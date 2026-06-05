#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

version="${1:-${VERSION:-v0.1.0}}"
manifest="release/manifest/${version}.json"
latest="release/manifest/latest.json"
required=(
  "docs/RELEASE_MANIFEST-${version}.md"
  "docs/EVIDENCE-20260601.md"
  "docs/RETROSPECTIVE-GOAL-20260601-001.md"
  "docs/VERSION_MATRIX.md"
  "$manifest"
  "$manifest.sha256"
  "$latest"
  "$latest.sha256"
  "docs/evidence/20260601/dependencies.txt"
  "docs/evidence/20260601/go-test.txt"
  "docs/evidence/20260601/go-test-race.txt"
  "docs/evidence/20260601/go-vet.txt"
  "docs/evidence/20260601/gofmt.txt"
  "docs/evidence/20260601/migration-up-down-up.txt"
  "docs/evidence/20260601/no-xgo-deps.txt"
  "docs/evidence/20260601/postgres-integration.txt"
  "docs/evidence/20260601/secret-scan.txt"
)

for path in "${required[@]}"; do
  if [[ ! -s "$path" ]]; then
    echo "missing required release evidence: $path" >&2
    exit 1
  fi
done

verify_sha256() {
  local file="$1"
  local checksum_file="$file.sha256"
  local expected actual
  expected="$(awk '{print $1}' "$checksum_file")"
  actual="$(sha256sum "$file" | awk '{print $1}')"
  if [[ -z "$expected" || "$expected" != "$actual" ]]; then
    echo "stale release manifest checksum: $checksum_file" >&2
    exit 1
  fi
}

verify_sha256 "$manifest"
verify_sha256 "$latest"

python3 - "$manifest" "$latest" "$version" <<'PY'
import json
import sys
manifest_path, latest_path, version = sys.argv[1:]
for path in (manifest_path, latest_path):
    with open(path, encoding="utf-8") as fh:
        data = json.load(fh)
    if data.get("module") != "github.com/ZoneCNH/postgresx":
        raise SystemExit(f"{path} has unexpected module")
    if data.get("version") != version:
        raise SystemExit(f"{path} has unexpected version")
    if data.get("layer") != "L2":
        raise SystemExit(f"{path} has unexpected layer")
    if data.get("role") != "postgresql_infrastructure_adapter":
        raise SystemExit(f"{path} has unexpected role")
    required_hashes = {"api", "config", "health", "metrics", "errors"}
    hashes = data.get("contract_hashes") or {}
    missing = sorted(k for k in required_hashes if not str(hashes.get(k, "")).startswith("sha256:"))
    if missing:
        raise SystemExit(f"{path} missing contract hashes: {', '.join(missing)}")
PY

if rg -n 'github[.]com/[b]ytechainx|github[.]com/ZoneCNH/postgresx/pkg/postgresx/(examples|contracts)|go get github[.]com/ZoneCNH/postgresx/pkg/postgresx' \
  README.md docs contracts release .github --glob '!docs/goal.md' --glob '!docs/evidence/20260601/*'; then
  echo "release evidence contains stale module/package references" >&2
  exit 1
fi

if rg -n --hidden \
  --glob '!.git/**' \
  --glob '!docs/goal.md' \
  --glob '!docs/evidence/20260601/*' \
  --glob '!scripts/ci/release_evidence_check.sh' \
  --glob '!scripts/ci/secret_scan.sh' \
  'postgres(?:ql)?://[^[:space:]'"'"'"`]+:[^[:space:]'"'"'"`]+@|PGPASSWORD=[^[:space:]'"'"'"`]+|password=[^[:space:]'"'"'"`]+' \
  README.md docs contracts release scripts .github .agent 2>/dev/null; then
  echo "release evidence contains secret-like PostgreSQL credential material" >&2
  exit 1
fi

echo "release evidence check passed for $version"
