# configx / observex Boundary

`postgresx` is a foundation library. The core package
`github.com/ZoneCNH/postgresx/pkg/postgresx` must not import application wiring
packages such as `configx` or `observex`.

Allowed boundaries:

- callers pass explicit `postgresx.Config`, `pgxpool.Config`, or
  `postgresx.MigrationConfig` values;
- metrics/tracing are represented by small interfaces and no-op defaults;
- callers decide where configuration and observability adapters come from;
- no env/file secret loading occurs inside the core package.

Release enforcement:

```sh
GOWORK=off make boundary
GOWORK=off go list -deps ./pkg/postgresx
```

The boundary check fails if core files mention `configx` or `observex`, if the
module depends on an application module, or if business-domain schema terms enter
library code.
