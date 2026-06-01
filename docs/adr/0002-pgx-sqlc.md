# ADR-0002: pgx and sqlc

## Status

Accepted

## Decision

Use pgx as the PostgreSQL driver and expose a sqlc-compatible `DBTX` interface.

## Consequences

Application code remains SQL-first and can pass either a pool or transaction to generated queries.

