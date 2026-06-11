# Version Matrix

| Component | Version / Contract | Notes |
| --- | --- | --- |
| Go | 1.24.4 | declared by `go.mod` |
| pgx | v5.7.5 | direct PostgreSQL driver and pool dependency |
| foundationx | v0.1.0 | config redaction, error, logging, metrics contracts |
| sqlc | tool contract only | callers can target the exported `Queryer` interface |
| workspace mode | `GOWORK=off` | required for postgresx verification evidence |
| downstream adoption | not proven | requires current consumer-side dependency, compile, test, and release evidence |

postgresx intentionally keeps sqlc as a caller-owned tool contract rather than
a library dependency. Domain SQL and generated query packages belong to the
consumer repository that adopts this module.
