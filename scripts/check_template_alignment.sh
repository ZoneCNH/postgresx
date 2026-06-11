#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

status=0
required_paths=(
  "go.mod"
  "pkg/postgresx"
  "contracts/config.schema.json"
  "contracts/error.schema.json"
  "contracts/health.schema.json"
  "contracts/metrics.md"
  "examples/basic/main.go"
  "testkit/postgres.go"
  "docs/EVIDENCE-20260601.md"
  "docs/RELEASE_MANIFEST-v0.1.0.md"
  "docs/RETROSPECTIVE-GOAL-20260601-001.md"
  ".agent"
  ".devcontainer/devcontainer.json"
  ".dockerignore"
  ".github/CODEOWNERS"
  ".github/ISSUE_TEMPLATE/contract-change.yml"
  ".github/ISSUE_TEMPLATE/generator-regression.yml"
  ".github/ISSUE_TEMPLATE/release-evidence-gap.yml"
  ".github/ISSUE_TEMPLATE/standard-doc-gap.yml"
  ".github/dependabot.yml"
  ".github/pull_request_template.md"
  ".github/rulesets/protect-main.json"
  ".github/rulesets/protect-release-tags.json"
  ".github/workflows/docker-contract.yml"
  ".github/workflows/security.yml"
  ".golangci.yml"
  ".githooks/pre-commit"
  ".githooks/pre-push"
  ".tool-versions"
  "Dockerfile"
  "docker-compose.yml"
  "release/manifest"
  "renovate.json"
  "scripts/docker/check_toolchain.sh"
  "scripts/docker/docker_gate.sh"
  "scripts/docker/prefetch_tools.sh"
)

for path in "${required_paths[@]}"; do
  if [[ ! -e "$path" ]]; then
    echo "missing template-required path: $path" >&2
    status=1
  fi
done

module="$(GOWORK=off go list -m)"
if [[ "$module" != "github.com/ZoneCNH/postgresx" ]]; then
  echo "unexpected module path: $module" >&2
  status=1
fi

if find . -maxdepth 1 -name '*.go' -print | grep -q .; then
  echo "root package Go files are not allowed; core must live in pkg/postgresx" >&2
  status=1
fi

required_make_targets=(
  "build"
  "build-check"
  "govulncheck"
  "security"
  "docker-toolchain-check"
  "docker-build"
  "docker-build-check"
  "docker-ci"
  "docker-release-check"
  "docker-release-final-check"
  "docker-runtime-check"
  "docker-drift-check"
  "docker-contract"
)

for target in "${required_make_targets[@]}"; do
  if ! rg -q "^${target}:" Makefile; then
    echo "missing template-required Make target: $target" >&2
    status=1
  fi
done

go_version="$(awk '$1 == "go" { print $2; exit }' go.mod)"
go_minor="${go_version%.*}"
if [[ -z "$go_version" || -z "$go_minor" ]]; then
  echo "could not determine Go version from go.mod" >&2
  status=1
else
  if ! rg -q "^ARG GO_VERSION=${go_minor}$" Dockerfile; then
    echo "Dockerfile GO_VERSION must match go.mod minor version: $go_minor" >&2
    status=1
  fi
  if ! rg -Fq 'GO_VERSION: ${GO_VERSION:-'"${go_minor}"'}' docker-compose.yml; then
    echo "docker-compose.yml GO_VERSION default must match go.mod minor version: $go_minor" >&2
    status=1
  fi
  if ! rg -Fq 'go_version="${GO_VERSION:-'"${go_minor}"'}"' scripts/docker/docker_gate.sh; then
    echo "docker gate GO_VERSION default must match go.mod minor version: $go_minor" >&2
    status=1
  fi
  if ! rg -Fq "golang ${go_version}" .tool-versions; then
    echo ".tool-versions Go version must match go.mod: $go_version" >&2
    status=1
  fi
fi

stale_module_path='github.com/bytechainx''/postgresx'
nested_non_core_path='github.com/ZoneCNH/postgresx/pkg/postgresx/''(examples|testkit|contracts|docs|internal)'
if rg -n "$stale_module_path|$nested_non_core_path" \
  --glob '!go.sum' \
  --glob '!docs/goal.md' \
  .; then
  echo "stale module/package path found" >&2
  status=1
fi

if [[ "$status" -ne 0 ]]; then
  exit "$status"
fi

echo "template alignment check passed"
