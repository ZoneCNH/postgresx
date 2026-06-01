# x.go Integration

## Current Checkout Status

As of the 2026-06-01 documentation alignment audit, the mounted `/home/x.go`
checkout is not integrated with `postgresx`.

Current evidence:

- `/home/x.go/go.mod` does not require `github.com/bytechainx/postgresx`.
- With `/home/go.work` active, `go list -m github.com/bytechainx/postgresx`
  returns the sibling workspace module.
- With `GOWORK=off`, `go list -m github.com/bytechainx/postgresx` in
  `/home/x.go` fails because the module is not a known dependency.
- `/home/x.go/pkg/adapter/db/postgres/postgres.go` still imports and owns
  `github.com/jackc/pgx/v5/pgxpool` directly.
- No `sqlc.yaml` or `sqlc.yml` file is present under `/home/x.go` at max depth
  3, and the collection-status files listed in earlier evidence are absent.

## Intended Dependency Direction

When x.go adopts this module, the dependency direction must be:

```text
x.go -> postgresx -> PostgreSQL
```

`postgresx` must never import `x.go`.

## Adoption Contract

The x.go PostgreSQL adapter can keep x.go-specific DSN parsing,
`pgxpool.Config` mutation, logging, and port compatibility. The planned
adoption path is to open the prepared pool through `postgresx.OpenPool(ctx,
poolCfg)`, store the returned `postgresx.Client`, and expose the existing pool
boundary to x.go callers through `client.RawPool()` only where that boundary is
already intentional.

Business-owned x.go packages should depend on `postgresx.DBTX` and
`postgresx.TxRunner` instead of owning `pgxpool` directly.

## collection_status Plan

If x.go reintroduces the collection-status persistence slice, it should remain
owned by x.go and include files equivalent to:

- `internal/market_data/server/state/collection_status.go`
- `internal/market_data/server/state/sqlc/collection_status.sql`
- `internal/market_data/server/state/generated/collection_status.sql.go`
- `migrations/000010_collection_status.up.sql`
- `migrations/000010_collection_status.down.sql`
- `internal/market_data/server/state/collection_status_test.go`
- `sqlc.yaml`

## Verified Commands

From `/home/postgresx`:

```sh
go test ./...
```

Current `/home/x.go` recheck:

```sh
go list -m github.com/bytechainx/postgresx
GOWORK=off go list -m github.com/bytechainx/postgresx
rg -n "postgresx|collection_status|sqlc" go.mod internal/market_data/server/state pkg/adapter/db/postgres migrations
```

With `/home/go.work` active, `go list -m github.com/bytechainx/postgresx` returns
the workspace module:

```text
github.com/bytechainx/postgresx
```

With `GOWORK=off`, x.go's own module currently reports:

```text
go: module github.com/bytechainx/postgresx: not a known dependency
```

## Configuration

Read DSNs from the owning application environment or secret manager. Do not hardcode secret paths or credential-bearing DSNs in source, tests, logs, or docs.
