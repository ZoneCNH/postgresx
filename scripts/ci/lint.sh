#!/usr/bin/env bash
set -euo pipefail

GO="${GO:-go}"
GOWORK=off "$GO" vet ./...

if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "golangci-lint is required; install it or run inside the Docker toolchain" >&2
  exit 1
fi

GOWORK=off golangci-lint run ./...
