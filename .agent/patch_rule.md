# postgresx v0.1.0 patch rule

- Module: github.com/ZoneCNH/postgresx
- Core package: github.com/ZoneCNH/postgresx/pkg/postgresx
- Boundary: explicit caller-owned config and secrets; no business schema, ORM, global database singleton, implicit env/file secret loading, or core config/observability dependencies.
- Verification: run all release gates with GOWORK=off before publishing or integrating changes.
