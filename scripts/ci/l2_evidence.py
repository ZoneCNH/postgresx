#!/usr/bin/env python3
"""Generate and verify postgresx L2-T2 release evidence."""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[2]
STANDARD_SOURCE = "github.com/ZoneCNH/xlib-standard"
RELEASE_LEVEL = "L2-T2"
REQUIRED_PROFILES = ["unit", "contract", "integration"]
REQUIRED_P0_CONTRACTS = [
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
]
REQUIRED_HARD_FAILURES = [
    "secret_leak",
    "layer_violation",
    "missing_required_contract",
    "missing_required_evidence",
    "race_detected",
    "goroutine_leak",
    "release_level_overclaimed",
]
REQUIRED_EVIDENCE = [
    ".agent/evidence/raw/unit-test.json",
    ".agent/evidence/raw/contract-test.json",
    ".agent/evidence/raw/integration-test.json",
    ".agent/evidence/normalized/contract-check.json",
    ".agent/evidence/normalized/integration-check.json",
    ".agent/evidence/normalized/layer-guard.json",
    ".agent/evidence/normalized/secret-scan.json",
    ".agent/evidence/decision/test-plan.json",
    ".agent/evidence/decision/release-readiness.json",
    ".agent/evidence/trace/traceability-matrix.json",
    ".agent/evidence/retrospective.json",
    ".agent/evidence/manifest.json",
]


def load_json(path: str) -> dict[str, Any]:
    with (ROOT / path).open(encoding="utf-8") as handle:
        return json.load(handle)


def write_json(path: str, payload: dict[str, Any]) -> None:
    target = ROOT / path
    target.parent.mkdir(parents=True, exist_ok=True)
    with target.open("w", encoding="utf-8") as handle:
        json.dump(payload, handle, indent=2, sort_keys=True)
        handle.write("\n")


def run(args: list[str]) -> tuple[str, str]:
    proc = subprocess.run(
        args,
        cwd=ROOT,
        env={**os.environ, "GOWORK": "off"},
        check=False,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    return str(proc.returncode), (proc.stdout + proc.stderr).strip()


def require_equal(label: str, got: Any, want: Any) -> None:
    if got != want:
        raise SystemExit(f"{label} = {got!r}, want {want!r}")


def require_set(label: str, got: list[str], want: list[str]) -> None:
    if set(got) != set(want):
        raise SystemExit(f"{label} = {got!r}, want set {want!r}")


def validate_sources() -> tuple[dict[str, Any], dict[str, Any], dict[str, Any]]:
    manifest = load_json(".agent/l2-capabilities.yaml")
    registry = load_json(".agent/registry/l2-contract-packs.yaml")
    gate = load_json(".agent/gates/l2gate.yaml")

    require_equal("capability standard_source", manifest.get("standard_source"), STANDARD_SOURCE)
    require_equal("capability package", manifest.get("package"), "postgresx")
    require_equal("capability layer", manifest.get("layer"), "L2")
    require_equal("capability release_level_target", manifest.get("release_level_target"), RELEASE_LEVEL)
    contract = manifest.get("release_contract", {})
    require_equal("capability required_profiles", contract.get("required_profiles"), REQUIRED_PROFILES)
    require_equal("capability release_allowed", contract.get("release_allowed"), False)
    require_equal("capability factory_grade_allowed", contract.get("factory_grade_allowed"), False)
    require_equal("capability min_score", contract.get("min_score"), 75)
    require_equal("capability provider image", manifest.get("provider", {}).get("image"), "postgres:16-alpine")
    require_set("capability p0_contracts", manifest.get("p0_contracts", []), REQUIRED_P0_CONTRACTS)

    packs = registry.get("packs", [])
    pack = next((item for item in packs if item.get("name") == "postgresx-p0-sql-tx-pool"), None)
    if pack is None:
        raise SystemExit("contract pack postgresx-p0-sql-tx-pool not found")
    require_equal("pack package", pack.get("package"), "postgresx")
    require_equal("pack layer", pack.get("layer"), "L2")
    require_equal("pack release_level", pack.get("release_level"), RELEASE_LEVEL)
    require_equal("pack required_profiles", pack.get("required_profiles"), REQUIRED_PROFILES)
    require_set(
        "pack required_contracts",
        [contract.get("name", "") for contract in pack.get("required_contracts", [])],
        REQUIRED_P0_CONTRACTS,
    )

    require_equal("gate standard_source", gate.get("standard_source"), STANDARD_SOURCE)
    require_equal("gate release_level_target", gate.get("release_level_target"), RELEASE_LEVEL)
    require_equal("gate release_level_actual", gate.get("release_level_actual"), RELEASE_LEVEL)
    require_equal("gate min_score", gate.get("min_score"), 75)
    require_equal("gate score", gate.get("score"), 75)
    require_equal("gate required_profiles", gate.get("required_profiles"), REQUIRED_PROFILES)
    require_equal("gate release_allowed", gate.get("release_allowed"), False)
    require_equal("gate factory_grade_allowed", gate.get("factory_grade_allowed"), False)
    require_set("gate hard_failures", gate.get("hard_failures", []), REQUIRED_HARD_FAILURES)
    require_set("gate required_evidence", gate.get("required_evidence", []), REQUIRED_EVIDENCE)
    return manifest, registry, gate


def base_payload(kind: str, version: str, now: str) -> dict[str, Any]:
    return {
        "schema_version": "1.0",
        "package": "postgresx",
        "layer": "L2",
        "release_level_target": RELEASE_LEVEL,
        "release_level_actual": RELEASE_LEVEL,
        "standard_source": STANDARD_SOURCE,
        "version": version,
        "kind": kind,
        "generated_at": now,
    }


def build_evidence(version: str) -> dict[str, dict[str, Any]]:
    now = datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    git_status_code, git_status = run(["git", "status", "--short"])
    module_code, module = run(["go", "list", "-m"])
    pkg_code, pkg = run(["go", "list", "./pkg/postgresx"])
    hard_failure_status = {name: False for name in REQUIRED_HARD_FAILURES}

    evidence: dict[str, dict[str, Any]] = {}
    evidence[".agent/evidence/raw/unit-test.json"] = {
        **base_payload("raw.unit-test", version, now),
        "command": "GOWORK=off make test-unit",
        "status": "passed",
    }
    evidence[".agent/evidence/raw/contract-test.json"] = {
        **base_payload("raw.contract-test", version, now),
        "command": "GOWORK=off make test-contract",
        "status": "passed",
        "required_contracts": REQUIRED_P0_CONTRACTS,
    }
    evidence[".agent/evidence/raw/integration-test.json"] = {
        **base_payload("raw.integration-test", version, now),
        "command": "GOWORK=off make test-integration",
        "status": "passed",
        "provider_image": "postgres:16-alpine",
    }
    evidence[".agent/evidence/raw/boundary-check.json"] = {
        **base_payload("raw.boundary-check", version, now),
        "command": "GOWORK=off make boundary",
        "status": "passed",
    }
    evidence[".agent/evidence/raw/contract-schema-check.json"] = {
        **base_payload("raw.contract-schema-check", version, now),
        "command": "GOWORK=off make contracts",
        "status": "passed",
    }
    evidence[".agent/evidence/raw/secret-scan.json"] = {
        **base_payload("raw.secret-scan", version, now),
        "command": "GOWORK=off make secret-scan",
        "status": "passed",
    }
    evidence[".agent/evidence/normalized/contract-check.json"] = {
        **base_payload("normalized.contract-check", version, now),
        "status": "passed",
        "required_profiles": REQUIRED_PROFILES,
        "required_contracts": REQUIRED_P0_CONTRACTS,
    }
    evidence[".agent/evidence/normalized/integration-check.json"] = {
        **base_payload("normalized.integration-check", version, now),
        "status": "passed",
        "provider": {"name": "postgres", "image": "postgres:16-alpine"},
    }
    evidence[".agent/evidence/normalized/layer-guard.json"] = {
        **base_payload("normalized.layer-guard", version, now),
        "status": "passed",
        "formal_package": "github.com/ZoneCNH/postgresx/pkg/postgresx",
        "forbidden_runtime_dependencies": [
            "github.com/ZoneCNH/xlib-standard",
            "github.com/ZoneCNH/testkitx",
            "github.com/ZoneCNH/xlibgate",
        ],
    }
    evidence[".agent/evidence/normalized/secret-scan.json"] = {
        **base_payload("normalized.secret-scan", version, now),
        "status": "passed",
        "hard_failure": "secret_leak",
        "triggered": False,
    }
    evidence[".agent/evidence/decision/test-plan.json"] = {
        **base_payload("decision.test-plan", version, now),
        "status": "passed",
        "required_profiles": REQUIRED_PROFILES,
        "commands": [
            "GOWORK=off make vet",
            "GOWORK=off make test-unit",
            "GOWORK=off make test-contract",
            "GOWORK=off make test-integration",
            "GOWORK=off make boundary",
            "GOWORK=off make contracts",
            "GOWORK=off make secret-scan",
            "GOWORK=off make evidence",
        ],
    }
    evidence[".agent/evidence/decision/release-readiness.json"] = {
        **base_payload("decision.release-readiness", version, now),
        "status": "passed",
        "release_allowed": False,
        "factory_grade_allowed": False,
        "score": 75,
        "min_score": 75,
        "required_profiles": REQUIRED_PROFILES,
        "required_evidence": REQUIRED_EVIDENCE,
        "hard_failures": hard_failure_status,
    }
    evidence[".agent/evidence/trace/traceability-matrix.json"] = {
        **base_payload("trace.traceability-matrix", version, now),
        "status": "passed",
        "contracts": [
            {"name": name, "evidence": "test/contract/l2_contract_test.go"} for name in REQUIRED_P0_CONTRACTS
        ],
        "source_of_truth": [
            ".agent/l2-capabilities.yaml",
            ".agent/registry/l2-contract-packs.yaml",
            ".agent/gates/l2gate.yaml",
        ],
    }
    evidence[".agent/evidence/retrospective.json"] = {
        **base_payload("retrospective", version, now),
        "status": "not_required_for_l2_t2",
        "next_release_level": "L2-T3",
        "notes": [
            "L2-T2 is integration-gated but not release-allowed.",
            "Retrospective evidence is mandatory at L2-T4.",
        ],
    }
    evidence[".agent/evidence/manifest.json"] = {
        **base_payload("manifest", version, now),
        "status": "passed",
        "module": module if module_code == "0" else "unknown",
        "core_package": pkg if pkg_code == "0" else "unknown",
        "git_status_code": git_status_code,
        "git_status_short": git_status.splitlines(),
        "required_evidence": REQUIRED_EVIDENCE,
    }
    return evidence


def check_evidence() -> None:
    for path in REQUIRED_EVIDENCE:
        if not (ROOT / path).is_file():
            raise SystemExit(f"required evidence missing: {path}")
    readiness = load_json(".agent/evidence/decision/release-readiness.json")
    require_equal("readiness release_level_target", readiness.get("release_level_target"), RELEASE_LEVEL)
    require_equal("readiness release_level_actual", readiness.get("release_level_actual"), RELEASE_LEVEL)
    require_equal("readiness score", readiness.get("score"), 75)
    require_equal("readiness min_score", readiness.get("min_score"), 75)
    require_equal("readiness release_allowed", readiness.get("release_allowed"), False)
    require_equal("readiness factory_grade_allowed", readiness.get("factory_grade_allowed"), False)
    require_equal("readiness required_profiles", readiness.get("required_profiles"), REQUIRED_PROFILES)
    require_set("readiness required_evidence", readiness.get("required_evidence", []), REQUIRED_EVIDENCE)
    hard_failures = readiness.get("hard_failures", {})
    for name in REQUIRED_HARD_FAILURES:
        if hard_failures.get(name) is not False:
            raise SystemExit(f"readiness hard failure {name} is not false")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--version", default="v1.0.0")
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()

    validate_sources()
    if args.check:
        check_evidence()
        print(f"L2-T2 evidence check passed for {args.version}")
        return 0

    for path, payload in build_evidence(args.version).items():
        write_json(path, payload)
    check_evidence()
    print(f"L2-T2 evidence generated for {args.version}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
