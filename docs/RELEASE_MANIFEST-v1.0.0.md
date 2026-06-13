# Release Manifest — postgresx v1.0.0

## Identity

- Module: `github.com/ZoneCNH/postgresx`
- Core package: `github.com/ZoneCNH/postgresx/pkg/postgresx`
- Status: published immutable `v1.0.0` tag with restored tag snapshot
  manifest metadata, plus separate post-tag local `L2-T3 / 85` evidence on
  the `postgresx` branch
- Go verification mode: `GOWORK=off`

## Publication evidence

- GitHub release: `https://github.com/ZoneCNH/postgresx/releases/tag/v1.0.0`
- Tag object: `refs/tags/v1.0.0` resolves to `5c3e3a6`.
- Tag commit: `refs/tags/v1.0.0^{}` resolves to `310a249`.
- Release snapshot commit metadata: `release/manifest/v1.0.0.json` records
  `9eaf770`.
- Release snapshot tree: `a45b1813f4ba5c0cb9a5b90e80b75f970078616b`.
- Snapshot ancestry blocker: commit `9eaf770` resolves in the local object
  database, but is not an ancestor of the current `HEAD` or the immutable tag
  commit `310a249`; the current `release-evidence-check` contract therefore
  rejects this manifest.
- Release metadata `targetCommitish` remains `main`; the tag object and resolved
  tag commit are the authoritative published release identity.

## Local gate evidence

- Tag manifest gate: `L2-T2 / 75`
- Tag manifest release decision: `release_allowed=false`
- Current branch evidence gate: `L2-T3 / 85`
- Current branch release decision: `release_allowed=true`
- Factory decision for both snapshots: `factory_grade_allowed=false`
- Current branch required profiles: unit, contract, integration, chaos,
  benchmark, and local downstream compile smoke
- Tag evidence check: `GOWORK=off VERSION=v1.0.0 make release-evidence-check`
  currently fails with `release manifest source commit is not an ancestor of
  HEAD: 9eaf770`
- Current branch evidence generator: `GOWORK=off VERSION=v1.0.0 make release-check`
- Tag manifest generated at: `2026-06-13T01:10:36Z`
- Current branch evidence generated at: `2026-06-13T07:52:51Z`

The downstream smoke proves import, compile, configuration, and `Queryer`
boundary compatibility from a temporary consumer module. It is not production
consumer adoption evidence, so adoption remains smoke-level rather than
factory-grade evidence.

The integration evidence was produced on 2026-06-13 with PostgreSQL connection
fields read from the local SRE secret Markdown document and assembled into a
DSN inside one shell process. That DSN was passed only through
`POSTGRESX_INTEGRATION_DSN` / `POSTGRES_TEST_DSN` with
`POSTGRESX_REQUIRE_INTEGRATION=1`; the raw DSN and credential-bearing values
are not written to release docs, manifests, or evidence artifacts.

## Manifest semantics

`release/manifest/v1.0.0.json` and `release/manifest/latest.json` intentionally
preserve the source metadata from the published `v1.0.0` snapshot instead of
silently replacing it with post-tag branch evidence. Their `commit` field is
`9eaf770`, which resolves in Git but is outside the current `HEAD` and
`v1.0.0` tag ancestry. This preserves the immutable snapshot semantics, but it
does not satisfy `release-evidence-check` until the release-history or manifest
contract is explicitly resolved.

Running `release-check` or `make evidence` on the current `postgresx` branch
will regenerate manifests from the post-tag branch head and should be treated
as successor-release input, for example `v1.0.1`, unless an explicit
release-history decision authorizes retagging `v1.0.0`.

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
- This snapshot does not publish or embed PostgreSQL credentials, DSNs, or raw
  SRE secret document fields.
- This module does not own domain schema, repositories, application services, or
  production DSNs.

## Evidence index

- `release/manifest/v1.0.0.json`
- `release/manifest/v1.0.0.json.sha256`
- `release/manifest/latest.json`
- `release/manifest/latest.json.sha256`
- `docs/POSTGRESX-SCORE-AUDIT-20260613.md`
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
available branch gate evidence. Publishing the post-tag `L2-T3 / 85` evidence
as a release requires a successor tag such as `v1.0.1`, an approved
manifest-contract decision for squashed source metadata, or explicit
release-history authorization to retag `v1.0.0`.
