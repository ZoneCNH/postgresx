# ADR-0003: Consumer-owned migrations

## Status

Accepted

## Context

postgresx must support migration execution and dirty-state checks without owning consumer schemas. Applications and future services need to evolve their own business tables independently.

## Decision

postgresx provides migration runner primitives and CI-friendly checks. Consumer repositories own migration files, schema design, versioning policy, and rollout order.

postgresx examples may demonstrate usage, but they must not become the source of truth for application business schema.

## Consequences

- The postgresx module remains independent from application modules.
- Applications can add business tables in their own migration trees.
- Dirty migration handling is tested in postgresx, while production schema responsibility stays with applications.
