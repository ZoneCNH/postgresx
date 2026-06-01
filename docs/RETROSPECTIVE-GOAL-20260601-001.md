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

- Remote tag creation and published release notes remain outside this local task.
- Live PostgreSQL validation requires a DSN or Docker-capable runner; local gates
  document skip behavior when neither is available.
