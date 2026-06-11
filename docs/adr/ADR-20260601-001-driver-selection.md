# ADR-20260601-001: Driver Selection

Decision: use pgx/v5 for PostgreSQL connectivity, pooling, and native error information.

Consequences: `postgresx` can normalize pgx errors into foundationx errors while avoiding ORM responsibilities.
