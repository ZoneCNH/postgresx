# ADR-0004: Normalized Errors and Secret Masking

## Status

Accepted

## Decision

Normalize PostgreSQL failures into stable `ErrorKind` values and mask credential-bearing text before exposing public error strings.

## Consequences

Callers can make retry and not-found decisions without depending on driver-specific errors, while logs remain safer by default.

