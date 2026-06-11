# Migration Rules

- Applications own files under their migration tree, for example
  `db/postgres/migrations` or `migrations`.
- Each schema change must have one `.up.sql` file and one matching `.down.sql` file.
- CI must prove `up`, `down`, and `up again` for live database changes.
- Dirty migrations are deployment blockers.
- `postgresx` must not contain business table definitions such as collection status, macro cursors, or regime snapshots.
- Migration logs and errors must not include credential-bearing DSNs.
