# ADR-0003: Migrations Owned by Applications

## Status

Accepted

## Decision

Applications own schema migrations. `postgresx` provides a runner only.

## Consequences

The foundation library remains free of business schema and can be reused across x.go modules.

