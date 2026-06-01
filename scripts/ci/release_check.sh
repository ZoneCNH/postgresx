#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
version="${1:-${VERSION:-v0.1.0}}"

GOWORK=off make ci
make integration

module="$(GOWORK=off "$GO" list -m)"
if [[ "$module" != "github.com/ZoneCNH/postgresx" ]]; then
  echo "unexpected module path: $module" >&2
  exit 1
fi

if ! GOWORK=off "$GO" list ./pkg/postgresx >/dev/null; then
  echo "core package github.com/ZoneCNH/postgresx/pkg/postgresx is not listable" >&2
  exit 1
fi

legacy_org='byte''chainx'
forbidden_dep="github.com/(${legacy_org}|ZoneCNH)/x[.]go"
if GOWORK=off "$GO" list -deps ./... | rg -n "$forbidden_dep"; then
  echo "postgresx must not depend on the forbidden downstream application module" >&2
  exit 1
fi

bash ./scripts/ci/release_evidence_check.sh "$version"

echo "release check passed for $version"
