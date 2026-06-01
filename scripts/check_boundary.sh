#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v rg >/dev/null 2>&1; then
  echo "rg is required for boundary checks" >&2
  exit 1
fi

GO="${GO:-go}"
status=0
legacy_org='byte''chainx'
forbidden_dep="github.com/(${legacy_org}|ZoneCNH)/x[.]go"

if GOWORK=off "$GO" list -deps ./... | rg -n 'github.com/([b]ytechainx|ZoneCNH)/x\.go'; then
  echo "boundary violation: postgresx must not depend on application module" >&2
  status=1
fi

scan_paths=()
for dir in pkg contracts internal examples testkit; do
  if [[ -e "$dir" ]]; then
    scan_paths+=("$dir")
  fi
done

if [[ "${#scan_paths[@]}" -gt 0 ]] && rg -n \
  'MacroRegime|MarketRegime|TradingSignal|BTCUSDT|ETHUSDT|Kline|OrderBook|Position|RiskGate|MarketData|MacroData' \
  "${scan_paths[@]}"; then
  echo "boundary violation: business-domain terms found in postgresx library code" >&2
  status=1
fi

if rg -n 'github[.]com/[b]ytechainx|github[.]com/ZoneCNH/postgresx/pkg/postgresx/(examples|contracts)' \
  go.mod go.sum pkg contracts internal examples testkit scripts .github README.md \
  --glob '!docs/goal.md'; then
  echo "boundary violation: stale module/package reference found" >&2
  status=1
fi

if rg -n 'configx|observex' pkg/postgresx; then
  echo "boundary violation: core package must not depend on configx or observex" >&2
  status=1
fi

if rg -n 'os[.]Getenv|os[.]LookupEnv|godotenv|production[.]yaml|config[.]local[.]yaml|/home/k8s/secrets' pkg/postgresx; then
  echo "boundary violation: core package must not load secrets from env or files implicitly" >&2
  status=1
fi

if [[ "$status" -ne 0 ]]; then
  exit "$status"
fi

echo "boundary check passed"
