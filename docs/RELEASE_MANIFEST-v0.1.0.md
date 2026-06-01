# Release Manifest — postgresx v0.1.0

Status: local release candidate prepared for `github.com/ZoneCNH/postgresx`.

## Identity

- Go module `github.com/ZoneCNH/postgresx`
- pgxpool client lifecycle
- sqlc `DBTX`
- transaction and retry helpers
- migration runner
- normalized errors with explicit SQLSTATE mappings
- health and pool stats
- metrics/tracing adapter interfaces
- public `Mask` and `MaskDSN` secret masking
- testkit and examples
- CI scripts and GitHub Actions workflow
- integration documentation describing the intended x.go adoption boundary and
  current checkout status

## Release Contents

- pgx-backed client lifecycle and explicit pool ownership.
- sqlc-compatible `DBTX` execution boundary.
- Explicit transaction helper and retry policy.
- Migration runner for caller-owned migration trees.
- Error normalization, health checks, stats, metrics/tracing interfaces, and DSN masking.
- Testkit helpers and examples outside the core package.

## Exclusions

- `/home/x.go/go.mod` does not require `github.com/ZoneCNH/postgresx/pkg/postgresx`.
- `/home/x.go/pkg/adapter/db/postgres/postgres.go` still imports
  `github.com/jackc/pgx/v5/pgxpool` directly.
- `GOWORK=off go list -m github.com/ZoneCNH/postgresx/pkg/postgresx` in `/home/x.go`
  reports that `github.com/ZoneCNH/postgresx/pkg/postgresx` is not a known dependency.
- No `sqlc.yaml` or `sqlc.yml` file is present under `/home/x.go` at max depth
  3, and the previously documented `collection_status` files are absent from
  the current checkout.

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
