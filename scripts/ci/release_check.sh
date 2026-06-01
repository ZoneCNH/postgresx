#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
version="${1:-${VERSION:-v0.1.0}}"

GOWORK=off make ci
GOWORK=off make integration

module="$(GOWORK=off "$GO" list -m)"
if [[ "$module" != "github.com/ZoneCNH/postgresx/pkg/postgresx" ]]; then
  echo "unexpected module path: $module" >&2
  exit 1
fi

if ! GOWORK=off "$GO" list ./pkg/postgresx >/dev/null; then
  echo "core package github.com/ZoneCNH/postgresx/pkg/postgresx is not listable" >&2
  exit 1
fi

if GOWORK=off "$GO" list -deps ./... | rg -n 'github.com/([b]ytechainx|ZoneCNH)/x\.go'; then
  echo "postgresx must not depend on application module" >&2
  exit 1
fi

bash ./scripts/ci/release_evidence_check.sh "$version"

echo "release check passed for $version"
