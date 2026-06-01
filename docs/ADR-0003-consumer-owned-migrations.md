# ADR-0003: Consumers Own Business Migrations

## Status

Accepted

## Context

postgresx must support migration execution and dirty-state checks without owning consumer schemas. x.go and future services need to evolve their own business tables independently.

## Decision

postgresx provides migration runner primitives and CI-friendly checks. Consumer repositories own migration files, schema design, versioning policy, and rollout order.

postgresx examples may demonstrate usage, but they must not become the source of truth for x.go business schema.

## Consequences

- The postgresx module remains independent from x.go.
- x.go can add tables such as `market_data.collection_status` in its own migration tree.
- Dirty migration handling is tested in postgresx, while production schema responsibility stays with applications.
