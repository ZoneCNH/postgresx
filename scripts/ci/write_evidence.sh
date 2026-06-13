#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"
version="${1:-${VERSION:-v1.0.0}}"
out="docs/evidence/20260601"
mkdir -p "$out"
GOWORK=off go list -m > "$out/go-list-module.txt"
GOWORK=off go list ./... > "$out/go-list-packages.txt"
GOWORK=off go test ./... > "$out/go-test.txt"
bash ./scripts/check_boundary.sh > "$out/boundary.txt"
bash ./scripts/check_contracts.sh > "$out/contracts.txt"
bash ./scripts/ci/secret_scan.sh > "$out/secret-scan.txt"
bash ./scripts/ci/release_evidence_check.sh "$version" > "$out/release-evidence-check.txt"
echo "evidence refreshed in $out"
