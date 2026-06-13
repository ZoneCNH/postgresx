# Release Manifest — postgresx v1.0.0

## Identity

- Module: `github.com/ZoneCNH/postgresx`
- Core package: `github.com/ZoneCNH/postgresx/pkg/postgresx`
- Status: published `v1.0.0` evidence snapshot
- Go verification mode: `GOWORK=off`

## Publication evidence

- GitHub release: `https://github.com/ZoneCNH/postgresx/releases/tag/v1.0.0`
- Remote branch: `refs/heads/postgresx` resolves to `310a249`.
- Remote tag: `refs/tags/v1.0.0^{}` resolves to `310a249`.
- Release metadata `targetCommitish` remains `main`; the tag object and resolved
  tag commit are the authoritative release identity.

## Manifest semantics

`release/manifest/v1.0.0.json` and `release/manifest/latest.json` record the
source snapshot used to generate the release evidence. Their `commit` and
`tree_sha` fields must resolve in Git, the tree must match the resolved commit,
and the resolved commit must be an ancestor of both the evidence-carrying HEAD
and the `v1.0.0` tag when the tag is present. They are not required to equal a
later evidence-maintenance commit.

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
- This snapshot does not prove current CI status or production soak.
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

## Remaining release hardening

Before treating this as a production-adopted release, add fresh CI evidence,
production soak evidence, and consumer adoption evidence from a current consumer
checkout. Do not rewrite or retag `v1.0.0` without an explicit release-history
approval.
