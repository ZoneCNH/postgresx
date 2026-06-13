# Postgresx Metrics Contract

`postgresx` exposes metrics through a small application-owned adapter. The
package does not choose a metrics backend, exporter, registry, or tracing SDK.

## Metrics

Applications pass a `Metrics` implementation through `WithMetrics`.

- `IncCounter(name string, labels map[string]string)` records monotonic counts.
- `ObserveHistogram(name string, value float64, labels map[string]string)` records durations.
- `SetGauge(name string, value float64, labels map[string]string)` records pool state.

## Metric Names

- `postgresx.query.total` is incremented by `Ping`, `Exec`, and `Query`.
- `postgresx.query.duration_seconds` records `Ping`, `Exec`, and `Query` duration.
- `postgresx.tx.total` is incremented by `WithTx` and `WithTxOptions`.
- `postgresx.tx.duration_seconds` records transaction duration.
- `postgresx.health.total` is incremented by `Check` and `HealthCheck`.
- `postgresx.health.latency_seconds` records health-check duration.
- `postgresx.pool.connections` records `Stats` gauges for `total`, `idle`, `acquired`, `constructing`, and `max`.

Adapters must not log or export raw credentials.
