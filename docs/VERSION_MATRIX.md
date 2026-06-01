# Version Matrix

Generated: 2026-06-01.

| Component | Version | Source |
| --- | --- | --- |
| Go | 1.26.3 | Local toolchain, `/home/go.work`, and `/home/postgresx/go.mod` |
| `github.com/jackc/pgx/v5` | v5.9.2 | `/home/postgresx/go.mod` |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | `/home/postgresx/go.mod` |
| `github.com/sqlc-dev/sqlc` | v1.31.1 | Historical x.go integration evidence; not a postgresx module dependency |

`postgresx` intentionally keeps sqlc as a tool contract rather than a library dependency. The public `DBTX` interface is compatible with pgx/sqlc generated code.

The current `/home/x.go` checkout is not the source of a postgresx requirement:
`GOWORK=off go list -m github.com/ZoneCNH/postgresx/pkg/postgresx` in `/home/x.go`
currently reports that the module is not a known dependency.
