# ADR-0002: Prefer sqlc-compatible callers over an ORM

## Status

Accepted

## Context

postgresx must provide PostgreSQL infrastructure without owning domain models,
repositories, generated queries, or schema-specific migrations.

## Decision

postgresx exposes a sqlc-compatible `Queryer` interface and transaction helpers
instead of defining models, query builders, or repository abstractions.

Consumer code should generate or write domain-specific SQL repositories against
`postgresx.Queryer` and use `postgresx.WithTx` for transaction boundaries.

## Consequences

- The module stays schema-agnostic.
- Callers keep full control over SQL, migrations, and generated packages.
- Documentation and examples must not imply that postgresx owns domain query
  code.
