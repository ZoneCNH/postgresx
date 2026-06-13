#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

version="${1:-${VERSION:-v1.0.0}}"
manifest="release/manifest/${version}.json"
status=0

fail() {
  printf 'BLOCKER %s\n' "$1"
  status=1
}

short_sha() {
  local value="${1:-}"
  if [[ -z "$value" ]]; then
    printf 'missing'
    return 0
  fi
  git rev-parse --short "$value" 2>/dev/null || printf '%s' "${value:0:12}"
}

check_ancestor() {
  local ancestor="$1"
  local descendant="$2"
  local label="$3"

  if git merge-base --is-ancestor "$ancestor" "$descendant"; then
    printf '%s=true\n' "$label"
  else
    fail "${label}=false ancestor=$(short_sha "$ancestor") descendant=$(short_sha "$descendant")"
  fi
}

printf 'release blocker report for %s\n' "$version"
printf 'manifest: %s\n' "$manifest"

if [[ ! -s "$manifest" ]]; then
  fail "manifest_present=false path=$manifest"
  printf 'release blockers present for %s\n' "$version"
  exit "$status"
fi

mapfile -t manifest_fields < <(python3 - "$manifest" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as handle:
    data = json.load(handle)

print(data.get("commit", ""))
print(data.get("tree_sha", ""))
PY
)

manifest_commit="${manifest_fields[0]:-}"
manifest_tree="${manifest_fields[1]:-}"
head_commit="$(git rev-parse HEAD)"
tag_ref="refs/tags/${version}"

printf 'HEAD: %s\n' "$(short_sha "$head_commit")"
printf 'manifest_commit: %s\n' "${manifest_commit:-missing}"
printf 'manifest_tree: %s\n' "${manifest_tree:-missing}"

if [[ -z "$manifest_commit" ]]; then
  fail "manifest_commit_present=false"
  resolved_manifest_commit=""
else
  resolved_manifest_commit="$(git rev-parse --verify --quiet "${manifest_commit}^{commit}" 2>/dev/null || true)"
  if [[ -z "$resolved_manifest_commit" ]]; then
    fail "manifest_source_present=false commit=$manifest_commit"
  else
    printf 'manifest_source_present=true commit=%s\n' "$(short_sha "$resolved_manifest_commit")"
  fi
fi

if [[ -z "$manifest_tree" ]]; then
  fail "manifest_tree_present=false"
elif [[ -n "${resolved_manifest_commit:-}" ]]; then
  resolved_manifest_tree="$(git rev-parse "${resolved_manifest_commit}^{tree}")"
  if [[ "$resolved_manifest_tree" == "$manifest_tree" ]]; then
    printf 'manifest_tree_matches_source=true tree=%s\n' "$(short_sha "$manifest_tree")"
  else
    fail "manifest_tree_matches_source=false manifest_tree=$(short_sha "$manifest_tree") source_tree=$(short_sha "$resolved_manifest_tree")"
  fi
fi

tag_object="$(git rev-parse --verify --quiet "$tag_ref" 2>/dev/null || true)"
tag_commit="$(git rev-parse --verify --quiet "${tag_ref}^{commit}" 2>/dev/null || true)"

if [[ -z "$tag_commit" ]]; then
  fail "tag_present=false ref=$tag_ref"
else
  printf 'tag_object: %s\n' "$(short_sha "$tag_object")"
  printf 'tag_commit: %s\n' "$(short_sha "$tag_commit")"
fi

if [[ -n "${resolved_manifest_commit:-}" ]]; then
  check_ancestor "$resolved_manifest_commit" "$head_commit" "manifest_source_ancestor_of_head"
  if [[ -n "$tag_commit" ]]; then
    check_ancestor "$resolved_manifest_commit" "$tag_commit" "manifest_source_ancestor_of_tag"
    check_ancestor "$tag_commit" "$head_commit" "tag_commit_ancestor_of_head"
  fi
fi

if [[ "$status" -ne 0 ]]; then
  printf 'release blockers present for %s\n' "$version"
  printf 'safe closure options:\n'
  printf ' - publish a successor version from the current branch and generate fresh manifest/checksums\n'
  printf ' - approve a manifest-contract change for restored snapshot source metadata\n'
  printf ' - explicitly authorize a %s tag/release rewrite\n' "$version"
  exit "$status"
fi

printf 'no release blockers found for %s\n' "$version"
