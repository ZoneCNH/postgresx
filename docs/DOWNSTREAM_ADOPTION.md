# Downstream Adoption Boundary

postgresx is an infrastructure adapter. It provides PostgreSQL configuration,
pool lifecycle, health checks, transactions, migrations, error mapping, and a
narrow query interface. Downstream services own schema, generated query code,
repositories, application transactions, and release wiring.

## Current status

The repository includes a local downstream smoke check that creates a temporary
consumer module, imports `github.com/ZoneCNH/postgresx`, compiles against the
public `Queryer` boundary, and verifies caller-owned configuration/redaction
behavior. This supports the local `L2-T3 / 85` gate.

That smoke check is not production adoption proof. Documentation and release
notes must avoid claiming that a downstream service has adopted the module until
a fresh consumer checkout records dependency, compile, test, and release
evidence.

## Dependency direction

The allowed direction is one-way:

1. A downstream service may require `github.com/ZoneCNH/postgresx`.
2. postgresx must not require downstream service modules or peer L2 adapters.
3. Domain SQL, migrations, and generated code remain outside postgresx.

## Caller-owned responsibilities

Downstream callers are responsible for:

- loading env and secret material before constructing `postgresx.Config`;
- deciding schema names, migration contents, and generated query packages;
- deciding transaction scope at the application boundary;
- wiring observability adapters through the exported logger and metrics options;
- storing consumer-specific release evidence with the consumer release.

## Proof required before an adoption claim

A valid downstream adoption claim needs fresh evidence from the consumer
checkout:

- dependency pin or workspace replacement for this module version;
- compile-only check for packages that use the `Queryer` boundary;
- targeted tests for repository and transaction call sites;
- import-boundary evidence showing postgresx has no reverse dependency;
- release manifest entry linking the consumer release to the postgresx version.
