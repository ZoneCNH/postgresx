# Consumer Integration Guide

postgresx is adopted by a consumer service as an infrastructure dependency, not
as a domain layer.

## Current status

This repository includes a local consumer smoke check for import, compile,
configuration redaction, and `postgresx.Queryer` compatibility. Treat the steps
below as integration guidance until a real consumer checkout records dependency,
compile, test, import-boundary, and release evidence.

## Recommended flow

1. Add a dependency on the intended postgresx version.
2. Load env and secret material in consumer-owned startup code.
3. Construct `postgresx.Config` explicitly and call `Validate` before opening a
   client.
4. Pass `client.Queryer()` or a transaction `Tx` into generated repositories
   that target the `postgresx.Queryer` interface.
5. Store domain SQL and migrations in the consumer repository.
6. Record consumer-side compile and test evidence with the consumer release.

## Boundary

postgresx must not import consumer modules, peer L2 adapters, domain schema, or
application services. Consumer code owns all production wiring and secret
material.
