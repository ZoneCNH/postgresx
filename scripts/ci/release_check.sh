#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
export POSTGRESX_REQUIRE_INTEGRATION=1

make ci
make integration

module="$(GOWORK=off "$GO" list -m)"
if [[ "$module" != "github.com/bytechainx/postgresx" ]]; then
  echo "unexpected module path: $module" >&2
  exit 1
fi

if GOWORK=off "$GO" list -deps ./... | rg -n 'github.com/(bytechainx|ZoneCNH)/x\.go'; then
  echo "postgresx must not depend on x.go" >&2
  exit 1
fi

echo "release check passed"
