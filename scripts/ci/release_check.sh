#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
VERSION="${VERSION:-v0.1.0}"
export POSTGRESX_REQUIRE_INTEGRATION="${POSTGRESX_REQUIRE_INTEGRATION:-1}"

GOWORK=off make ci-extended
GOWORK=off make integration
GOWORK=off make evidence VERSION="$VERSION"
GOWORK=off make release-evidence-check
GOWORK=off make release-final-check

module="$(GOWORK=off "$GO" list -m)"
if [[ "$module" != "github.com/ZoneCNH/postgresx" ]]; then
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

bash ./scripts/ci/release_evidence_check.sh "$VERSION"

echo "release check passed for $VERSION"
