# Public API Contract

The supported import path is `github.com/ZoneCNH/postgresx/pkg/postgresx`.

Stable v0.1.0 surfaces:

- `Config`, `SanitizedConfig`, `DefaultConfig`, `Validate`, `Sanitize`, `DSN`,
  and `RedactedDSN` for explicit caller-owned configuration.
- `New`, `Open`, `Client`, `Close`, `Ping`, `Check`, `Stats`, `Exec`, `Query`,
  `QueryRow`, and `Queryer` for pool-backed access without exposing a global
  database handle.
- `CommandTag`, `Row`, `Rows`, and `Queryer` as the narrow sqlc-compatible
  execution boundary.
- `Tx`, `TxFunc`, `TxOptions`, `WithTx`, and `WithTxOptions` for explicit
  transaction ownership and retry behavior.
- `Migration`, `MigrationSource`, `AppliedMigration`, `MigrationRunner`,
  `NewMigrationRunner`, `Up`, and `Applied` for caller-owned migration SQL.
- `MapError`, `IsRetryable`, `Logger`, `Metrics`, `Field`, `Option`,
  `WithLogger`, `WithMetrics`, and `WithClock` for foundationx-compatible
  errors and caller-owned observability adapters.

The core package must not read env, read secret files, own application schema,
choose an ORM, expose a global database handle, or import application modules.
