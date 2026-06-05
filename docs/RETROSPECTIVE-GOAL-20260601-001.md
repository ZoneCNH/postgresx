# Retrospective — POSTGRESX Goal 2026-06-01

## Outcome

postgresx converged toward a standalone PostgreSQL L2 factory module with a
small root package, caller-owned configuration, pool lifecycle helpers,
transaction helpers, migrations, error mapping, and evidence-focused release
documentation.

## Decisions that held

- Keep postgresx infrastructure-only: no domain schema, repository layer,
  application service orchestration, or implicit env/secret loading.
- Keep generated SQL and domain migrations outside this repository.
- Use a narrow `Queryer` interface so generated callers can compile against the
  adapter without depending on a concrete pool type.
- Run Go verification with `GOWORK=off` to avoid accidental workspace leakage.
- Treat downstream adoption as a separate proof requirement rather than a local
  documentation assertion.

## Corrections made during documentation convergence

- Public API documentation was narrowed to exported symbols that exist in the
  current root package.
- Error mapping documentation was aligned with the current foundationx error
  kinds emitted by `MapError`.
- Release evidence now separates local module verification from missing
  downstream adoption proof.
- Manifest evidence paths were renamed away from consumer-specific wording.
- Example live mode now requires caller-provided connection env values instead
  of silently using a default password.

## Remaining risks

- Remote release publication, checksum artifacts, and contract hashes are not
  represented by current evidence files.
- Downstream adoption remains unproven until a consumer checkout records a
  fresh dependency pin, compile/test evidence, import-boundary evidence, and
  release-manifest linkage.
- README-level API wording may still need a separate owner pass if it is treated
  as part of the release surface.

## Follow-up rule

Future modifiers should update docs, contracts, examples, and evidence in the
same change whenever exported API names or release claims change. Do not add a
consumer adoption claim without current consumer-side evidence.
