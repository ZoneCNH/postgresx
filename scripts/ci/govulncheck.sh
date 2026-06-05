#!/usr/bin/env bash
set -euo pipefail

enable="${XLIB_ENABLE_VULNCHECK:-0}"
force="${XLIB_FORCE_VULNCHECK:-0}"
interval_hours="${XLIB_VULNCHECK_INTERVAL_HOURS:-168}"
state_file="${XLIB_VULNCHECK_STATE:-.cache/security/govulncheck-last-run}"

if [[ "$enable" != "1" && "$force" != "1" ]]; then
  echo "govulncheck disabled; set XLIB_ENABLE_VULNCHECK=1 or XLIB_FORCE_VULNCHECK=1 to run it"
  exit 0
fi

if ! [[ "$interval_hours" =~ ^[0-9]+$ ]] || [[ "$interval_hours" -eq 0 ]]; then
  echo "XLIB_VULNCHECK_INTERVAL_HOURS must be a positive integer" >&2
  exit 2
fi

if [[ "$force" != "1" && -f "$state_file" ]]; then
  now="$(date +%s)"
  last="$(stat -c %Y "$state_file")"
  interval_seconds=$((interval_hours * 3600))
  if (( now - last < interval_seconds )); then
    echo "govulncheck skipped; last successful run is within ${interval_hours}h"
    exit 0
  fi
fi

if ! command -v govulncheck >/dev/null 2>&1; then
  echo "govulncheck is required; install with: GOWORK=off go install golang.org/x/vuln/cmd/govulncheck@v1.1.4" >&2
  exit 1
fi

GOWORK=off govulncheck ./...
mkdir -p "$(dirname "$state_file")"
touch "$state_file"
