# postgresx API

The supported import path is `github.com/ZoneCNH/postgresx/pkg/postgresx`.

Stable v0.1.0 surfaces:

- `Config`, `DefaultConfig`, `Validate`, `DSN`, and `RedactedDSN` for explicit
  caller-owned configuration.
- `Open`, `New`, `OpenPool`, `Client`, `Close`, `Ping`, `Check`, `Stats`,
  `Exec`, `Query`, `QueryRow`, `Queryer`, `DB`, and `RawPool`.
- `WithTx`, `WithTxOptions`, `WithinTx`, `Tx`, and `DBTX` for explicit
  transaction and sqlc boundaries.
- `Migration`, `MigrationSource`, `MigrationRunner`, `Up`, `Applied`, and
  rollback-on-failure behavior for caller-owned migration SQL.
- `MapError`, `IsRetryable`, `MaskDSN`, `Metrics`, `Tracer`, and options for
  foundationx-compatible errors and application-owned observability adapters.

The core package must not read env, read secret files, own application schema,
choose an ORM, expose a global database handle, or import application modules.
