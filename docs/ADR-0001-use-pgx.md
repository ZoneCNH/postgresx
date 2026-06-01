# ADR-0001: Use pgx as the PostgreSQL Driver

## Status

Accepted

## Context

postgresx needs a small, explicit PostgreSQL foundation for connection pooling, transactions, health checks, error normalization, metrics, tracing, migrations, and sqlc integration. The library must stay independent from application modules and must not become an ORM.

## Decision

Use `github.com/jackc/pgx/v5` and `pgxpool` as the PostgreSQL driver and pool layer for postgresx.

postgresx exposes stable wrapper contracts such as `Config`, `Client`, `DBTX`, `TxRunner`, `WithinTx`, `Migrator`, `HealthChecker`, and adapter interfaces while preserving access to `pgxpool.Pool` for advanced callers.

## Consequences

- Applications get native PostgreSQL behavior, pooling, SQLSTATE visibility, and transaction control.
- sqlc-generated repositories can depend on `postgresx.DBTX` instead of a concrete pool.
- Callers remain responsible for business schema and query ownership.
- postgresx must normalize pgx and PostgreSQL errors at its boundary.
