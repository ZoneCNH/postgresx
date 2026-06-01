#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! rg -n 'github.com/ZoneCNH/foundationx/pkg/foundationx' pkg/postgresx contracts testkit examples >/dev/null; then
  echo "foundationx API import was not found in expected package surfaces" >&2
  exit 1
fi

if rg -n 'github.com/ZoneCNH/(configx|observex)|github.com/bytechainx/(configx|observex)' pkg/postgresx; then
  echo "core package must not import configx or observex" >&2
  exit 1
fi

GOWORK=off go list -deps ./pkg/postgresx | rg -q '^github.com/ZoneCNH/foundationx/pkg/foundationx$' || {
  echo "pkg/postgresx does not resolve foundationx API dependency" >&2
  exit 1
}

echo "foundationx API check passed"
