#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
export POSTGRESX_REQUIRE_INTEGRATION=1

GOWORK=off make ci
make integration

module="$(GOWORK=off "$GO" list -m)"
if [[ "$module" != "github.com/ZoneCNH/postgresx" ]]; then
  echo "unexpected module path: $module" >&2
  exit 1
fi

if [[ "$(GOWORK=off "$GO" list ./pkg/postgresx)" != "github.com/ZoneCNH/postgresx/pkg/postgresx" ]]; then
  echo "unexpected core package path" >&2
  exit 1
fi

if GOWORK=off "$GO" list -deps ./... | rg -n 'github.com/(bytechainx|ZoneCNH)/x\.go'; then
  echo "postgresx must not depend on x.go" >&2
  exit 1
fi

GOWORK=off make release-evidence-check

echo "release check passed"
