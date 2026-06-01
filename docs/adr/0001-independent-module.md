# ADR-0001: Independent Module

## Status

Accepted

## Decision

`postgresx` is an independent Go module at `github.com/ZoneCNH/postgresx/pkg/postgresx`.

## Consequences

Application modules can import the library, but `postgresx` must never import application modules.

