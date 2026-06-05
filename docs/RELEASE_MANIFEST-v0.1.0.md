# Release Manifest — postgresx v0.1.0

## Identity

- Module: `github.com/ZoneCNH/postgresx`
- Core package: `github.com/ZoneCNH/postgresx/pkg/postgresx`
- Status: local release-preparation snapshot, not remote publication proof
- Go verification mode: `GOWORK=off`

## Included surfaces

- caller-owned `Config` construction, validation, DSN rendering, and redaction;
- `New`/`Open` client lifecycle, ping/check, pool stats, and query methods;
- sqlc-compatible `Queryer` execution boundary;
- `WithTx` and `WithTxOptions` transaction helpers;
- caller-owned migration runner over embedded or filesystem SQL sources;
- foundationx-compatible error mapping, logging, and metrics options;
- examples that keep env and secret loading outside the core package.

## Explicit non-claims

- This snapshot does not include current downstream adoption proof.
- This snapshot does not prove remote tag publication or CI status.
- This snapshot does not include checksum or contract-hash artifacts.
- This module does not own domain schema, repositories, application services, or
  production DSNs.

## Evidence index

- `release/manifest/v0.1.0.json`
- `release/manifest/latest.json`
- `docs/EVIDENCE-20260601.md`
- `docs/evidence/20260601/go-test.txt`
- `docs/evidence/20260601/go-test-race.txt`
- `docs/evidence/20260601/go-vet.txt`
- `docs/evidence/20260601/secret-scan.txt`
- `docs/evidence/20260601/no-consumer-deps.txt`
- `docs/evidence/20260601/dependencies.txt`

## Remaining release hardening

Before treating this as a published release, add fresh CI evidence, tag and
remote publication evidence, checksum files, and contract-hash files. Add
consumer adoption evidence only after a current consumer checkout proves it.
