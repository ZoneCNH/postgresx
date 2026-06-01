#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

version="${1:-${VERSION:-v0.1.0}}"
plain_version="${version#v}"
required=(
  "docs/RELEASE_MANIFEST-${version}.md"
  "docs/EVIDENCE-20260601.md"
  "docs/RETROSPECTIVE-GOAL-20260601-001.md"
  "release/manifest/latest.json"
  "release/manifest/${version}.json"
  ".agent/postgresx-v0.1.0.md"
  "contracts/config.schema.json"
  "contracts/error.schema.json"
  "contracts/health.schema.json"
)

for path in "${required[@]}"; do
  if [[ ! -s "$path" ]]; then
    echo "missing required release evidence: $path" >&2
    exit 1
  fi
done

python3 - "$version" "$plain_version" <<'PY'
import json, sys
from pathlib import Path
version, plain = sys.argv[1:3]
latest = json.loads(Path('release/manifest/latest.json').read_text())
versioned = json.loads(Path(f'release/manifest/{version}.json').read_text())
if latest != versioned:
    raise SystemExit('latest release manifest differs from versioned manifest')
if latest.get('version') != version or latest.get('module') != 'github.com/ZoneCNH/postgresx':
    raise SystemExit('release manifest has wrong version or module')
if latest.get('core_package') != 'github.com/ZoneCNH/postgresx/pkg/postgresx':
    raise SystemExit('release manifest has wrong core package')
if latest.get('go_module_version') not in (plain, version):
    raise SystemExit('release manifest has wrong Go module version')
PY

legacy_org='byte''chainx'
if rg -n "${legacy_org}|x[.]go|/home/k8s/secrets|production[.]yaml|config[.]local[.]yaml" \
  README.md docs release .agent contracts examples testkit pkg \
  --glob '!docs/goal.md'; then
  echo "release evidence contains forbidden legacy or implicit-secret references" >&2
  exit 1
fi

echo "release evidence check passed for $version"
