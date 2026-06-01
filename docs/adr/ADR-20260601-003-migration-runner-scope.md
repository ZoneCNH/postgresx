# ADR-20260601-003: Migration Runner Scope

Decision: include a small migration runner for ordered SQL steps and dirty-state checks.

Consequences: applications retain ownership of migration files and schemas; postgresx owns execution mechanics only.
