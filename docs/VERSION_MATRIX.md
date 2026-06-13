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
| published `v1.0.0` tag | `L2-T2 / 75` | immutable restored snapshot; `release_allowed=false` |
| local L2 gate | `L2-T3 / 85` | current branch only; unit, contract, integration, chaos, benchmark smoke, and downstream compile smoke |
| release permission | `release_allowed=true` | local gate only; public tag is not moved by evidence maintenance |
| factory grade | `factory_grade_allowed=false` | requires external CI, production soak, and real downstream adoption evidence |
| release evidence check | blocked for `v1.0.0` | source commit `9eaf770` resolves but is outside current `HEAD` and tag ancestry |
| release blocker diagnostic | `release-blockers` | non-mutating ancestry report for manifest commit, current `HEAD`, and tag commit |
| release manifest | dual-scope evidence | tag manifest preserves snapshot metadata; current-branch evidence must ship through a successor release or approved manifest-contract change |
| workspace mode | `GOWORK=off` | required for postgresx verification evidence |
| downstream adoption | local smoke only | production adoption still requires current consumer-side dependency, compile, test, and release evidence |
| external CI | blocked | GitHub Actions cannot start while the account is locked for billing |

postgresx intentionally keeps sqlc as a caller-owned tool contract rather than
a library dependency. Domain SQL and generated query packages belong to the
consumer repository that adopts this module.
