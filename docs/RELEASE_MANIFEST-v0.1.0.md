# Release Manifest — postgresx v0.1.0

Status: local release candidate prepared for `github.com/ZoneCNH/postgresx`.

## Identity

- Module: `github.com/ZoneCNH/postgresx`
- Core package: `github.com/ZoneCNH/postgresx/pkg/postgresx`
- Version: `v0.1.0`
- Date: 2026-06-01
- Machine-readable manifest: `release/manifest/v0.1.0.json`

## Release Contents

- pgx-backed client lifecycle and explicit pool ownership.
- sqlc-compatible `DBTX` execution boundary.
- Explicit transaction helper and retry policy.
- Migration runner for caller-owned migration trees.
- Error normalization, health checks, stats, metrics/tracing interfaces, and DSN masking.
- Testkit helpers and examples outside the core package.

## Exclusions

- No ORM or business schema ownership.
- No implicit env/file secret loading.
- No package-global DB/singleton.
- No core dependency on application wiring packages such as `configx` or `observex`.
- No dependency on application modules.

## Required Gates

```sh
GOWORK=off make ci
GOWORK=off make ci-extended
GOWORK=off make release-evidence-check
GOWORK=off make release-preflight
```

`make release-check` additionally sets required live integration semantics through
`scripts/ci/release_check.sh`.

## Evidence

Primary evidence is recorded in `docs/EVIDENCE-20260601.md` and raw command logs
under `docs/evidence/20260601/`.

## Publication Notes

This worker prepared local release artifacts and verification gates. It did not
create or push a remote tag.
