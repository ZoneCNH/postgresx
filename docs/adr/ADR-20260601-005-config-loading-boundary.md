# ADR-20260601-005: Config Loading Boundary

Decision: `pkg/postgresx` accepts explicit configuration and never loads secrets from env or files.

Consequences: applications may use their own config stack, but core remains deterministic, testable, and reusable.
