# Release Manifest — postgresx v1.0.0

## Identity

- Module: `github.com/ZoneCNH/postgresx`
- Core package: `github.com/ZoneCNH/postgresx/pkg/postgresx`
- Status: published `v1.0.0` tag plus post-tag local `L2-T3 / 85`
  evidence on the `postgresx` branch
- Go verification mode: `GOWORK=off`

## Publication evidence

- GitHub release: `https://github.com/ZoneCNH/postgresx/releases/tag/v1.0.0`
- Tag object: `refs/tags/v1.0.0` resolves to `5c3e3a6`.
- Tag commit: `refs/tags/v1.0.0^{}` resolves to `310a249`.
- Release snapshot commit: `release/manifest/v1.0.0.json` records `7fe4cfd`.
- Release snapshot tree: `231164546c9c7b11d30287f1318c5f6a3b51442d`.
- Release metadata `targetCommitish` remains `main`; the tag object and resolved
  tag commit are the authoritative published release identity.

## Local gate evidence

- Current local gate: `L2-T3 / 85`
- Local release decision: `release_allowed=true`
- Factory decision: `factory_grade_allowed=false`
- Required profiles: unit, contract, integration, chaos, benchmark, and local
  downstream compile smoke
- Executable check: `GOWORK=off VERSION=v1.0.0 make release-check`
- Evidence generated at: `2026-06-13T07:52:51Z`

The downstream smoke proves import, compile, configuration, and `Queryer`
boundary compatibility from a temporary consumer module. It is not production
consumer adoption evidence, so adoption remains smoke-level rather than
factory-grade evidence.

The integration evidence was produced with a real PostgreSQL development DSN
loaded from the local SRE secret document through environment variables only.
The DSN and credentials are not written to release docs, manifests, or evidence
artifacts.

## Manifest semantics

`release/manifest/v1.0.0.json` and `release/manifest/latest.json` record the
source snapshot used to generate the post-tag release evidence. Their current
`commit` field is `7fe4cfd`, which resolves in Git and is expected to be an
ancestor of the evidence-carrying `postgresx` branch after this
documentation/evidence commit, but it is **not** an ancestor of the immutable
`v1.0.0` tag commit `310a249`. The local
`release-evidence-check` therefore fails until the release-history decision is
reconciled by an approved action, such as regenerating the manifest from the
tagged snapshot or cutting a successor release tag. Do not rewrite or retag
`v1.0.0` without explicit release-history approval.

## Included surfaces

- caller-owned `Config` construction, validation, DSN rendering, and redaction;
- `New`/`Open` client lifecycle, ping/check, pool stats, and query methods;
- sqlc-compatible `Queryer` execution boundary;
- `WithTx` and `WithTxOptions` transaction helpers;
- caller-owned migration runner over embedded or filesystem SQL sources;
- foundationx-compatible error mapping, including SQLSTATE `42P01`
  `undefined_table` to `not_found`;
- logging and metrics options;
- examples that keep env and secret loading outside the core package.

## Explicit non-claims

- This snapshot does not include production downstream adoption proof beyond
  the local compile smoke.
- This snapshot does not prove current GitHub Actions status or production
  soak.
- This snapshot does not publish or embed PostgreSQL credentials or DSNs.
- This module does not own domain schema, repositories, application services, or
  production DSNs.

## Evidence index

- `release/manifest/v1.0.0.json`
- `release/manifest/v1.0.0.json.sha256`
- `release/manifest/latest.json`
- `release/manifest/latest.json.sha256`
- `docs/EVIDENCE-20260601.md`
- `docs/evidence/20260601/go-test.txt`
- `docs/evidence/20260601/go-test-race.txt`
- `docs/evidence/20260601/go-vet.txt`
- `docs/evidence/20260601/secret-scan.txt`
- `docs/evidence/20260601/no-consumer-deps.txt`
- `docs/evidence/20260601/dependencies.txt`
- `.agent/evidence/raw/chaos-test.json`
- `.agent/evidence/raw/benchmark-smoke.json`
- `.agent/evidence/raw/downstream-smoke.json`
- `.agent/evidence/normalized/chaos-check.json`
- `.agent/evidence/normalized/benchmark-check.json`
- `.agent/evidence/normalized/adoption-check.json`
- `.agent/evidence/decision/release-readiness.json`

## Remaining release hardening

Before treating this as a factory-grade or production-adopted release, add fresh
external CI evidence, production soak evidence, and consumer adoption evidence
from a current consumer checkout. GitHub Actions is currently blocked outside
the repository by an account billing lock, so local evidence is the authoritative
available gate evidence. The current release-evidence blocker is the manifest
`7fe4cfd` versus tag `310a249` ancestry mismatch; resolve it through an
approved release-history action rather than rewriting or retagging `v1.0.0`.
