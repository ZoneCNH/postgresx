# ADR-0001: Use pgx as the PostgreSQL driver

## Status

Accepted

## Context

postgresx needs a production-grade PostgreSQL driver with connection pooling,
transaction control, error details, and low-level escape hatches for advanced
callers while keeping this module infrastructure-only.

## Decision

Use `github.com/jackc/pgx/v5` and `pgxpool` internally. Expose a stable wrapper
contract around `Config`, `Client`, `Queryer`, `Tx`, transaction helpers,
migration helpers, health checks, and observability adapters rather than
requiring callers to own pgx construction directly.

## Consequences

- postgresx can map pg errors into foundationx-compatible error kinds.
- sqlc-generated callers can depend on `postgresx.Queryer` instead of a concrete
  pool type.
- Advanced pgx behavior remains available through caller-owned SQL and options,
  without adding ORM or domain abstractions to this module.
