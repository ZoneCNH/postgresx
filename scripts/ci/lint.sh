#!/usr/bin/env bash
set -euo pipefail

GO="${GO:-go}"
GOWORK=off "$GO" vet ./...
