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
| local L2 gate | `L2-T3 / 85` | unit, contract, integration, chaos, benchmark smoke, and downstream compile smoke |
| release permission | `release_allowed=true` | local gate only; public tag is not moved by evidence maintenance |
| factory grade | `factory_grade_allowed=false` | requires external CI, production soak, and real downstream adoption evidence |
| release manifest | source-snapshot ancestry checked | `commit` / `tree_sha` must resolve and match the tagged source snapshot |
| workspace mode | `GOWORK=off` | required for postgresx verification evidence |
| downstream adoption | local smoke only | production adoption still requires current consumer-side dependency, compile, test, and release evidence |
| external CI | blocked | GitHub Actions cannot start while the account is locked for billing |

postgresx intentionally keeps sqlc as a caller-owned tool contract rather than
a library dependency. Domain SQL and generated query packages belong to the
consumer repository that adopts this module.
