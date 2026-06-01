# ADR-20260601-002: Transaction Semantics

Decision: expose explicit transaction callbacks and minimal query interfaces.

Consequences: callers control unit-of-work boundaries, and the library does not create hidden global transaction state.
