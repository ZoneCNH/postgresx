#!/usr/bin/env bash
set -euo pipefail

if ! command -v govulncheck >/dev/null 2>&1; then
  echo "govulncheck is required; install with: GOWORK=off go install golang.org/x/vuln/cmd/govulncheck@latest" >&2
  exit 1
fi

GOWORK=off govulncheck ./...
