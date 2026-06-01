# Version Matrix

| Component | Version / Path | Notes |
| --- | --- | --- |
| postgresx module | `github.com/ZoneCNH/postgresx` | v0.1.0 release candidate |
| postgresx core package | `github.com/ZoneCNH/postgresx/pkg/postgresx` | Public API surface |
| Go | `go.mod` `go 1.26.3` | All validation uses `GOWORK=off` |
| pgx | `github.com/jackc/pgx/v5 v5.9.2` | PostgreSQL driver/pool foundation |
| foundationx | `github.com/ZoneCNH/foundationx v0.1.1` | Foundation contract dependency |
| golang-migrate | `github.com/golang-migrate/migrate/v4` | Caller-owned migrations |

`postgresx` intentionally keeps sqlc as a tool contract rather than a library dependency. The public `DBTX` interface is compatible with pgx/sqlc generated code.

The current `/home/x.go` checkout is not the source of a postgresx requirement:
`GOWORK=off go list -m github.com/ZoneCNH/postgresx/pkg/postgresx` in `/home/x.go`
currently reports that the module is not a known dependency.
