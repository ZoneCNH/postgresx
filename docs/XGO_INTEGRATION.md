# Application Integration Boundary

`postgresx` is application-independent. Application repositories may import
`github.com/ZoneCNH/postgresx/pkg/postgresx`, but this module must not import
application code.

As of the 2026-06-01 documentation alignment audit, the mounted `/home/x.go`
checkout is not integrated with `postgresx`.

Current evidence:

- `/home/x.go/go.mod` does not require `github.com/ZoneCNH/postgresx`.
- With `/home/go.work` active, `go list -m github.com/ZoneCNH/postgresx`
  returns the sibling workspace module.
- With `GOWORK=off`, `go list -m github.com/ZoneCNH/postgresx` in
  `/home/x.go` fails because the module is not a known dependency.
- `/home/x.go/pkg/adapter/db/postgres/postgres.go` still imports and owns
  `github.com/jackc/pgx/v5/pgxpool` directly.
- No `sqlc.yaml` or `sqlc.yml` file is present under `/home/x.go` at max depth
  3, and the collection-status files listed in earlier evidence are absent.

## Intended Dependency Direction

When x.go adopts this module, the dependency direction must be:

```text
application -> postgresx -> PostgreSQL
```

## Adoption Pattern

- Application configuration loads DSNs/secrets outside postgresx.
- Application code passes explicit `postgresx.Config`, `pgxpool.Config`, or
  `postgresx.MigrationConfig` values.
- Application repositories own SQL, migrations, generated sqlc packages, and
  business repositories.
- postgresx exposes `DBTX`, `Client.DB()`, `Client.RawPool()`, and transaction
  helpers only as infrastructure boundaries.

## Verification

```sh
GOWORK=off make boundary
GOWORK=off go list -deps ./... | rg 'github.com/.*/application-module'
```

Current `/home/x.go` recheck:

```sh
go list -m github.com/ZoneCNH/postgresx
GOWORK=off go list -m github.com/ZoneCNH/postgresx
rg -n "postgresx|collection_status|sqlc" go.mod internal/market_data/server/state pkg/adapter/db/postgres migrations
```

With `/home/go.work` active, `go list -m github.com/ZoneCNH/postgresx` returns
the workspace module:

```text
github.com/ZoneCNH/postgresx/pkg/postgresx
```

With `GOWORK=off`, x.go's own module currently reports:

```text
go: module github.com/ZoneCNH/postgresx/pkg/postgresx: not a known dependency
```

## Configuration

Read DSNs from the owning application environment or secret manager. Do not hardcode secret paths or credential-bearing DSNs in source, tests, logs, or docs.
