# ADR-0002: Support sqlc-style callers, not ORM ownership

## Status

Accepted

## Context

The v1.0 scope explicitly excludes ORM features. Application services and other consumers own domain repositories, SQL, and migrations, while postgresx provides the foundation layer.

## Decision

postgresx exposes a sqlc-compatible `DBTX` interface and transaction runner instead of defining models, query builders, or repository abstractions.

Consumer code should generate or write domain-specific SQL repositories against `postgresx.DBTX` and use `postgresx.TxRunner` for transaction boundaries.

## Consequences

- postgresx stays small and reusable.
- Domain SQL remains close to the consuming application.
- Transaction ownership is explicit and testable.
- Adding ORM behavior is out of scope for the v1 foundation.
