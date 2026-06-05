#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
VERSION="${VERSION:-v0.1.0}"
export POSTGRESX_REQUIRE_INTEGRATION="${POSTGRESX_REQUIRE_INTEGRATION:-1}"

GOWORK=off make vet
GOWORK=off make test-unit
GOWORK=off make test-contract
GOWORK=off make test-integration
GOWORK=off make boundary
GOWORK=off make contracts
GOWORK=off make secret-scan
GOWORK=off make foundationx-api
GOWORK=off make template-alignment
GOWORK=off make evidence VERSION="$VERSION"

module="$(GOWORK=off "$GO" list -m)"
if [[ "$module" != "github.com/ZoneCNH/postgresx" ]]; then
  echo "unexpected module path: $module" >&2
  exit 1
fi

if ! GOWORK=off "$GO" list ./pkg/postgresx >/dev/null; then
  echo "core package github.com/ZoneCNH/postgresx/pkg/postgresx is not listable" >&2
  exit 1
fi

if GOWORK=off "$GO" list -deps ./pkg/postgresx | rg -n 'github.com/([b]ytechainx|ZoneCNH)/x\.go'; then
  echo "postgresx must not depend on application module" >&2
  exit 1
fi

if GOWORK=off "$GO" list -deps ./pkg/postgresx | rg -n 'github.com/ZoneCNH/(xlib-standard|testkitx|xlibgate)'; then
  echo "postgresx runtime package must not depend on L2 test or gate tooling" >&2
  exit 1
fi

python3 ./scripts/ci/l2_evidence.py --check --version "$VERSION"

echo "L2-T2 release check passed for $VERSION"
