# Retrospective: postgresx Goal

Current audit note: the 2026-06-01 documentation alignment pass found that the
mounted `/home/x.go` checkout does not currently require `postgresx`, still owns
`pgxpool` directly in `pkg/adapter/db/postgres/postgres.go`, and does not contain
the archived `collection_status` sqlc slice. The x.go integration outcome below
is historical evidence from a prior local candidate, not current checkout state.

## Decisions

- Kept `postgresx` independent from x.go.
- Used pgx/sqlc contracts directly instead of introducing an ORM abstraction.
- Kept migrations owned by applications and implemented only the runner.
- Added explicit secret masking and a CI secret scan.
- Exercised pluggable metrics and tracing with a live PostgreSQL integration
  test instead of relying only on API presence.
- Exposed `Mask` and `MaskDSN` from the root package so callers can sanitize
  values without importing internal packages.
- Aligned normalized error kind literals with the goal contract for unique
  violation, foreign-key violation, serialization failure, deadlock, and
  context cancellation.
- Added `OpenPool` so x.go can keep its existing pgxpool configuration and
  still adopt the postgresx lifecycle and health contract.

## Historical Integration Outcome

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
  `github.com/bytechainx/postgresx v0.1.0`.

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
