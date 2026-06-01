# Retrospective — GOAL-20260601-001

## Outcome

The local postgresx release surface was aligned with the v0.1.0 foundation goal:
module identity is `github.com/ZoneCNH/postgresx`, the public core package lives
under `pkg/postgresx`, and validation/release gates run with `GOWORK=off`.

## What Changed

- Moved the core package to `pkg/postgresx` while keeping examples, contracts,
  testkit, and internal helpers outside the core package.
- Normalized docs, examples, tests, contracts, and scripts to the ZoneCNH module
  identity.
- Added release/evidence Makefile gates and a machine-readable release manifest.
- Strengthened boundary checks for forbidden application dependencies,
  business-domain terms, and core `configx`/`observex` drift.
- Documented explicit secret/config/observability ownership boundaries.

## Decisions

- Keep migrations caller-owned; postgresx only runs migration sources supplied by
  the caller.
- Keep sqlc support as an interface boundary (`DBTX`) rather than generated code
  or business repositories inside the library.
- Keep metrics/tracing adapter-style and opt-in so the core stays independent of
  application observability packages.

## Remaining Work

- The prior `/home/x.go/pkg/adapter/db/postgres/postgres.go` candidate opened
  its existing pool configuration through `postgresx.OpenPool`.
- x.go remained the owner of business SQL, migrations, and generated query code.
- The prior `collection_status` candidate had an x.go-owned repository,
  migration, and test slice that depended on `postgresx.DBTX` and
  `postgresx.TxRunner`.
- The prior `collection_status` SQL was generated through x.go's local sqlc
  v1.31.1 configuration rather than embedded by postgresx.
- The prior x.go import-boundary gate included `B-POSTGRESX-001`, which forbade
  direct `pgxpool` imports under `internal/market_data/server/state/**` while
  preserving the existing adapter-owned pool boundary.
- Targeted x.go tests passed in the prior candidate for
  `./internal/market_data/server/state ./pkg/adapter/db/postgres ./internal/bootstrap`.
- x.go import-boundary tests and the full import-boundary gate passed in the
  prior candidate.
- `docs/EVIDENCE-20260601.md` records the command evidence and remaining
  release gaps.

## Remaining Evidence Needed

- Release commit/staging remains open: the current `/home/postgresx` `HEAD`
  tracks only `.gitignore`, `LICENSE`, and `README.md`, while the implementation
  candidate is still local working-tree state.
- Git tag `v0.1.0` and remote publication verification require explicit
  release authorization after that release commit/staging decision.

## Verified Evidence

- `docs/EVIDENCE-20260601.md` records a live PostgreSQL migration up/down/up
  gate with `dirty=false` checks.
- `docs/EVIDENCE-20260601.md` records the live metrics/tracing adapter
  integration check.
- `docs/EVIDENCE-20260601.md` records the historical live x.go
  `collection_status` PostgreSQL integration run.
- `docs/EVIDENCE-20260601.md` records the historical x.go `B-POSTGRESX-001`
  import-boundary test and repository-wide import-boundary gate run.
- The current `/home/x.go/go.mod` does not record
  `github.com/ZoneCNH/postgresx v0.1.0`.

## Output Patch

- Prompt patch: keep local implementation evidence separate from release
  publication authority; do not mark `v0.1.0` complete until tag and remote
  publication are verified.
- Release patch: do not cut a tag from the current initial `HEAD`; the tag must
  point at an authorized release commit that contains the implementation,
  docs, tests, and x.go integration evidence.
- Harness patch: release evidence must show whether live PostgreSQL migration
  gates ran against a real `POSTGRES_TEST_DSN` or were skipped by CI defaults.
- Rule patch: business-owned x.go database modules should depend on
  `postgresx.DBTX` and `postgresx.TxRunner`; direct `pgxpool` ownership stays
  in platform adapter boundaries or explicit escape hatches. Enforce this for
  `collection_status` through x.go's import-boundary gate rather than a broad
  grep that would fail existing adapter/infra pool ownership.
