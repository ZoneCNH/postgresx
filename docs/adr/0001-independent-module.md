# ADR-0001: Independent Module

## Status

Accepted

## Decision

`postgresx` is an independent Go module at `github.com/bytechainx/postgresx`.

## Consequences

`x.go` can import the library, but `postgresx` must never import `x.go`.

