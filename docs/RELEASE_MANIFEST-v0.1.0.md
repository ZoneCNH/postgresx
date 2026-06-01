# Release Manifest v0.1.0

Status: local `postgresx` implementation prepared; current `/home/x.go`
checkout is not integrated with `postgresx`; git tag and remote publication are
not yet verified.

## Included

- Go module `github.com/ZoneCNH/postgresx/pkg/postgresx`
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

## Verified Locally

- `/home/postgresx`: `go test ./...`
- `/home/postgresx`: `go vet ./...`
- `/home/postgresx`: `go test -race ./...`
- `/home/postgresx`: `./scripts/ci/secret_scan.sh`
- `/home/postgresx`: `POSTGRES_TEST_DSN=<local temporary PostgreSQL DSN> ./scripts/ci/migration_up_down_up.sh`
- `/home/postgresx`: `POSTGRES_TEST_DSN=<local temporary PostgreSQL DSN> go test -run 'Test(OpenPingStatsIntegration|MetricsTracingIntegration|WithinTxCommitRollbackIntegration|WithinTxRetryIntegration|MigratorUpDownUpIntegration)' -count=1 ./...`
- `/home/postgresx`: `make ci`
- `/home/postgresx`: GitHub Actions workflow provisions PostgreSQL 17 and
  exports `POSTGRES_TEST_DSN` before package tests, race tests, and migration
  tests.

## Current x.go Recheck

- `/home/x.go/go.mod` does not require `github.com/ZoneCNH/postgresx/pkg/postgresx`.
- `/home/x.go/pkg/adapter/db/postgres/postgres.go` still imports
  `github.com/jackc/pgx/v5/pgxpool` directly.
- `GOWORK=off go list -m github.com/ZoneCNH/postgresx/pkg/postgresx` in `/home/x.go`
  reports that `github.com/ZoneCNH/postgresx/pkg/postgresx` is not a known dependency.
- No `sqlc.yaml` or `sqlc.yml` file is present under `/home/x.go` at max depth
  3, and the previously documented `collection_status` files are absent from
  the current checkout.

## Evidence Archive

- `docs/EVIDENCE-20260601.md` records the local command evidence,
  dependency pins, live PostgreSQL migration checks, historical x.go
  collection-status evidence, the current x.go recheck, and remaining release
  gaps.
- `docs/evidence/20260601/` contains raw command logs for gofmt, vet,
  package tests, race tests, secret scan, live PostgreSQL migration, focused
  live integration tests, and historical `GOWORK=off` module/dependency checks.
- `docs/evidence/20260601/migration-up-down-up.txt` records
  `version=1 dirty=false` after the first up and again after the up/down/up
  cycle.

## Required Before Publication

- git tag `v0.1.0`
- remote tag and release publication verification
