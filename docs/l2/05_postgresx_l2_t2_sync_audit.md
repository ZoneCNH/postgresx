# Postgrex L2-T2 Sync Audit

Date: 2026-06-05
Worker: `worker-1`
Scope: manifest, contract-pack registry, release-level semantics, and evidence fields for the L2-T2 MVA.

## Inputs reviewed

- Team context: `/home/postgresx/.omx/context/postgresx-l2-factory-20260605T165052Z.md`
- Current release manifests: `release/manifest/v0.1.0.json`, `release/manifest/latest.json`
- Current manifest generator/checkers: `scripts/generate_manifest.sh`, `scripts/ci/release_evidence_check.sh`, `scripts/ci/release_check.sh`, `Makefile`
- Current contract surfaces: `contracts/*`, `docs/api.md`
- Current release evidence docs: `docs/RELEASE_MANIFEST-v0.1.0.md`, `docs/EVIDENCE-20260601.md`, `docs/RETROSPECTIVE-GOAL-20260601-001.md`, `docs/VERSION_MATRIX.md`

The task-referenced plan file `docs/l2/05_postgresx_execution_plan.md` is not present in this worker checkout or under `/home/postgresx`; this audit uses the OMX context file and current repository state as the available source of truth.

## Standard-owned semantics to preserve

Postgrex may declare L2-T2 readiness, but the release-level semantics must remain owned by `github.com/ZoneCNH/xlib-standard`, not redefined locally.

Required L2-T2 semantics from the team context:

- `release_level: L2-T2`
- `required_profiles: [unit, contract, integration]`
- `release_allowed: false`
- `factory_grade_allowed: false`
- `min_score: 75`
- hard-failure names:
  - `secret_leak`
  - `layer_violation`
  - `missing_required_contract`
  - `missing_required_evidence`
  - `race_detected`
  - `goroutine_leak`
  - `release_level_overclaimed`

## Required synchronization matrix

| Surface | Current state | Required sync |
| --- | --- | --- |
| `.agent/l2-capabilities.yaml` | Missing. `.agent/` contains only narrative stubs. | Add the capability manifest declaring Postgrex as an L2 PostgreSQL infrastructure adapter, provider `postgres`, supported SQL/transaction/pool capabilities, `standard_source: github.com/ZoneCNH/xlib-standard`, and `release_level: L2-T2`. It must reference standard-owned release semantics instead of copying custom local semantics into product code. |
| `.agent/registry/l2-contract-packs.yaml` | Missing. Current repo uses `contract_hashes` in release manifests only. | Add the contract-pack registry for required L2-T2 contract tests. It should map required packs to `test/contract/l2_contract_test.go`, existing contract docs/schemas, and provider adapter behavior without replacing `contract_hashes` until manifest/checker support is updated. |
| `.agent/gates/l2gate.yaml` | Missing. Current gates are legacy string values under manifest `gates`. | Add xlibgate-oriented gate config with the required profiles, hard-failure names, required evidence paths, required contract tests, min score, and decision output path `.agent/evidence/decision/release-readiness.json`. |
| `.agent/evidence/decision/release-readiness.json` | Missing. | Generate or check in the xlibgate decision output containing `release_level_actual`, `hard_failures`, `required_contract_tests`, `required_evidence`, score/min-score comparison, and release/factory-grade booleans. |
| `.agent/evidence/trace/traceability-matrix.json` | Missing. | Add requirement-to-evidence traceability linking capability manifest entries, contract-pack tests, Make targets, evidence artifacts, and release manifest fields. |
| `.agent/evidence/retrospective.json` | Missing. | Add structured retrospective evidence capturing the L2-T2 sync result, unresolved gaps, and standard-sync status. |
| `release/manifest/v0.1.0.json` and `release/manifest/latest.json` | Present, generic `layer: L2`. They include `standard_source`, `contract_hashes`, legacy `gates`, and legacy evidence commands. | Extend generator and manifests to include L2-T2 declaration and decision fields: `release_level`, `release_level_actual`, `required_profiles`, `release_allowed`, `factory_grade_allowed`, `min_score`, `hard_failures`, `required_contract_tests`, `required_evidence`, and references to `.agent/evidence/...`. Update matching `.sha256` files after any manifest change. |
| `scripts/generate_manifest.sh` | Generates legacy manifest fields only. | Teach the generator to emit the L2-T2 fields above from the manifest/gate definitions so hand-edited manifests do not drift. |
| `scripts/ci/release_evidence_check.sh` | Checks legacy release evidence, checksums, base manifest identity, and contract hashes. | Add validation for `.agent/l2-capabilities.yaml`, `.agent/registry/l2-contract-packs.yaml`, `.agent/gates/l2gate.yaml`, `.agent/evidence/decision/release-readiness.json`, `.agent/evidence/trace/traceability-matrix.json`, `.agent/evidence/retrospective.json`, and the manifest fields listed above. |
| `Makefile` | Missing expected MVA targets `l2-plan`, `test-unit`, `test-contract`, and `test-integration`; has legacy `evidence`, `release-evidence-check`, and `release-check`. | Add the expected targets and ensure `release-check` validates generated L2 evidence, required contract tests, required evidence, hard failures, and release-level overclaim behavior. |
| `test/contract/l2_contract_test.go` | Missing; no top-level `test/` directory exists. | Add L2 contract tests required by the contract-pack registry. This should complement existing `contracts/contracts_test.go`, not remove it. |
| `test/postgresxtest/` | Missing. | Add provider fixtures/helpers used by integration-profile L2 tests. |
| `docs/RELEASE_MANIFEST-v0.1.0.md`, `docs/EVIDENCE-20260601.md`, `docs/RETROSPECTIVE-GOAL-20260601-001.md`, `docs/VERSION_MATRIX.md` | Present but partly stale against current manifests and `go.mod`. | Update narrative release docs after generated manifests/evidence are updated. Known stale items: checksum/contract-hash language now conflicts with manifest artifacts; version matrix differs from `go.mod` for Go, pgx, and foundationx. |

## Current manifest fields that must remain synchronized

The existing manifest generator/checker already depends on these fields and they must stay stable while L2-T2 fields are added:

- `schema_version`
- `module`
- `package`
- `layer`
- `role`
- `standard_source`
- `version`
- `core_package`
- `provider_dependencies`
- `boundaries`
- `contract_hashes`
- `gates`
- `evidence`
- `integration.status` and `integration.evidence`
- `downstream_adoption.status` and `downstream_adoption.evidence`
- `commit`
- `tree_sha`
- `source_digest`
- `generated_at`

## Must-not-sync constraints

- Do not import `xlib-standard`, `testkitx`, or `xlibgate` into public library/runtime packages. L2 tooling may be used by tests, gates, and evidence generation only.
- Do not redefine L2-T2 semantics in Postgrex business logic. Postgrex should declare its level and point to the standard-owned contract/gate definitions.
- Do not edit release manifest JSON without regenerating matching `.sha256` files and rerunning release evidence checks.

## Risks and blockers

1. `docs/l2/05_postgresx_execution_plan.md` is absent, so the exact plan-level acceptance criteria cannot be audited directly from the requested document.
2. Required `.agent` YAML/JSON evidence surfaces are absent, so `release_level_actual`, `hard_failures`, `required_contract_tests`, and `required_evidence` cannot currently be proven.
3. Current release manifests only claim generic `layer: L2`; they do not encode L2-T2 fields or xlibgate decision output.
4. Current `release_evidence_check` can pass while all L2-T2 `.agent/evidence` files are missing.
5. Missing Makefile targets and missing top-level L2 contract tests mean the required unit/contract/integration profiles are not yet enforceable by the expected target names.
6. Existing release docs contain stale statements about missing checksums/contract hashes and stale dependency versions; these should be reconciled after generator/checker changes to avoid release evidence drift.
7. Integration evidence may require Docker or DSN availability. If unavailable, the evidence path must record an explicit skip reason rather than silently passing.

## Recommended handoff order

1. Add `.agent/l2-capabilities.yaml`, `.agent/registry/l2-contract-packs.yaml`, and `.agent/gates/l2gate.yaml` from xlib-standard templates.
2. Add `test/contract/l2_contract_test.go` and `test/postgresxtest/` fixtures to satisfy the registry.
3. Update `scripts/generate_manifest.sh`, `scripts/ci/release_evidence_check.sh`, and Makefile MVA targets so generated manifests and release checks validate the new fields.
4. Generate `.agent/evidence/decision/release-readiness.json`, `.agent/evidence/trace/traceability-matrix.json`, and `.agent/evidence/retrospective.json`.
5. Regenerate release manifests/checksums, then update release docs to match the generated artifacts.

## Verification commands for the final integrated MVA

```bash
GOWORK=off make test-unit
GOWORK=off make test-contract
GOWORK=off make test-integration
GOWORK=off make evidence VERSION=v0.1.0
GOWORK=off make release-evidence-check VERSION=v0.1.0
GOWORK=off make release-check VERSION=v0.1.0
```

Additional dependency-boundary check:

```bash
GOWORK=off go list -deps ./... | rg -n 'xlib-standard|testkitx|xlibgate'
```

This command should return no matches for public/runtime dependency closure.
