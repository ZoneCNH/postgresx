#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

version="${1:-${VERSION:-v1.0.0}}"
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

if ! cmp -s "$manifest" "$latest"; then
  echo "latest release manifest diverges from $manifest" >&2
  exit 1
fi

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

read -r manifest_commit manifest_tree < <(python3 - "$manifest" <<'PY'
import json
import sys
with open(sys.argv[1], encoding="utf-8") as fh:
    data = json.load(fh)
print(data.get("commit", ""), data.get("tree_sha", ""))
PY
)

if [[ ! "$manifest_commit" =~ ^[0-9a-f]{7,40}$ ]]; then
  echo "release manifest has invalid source commit: $manifest_commit" >&2
  exit 1
fi

if [[ ! "$manifest_tree" =~ ^[0-9a-f]{40}$ ]]; then
  echo "release manifest has invalid source tree: $manifest_tree" >&2
  exit 1
fi

resolved_manifest_commit="$(git rev-parse --verify --quiet "${manifest_commit}^{commit}" 2>/dev/null || true)"
if [[ -z "$resolved_manifest_commit" ]]; then
  echo "release manifest source commit is not present in git history: $manifest_commit" >&2
  exit 1
fi

resolved_manifest_tree="$(git rev-parse "${resolved_manifest_commit}^{tree}")"
if [[ "$resolved_manifest_tree" != "$manifest_tree" ]]; then
  echo "release manifest tree does not match source commit: $manifest_commit" >&2
  exit 1
fi

head_commit="$(git rev-parse HEAD)"
if ! git merge-base --is-ancestor "$resolved_manifest_commit" "$head_commit"; then
  echo "release manifest source commit is not an ancestor of HEAD: $manifest_commit" >&2
  exit 1
fi

if tag_commit="$(git rev-parse --verify --quiet "refs/tags/${version}^{commit}" 2>/dev/null)"; then
  if ! git merge-base --is-ancestor "$resolved_manifest_commit" "$tag_commit"; then
    echo "release manifest source commit is not an ancestor of tag ${version}: $manifest_commit" >&2
    exit 1
  fi
fi

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
