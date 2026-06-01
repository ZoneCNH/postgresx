# ADR-0004: Use Explicit Transaction Boundaries

## Status

Accepted

## Context

Repository methods need a stable way to run one or more SQL statements atomically without depending directly on `pgxpool.Pool`. Callers also need retry behavior for serialization failures and deadlocks.

## Decision

postgresx exposes explicit transaction helpers through `Client.WithinTx`, `Client.WithinTxRetry`, and the `TxRunner` interface.

Repositories accept `postgresx.DBTX` for query execution and `postgresx.TxRunner` when they need to create a transaction boundary. They should not begin transactions by reaching directly into `pgxpool.Pool`.

## Consequences

- Transaction behavior is visible in constructor signatures.
- Tests can fake `DBTX` and `TxRunner` without a live database.
- Retry policy stays centralized in postgresx.
- Consumer repositories can remain sqlc-compatible while still using safe transaction boundaries.
