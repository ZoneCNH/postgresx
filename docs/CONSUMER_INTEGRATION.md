# Consumer Integration Boundary

Consumer applications should depend on `github.com/ZoneCNH/postgresx/pkg/postgresx` and construct `postgresx.Config` from their own configuration boundary.

Recommended direction:

```text
application -> postgresx -> PostgreSQL
```

The dependency must not point back into application code. Applications own domain schemas, repositories, sqlc packages, migrations, and observability wiring. `postgresx` owns generic PostgreSQL runtime behavior only.
