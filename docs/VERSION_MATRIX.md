# Version Matrix

| Component | Version / Contract | Notes |
| --- | --- | --- |
| postgresx | v1.0.0 | release target and public API contract baseline |
| Go | 1.25.0 | declared by `go.mod` |
| pgx | v5.9.2 | direct PostgreSQL driver and pool dependency |
| foundationx | v0.1.1 | config redaction, error, logging, metrics contracts |
| sqlc | tool contract only | callers can target the exported `Queryer` interface |
| metrics | dotted `postgresx.*` names | locked by `contracts/metrics.md` and package tests |
| integration | real PostgreSQL required for release gate | inject `POSTGRESX_INTEGRATION_DSN` / `POSTGRES_TEST_DSN` through environment only |
| workspace mode | `GOWORK=off` | required for postgresx verification evidence |
| downstream adoption | not proven | requires current consumer-side dependency, compile, test, and release evidence |

postgresx intentionally keeps sqlc as a caller-owned tool contract rather than
a library dependency. Domain SQL and generated query packages belong to the
consumer repository that adopts this module.
