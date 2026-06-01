# Application Integration Boundary

`postgresx` is application-independent. Application repositories may import
`github.com/ZoneCNH/postgresx/pkg/postgresx`, but this module must not import
application code.

Expected direction:

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

The release gate currently checks the concrete forbidden application dependency
pattern configured in `scripts/check_boundary.sh` and
`scripts/ci/release_check.sh`.
