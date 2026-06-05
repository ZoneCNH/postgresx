#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

VERSION="${VERSION:-v0.1.0}"
OUT_DIR="release/manifest"
VERSION_FILE="$OUT_DIR/$VERSION.json"
LATEST_FILE="$OUT_DIR/latest.json"
mkdir -p "$OUT_DIR"

GO="${GO:-go}"
module="$(GOWORK=off "$GO" list -m)"
commit="$(git rev-parse --short HEAD 2>/dev/null || printf 'unknown')"
tree_sha="$(git rev-parse HEAD^{tree} 2>/dev/null || printf 'unknown')"
created_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
pgx_version="$(GOWORK=off "$GO" list -m -f '{{.Path}}@{{.Version}}' github.com/jackc/pgx/v5 2>/dev/null || printf 'github.com/jackc/pgx/v5@unknown')"

digest_file() {
  local path="$1"
  if [[ -f "$path" ]]; then
    sha256sum "$path" | awk '{print "sha256:" $1}'
  else
    printf 'missing'
  fi
}

source_digest="$(git ls-files -z '*.go' 'go.mod' 'go.sum' 'contracts/*' 'docs/api.md' 2>/dev/null | xargs -0 sha256sum 2>/dev/null | sha256sum | awk '{print "sha256:" $1}')"

python3 - "$VERSION" "$module" "$commit" "$tree_sha" "$source_digest" "$pgx_version" "$created_at" \
  "$(digest_file docs/api.md)" \
  "$(digest_file contracts/config.schema.json)" \
  "$(digest_file contracts/health.schema.json)" \
  "$(digest_file contracts/metrics.md)" \
  "$(digest_file contracts/error.schema.json)" >"$VERSION_FILE" <<'PY'
import json
import sys
(
    version,
    module,
    commit,
    tree_sha,
    source_digest,
    pgx_version,
    created_at,
    api_hash,
    config_hash,
    health_hash,
    metrics_hash,
    errors_hash,
) = sys.argv[1:]
manifest = {
    "schema_version": "1.0",
    "module": module,
    "package": "postgresx",
    "layer": "L2",
    "role": "postgresql_infrastructure_adapter",
    "standard_source": "github.com/ZoneCNH/xlib-standard",
    "release_level_target": "L2-T2",
    "release_level_actual": "L2-T2",
    "required_profiles": ["unit", "contract", "integration"],
    "release_allowed": False,
    "factory_grade_allowed": False,
    "min_score": 75,
    "version": version,
    "core_package": f"{module}/pkg/postgresx",
    "commit": commit,
    "tree_sha": tree_sha,
    "source_digest": source_digest,
    "provider_dependencies": {
        "pgx": pgx_version,
        "postgres_image": "postgres:16-alpine",
    },
    "capability_manifest": ".agent/l2-capabilities.yaml",
    "contract_pack_registry": ".agent/registry/l2-contract-packs.yaml",
    "l2_gate": ".agent/gates/l2gate.yaml",
    "hard_failures": [
        "secret_leak",
        "layer_violation",
        "missing_required_contract",
        "missing_required_evidence",
        "race_detected",
        "goroutine_leak",
        "release_level_overclaimed",
    ],
    "required_contract_tests": [
        "sql.exec",
        "sql.query_row",
        "sql.query_many",
        "sql.not_found",
        "sql.syntax_error",
        "sql.unique_violation",
        "sql.foreign_key_violation",
        "sql.context_timeout",
        "tx.commit",
        "tx.rollback",
        "tx.rollback_on_error",
        "pool.exhaustion",
    ],
    "boundaries": {
        "business_schema_or_orm": "forbidden",
        "configx_observex_core_dependency": "forbidden",
        "consumer_module_dependency": "forbidden",
        "core_env_file_secret_loading": "forbidden",
        "global_database_singleton": "forbidden",
    },
    "contract_hashes": {
        "api": api_hash,
        "config": config_hash,
        "health": health_hash,
        "metrics": metrics_hash,
        "errors": errors_hash,
    },
    "gates": {
        "fmt": "required",
        "vet": "required",
        "lint": "required",
        "test": "required",
        "race": "required",
        "boundary": "required",
        "contract": "required",
        "docs": "required",
        "security": "required",
        "integration": "required-or-skip-with-reason",
        "release_final": "required",
    },
    "evidence": [
        "GOWORK=off make vet",
        "GOWORK=off make test-unit",
        "GOWORK=off make test-contract",
        "GOWORK=off make test-integration",
        "GOWORK=off make boundary",
        "GOWORK=off make contracts",
        "GOWORK=off make secret-scan",
        "GOWORK=off make foundationx-api",
        "GOWORK=off make template-alignment",
        f"GOWORK=off make evidence VERSION={version}",
    ],
    "integration": {
        "status": "not_claimed",
        "evidence": "docs/evidence/20260601/postgres-integration.txt",
    },
    "downstream_adoption": {
        "status": "not_claimed",
        "evidence": "docs/CONSUMER_INTEGRATION.md",
    },
    "generated_at": created_at,
}
json.dump(manifest, sys.stdout, indent=2, sort_keys=True)
sys.stdout.write("\n")
PY

cp "$VERSION_FILE" "$LATEST_FILE"
sha256sum "$VERSION_FILE" | awk '{print $1}' >"$VERSION_FILE.sha256"
sha256sum "$LATEST_FILE" | awk '{print $1}' >"$LATEST_FILE.sha256"

echo "$VERSION_FILE"
echo "$VERSION_FILE.sha256"
echo "$LATEST_FILE"
echo "$LATEST_FILE.sha256"
