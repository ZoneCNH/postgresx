package postgresx_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
	"github.com/ZoneCNH/postgresx/testkit"
)

type migrationSource []postgresx.Migration

func (s migrationSource) List(ctx context.Context) ([]postgresx.Migration, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	migrations := make([]postgresx.Migration, len(s))
	copy(migrations, s)
	return migrations, nil
}

func TestOpenCloseIntegration(t *testing.T) {
	ctx := t.Context()
	dsn := currentIntegrationDSN(t)
	cfg, err := testkit.ConfigFromDSN(dsn, "postgresx-open-close")
	if err != nil {
		t.Fatalf("ConfigFromDSN() error = %v, want nil", err)
	}

	client, err := postgresx.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open() error = %v, want nil", err)
	}
	if err := client.Ping(ctx); err != nil {
		t.Fatalf("Ping() error = %v, want nil", err)
	}
	if err := client.Close(ctx); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
	if err := client.Close(ctx); err != nil {
		t.Fatalf("second Close() error = %v, want nil", err)
	}
	if err := client.Ping(ctx); !foundationx.IsKind(err, foundationx.ErrorKindConnection) {
		t.Fatalf("Ping() after Close() error = %v, want connection error", err)
	}
}

func TestClientQueryStatsAndHealthIntegration(t *testing.T) {
	const appName = "postgresx-client-query-health"
	ctx := t.Context()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: appName})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.WithoutCancel(ctx), t, client) })

	if client.Name() != "postgresx" {
		t.Fatalf("Name() = %q, want postgresx", client.Name())
	}
	if err := client.Ping(ctx); err != nil {
		t.Fatalf("Ping() error = %v, want nil", err)
	}

	if _, err := client.Exec(ctx, `CREATE TABLE postgresx_client_items (id BIGINT PRIMARY KEY, name TEXT NOT NULL)`); err != nil {
		t.Fatalf("create client items table: %v", err)
	}
	tag, err := client.Queryer().Exec(ctx, `INSERT INTO postgresx_client_items (id, name) VALUES ($1, $2), ($3, $4)`, int64(1), "alpha", int64(2), "beta")
	if err != nil {
		t.Fatalf("Queryer().Exec() error = %v, want nil", err)
	}
	if got := tag.RowsAffected(); got != 2 {
		t.Fatalf("RowsAffected() = %d, want 2", got)
	}

	var one string
	if err := client.Queryer().QueryRow(ctx, `SELECT name FROM postgresx_client_items WHERE id = $1`, int64(1)).Scan(&one); err != nil {
		t.Fatalf("Queryer().QueryRow() error = %v, want nil", err)
	}
	if one != "alpha" {
		t.Fatalf("selected name = %q, want alpha", one)
	}

	rows, err := client.Query(ctx, `SELECT name FROM postgresx_client_items ORDER BY id`)
	if err != nil {
		t.Fatalf("Query() error = %v, want nil", err)
	}
	defer rows.Close()
	names := make([]string, 0, 2)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Rows.Scan() error = %v, want nil", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows.Err() = %v, want nil", err)
	}
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Fatalf("names = %#v, want [alpha beta]", names)
	}

	err = client.QueryRow(ctx, `SELECT name FROM postgresx_client_items WHERE id = $1`, int64(99)).Scan(new(string))
	if !foundationx.IsKind(err, foundationx.ErrorKindNotFound) {
		t.Fatalf("missing QueryRow() error = %v, want not found", err)
	}

	stats := client.Stats()
	if stats.MaxConns == 0 {
		t.Fatalf("Stats().MaxConns = %d, want > 0", stats.MaxConns)
	}

	status := client.Check(ctx)
	if status.Name != "postgresx" {
		t.Fatalf("Check().Name = %q, want postgresx", status.Name)
	}
	if !status.IsHealthy() || status.Status != foundationx.HealthHealthy {
		t.Fatalf("Check() = %#v, want healthy", status)
	}
	if status.Metadata["application_name"] != appName {
		t.Fatalf("Check().Metadata[application_name] = %q, want %q", status.Metadata["application_name"], appName)
	}
	if status.Metadata["database"] == "" || status.Metadata["pool_max_conns"] == "" {
		t.Fatalf("Check().Metadata = %#v, want database and pool_max_conns", status.Metadata)
	}
}

func TestTransactionCommitRollbackAndReadOnlyIntegration(t *testing.T) {
	ctx := t.Context()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-transaction"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.WithoutCancel(ctx), t, client) })

	if _, err := client.Exec(ctx, `CREATE TABLE postgresx_tx_items (id BIGINT PRIMARY KEY, name TEXT NOT NULL)`); err != nil {
		t.Fatalf("create tx items table: %v", err)
	}

	if err := client.WithTx(ctx, func(ctx context.Context, tx postgresx.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO postgresx_tx_items (id, name) VALUES ($1, $2)`, int64(1), "committed"); err != nil {
			return err
		}
		var name string
		if err := tx.QueryRow(ctx, `SELECT name FROM postgresx_tx_items WHERE id = $1`, int64(1)).Scan(&name); err != nil {
			return err
		}
		if name != "committed" {
			return errors.New("transaction query returned unexpected value")
		}
		return nil
	}); err != nil {
		t.Fatalf("WithTx() commit error = %v, want nil", err)
	}
	assertIntegrationCount(ctx, t, client, `SELECT count(*) FROM postgresx_tx_items`, 1)

	rollbackErr := errors.New("force transaction rollback")
	err := client.WithTx(ctx, func(ctx context.Context, tx postgresx.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO postgresx_tx_items (id, name) VALUES ($1, $2)`, int64(2), "rolled-back"); err != nil {
			return err
		}
		return rollbackErr
	})
	if err == nil {
		t.Fatal("WithTx() rollback error = nil, want failure")
	}
	assertIntegrationCount(ctx, t, client, `SELECT count(*) FROM postgresx_tx_items`, 1)

	err = client.WithTxOptions(ctx, postgresx.TxOptions{ReadOnly: true}, func(ctx context.Context, tx postgresx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO postgresx_tx_items (id, name) VALUES ($1, $2)`, int64(3), "read-only")
		return err
	})
	if err == nil {
		t.Fatal("WithTxOptions(ReadOnly) write error = nil, want failure")
	}
	assertIntegrationCount(ctx, t, client, `SELECT count(*) FROM postgresx_tx_items`, 1)
}

func TestPostgresErrorMappingIntegration(t *testing.T) {
	ctx := t.Context()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-error-mapping"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.WithoutCancel(ctx), t, client) })

	if _, err := client.Exec(ctx, `CREATE TABLE postgresx_error_parents (id BIGINT PRIMARY KEY)`); err != nil {
		t.Fatalf("create error parents table: %v", err)
	}
	if _, err := client.Exec(ctx, `CREATE TABLE postgresx_error_children (id BIGINT PRIMARY KEY, parent_id BIGINT NOT NULL REFERENCES postgresx_error_parents(id), code TEXT NOT NULL CHECK (length(code) > 1))`); err != nil {
		t.Fatalf("create error children table: %v", err)
	}
	if _, err := client.Exec(ctx, `INSERT INTO postgresx_error_parents (id) VALUES ($1)`, int64(1)); err != nil {
		t.Fatalf("seed error parent: %v", err)
	}

	_, err := client.Exec(ctx, `INSERT INTO postgresx_error_parents (id) VALUES ($1)`, int64(1))
	if !foundationx.IsKind(err, foundationx.ErrorKindAlreadyExist) {
		t.Fatalf("duplicate key error = %v, want already_exists", err)
	}
	if postgresx.IsRetryable(err) {
		t.Fatalf("duplicate key retryable = true, want false")
	}

	_, err = client.Exec(ctx, `INSERT INTO postgresx_error_children (id, parent_id, code) VALUES ($1, $2, $3)`, int64(1), int64(999), "ok")
	if !foundationx.IsKind(err, foundationx.ErrorKindConflict) {
		t.Fatalf("foreign key error = %v, want conflict", err)
	}

	_, err = client.Exec(ctx, `INSERT INTO postgresx_error_children (id, parent_id, code) VALUES ($1, $2, $3)`, int64(2), int64(1), nil)
	if !foundationx.IsKind(err, foundationx.ErrorKindValidation) {
		t.Fatalf("not null error = %v, want validation", err)
	}

	_, err = client.Exec(ctx, `INSERT INTO postgresx_error_children (id, parent_id, code) VALUES ($1, $2, $3)`, int64(3), int64(1), "x")
	if !foundationx.IsKind(err, foundationx.ErrorKindValidation) {
		t.Fatalf("check constraint error = %v, want validation", err)
	}
}

func TestTestkitExplicitDSNIntegration(t *testing.T) {
	ctx := t.Context()
	dsn := currentIntegrationDSN(t)
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{
		DSN:             dsn,
		ApplicationName: "postgresx-testkit-explicit-dsn",
	})
	if err := fixture.Client().Ping(ctx); err != nil {
		t.Fatalf("explicit DSN fixture Ping() error = %v, want nil", err)
	}
}

func TestMigrationRunnerUpIntegration(t *testing.T) {
	ctx := t.Context()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-migration-up"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.WithoutCancel(ctx), t, client) })

	runner := postgresx.NewMigrationRunner(client)
	source := migrationSource{
		{Version: 2, Name: "seed_items", UpSQL: `INSERT INTO postgresx_integration_items (id, name) VALUES (1, 'alpha')`},
		{Version: 1, Name: "create_items", UpSQL: `CREATE TABLE postgresx_integration_items (id BIGINT PRIMARY KEY, name TEXT NOT NULL)`},
	}

	if err := runner.Up(ctx, source); err != nil {
		t.Fatalf("Up() error = %v, want nil", err)
	}
	assertAppliedMigrations(ctx, t, runner, []string{"create_items", "seed_items"})

	var count int
	if err := client.QueryRow(ctx, `SELECT count(*) FROM postgresx_integration_items`).Scan(&count); err != nil {
		t.Fatalf("count integration items: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}

	if err := runner.Up(ctx, source); err != nil {
		t.Fatalf("second Up() error = %v, want nil", err)
	}
	assertAppliedMigrations(ctx, t, runner, []string{"create_items", "seed_items"})
	if err := client.QueryRow(ctx, `SELECT count(*) FROM postgresx_integration_items`).Scan(&count); err != nil {
		t.Fatalf("count integration items after second up: %v", err)
	}
	if count != 1 {
		t.Fatalf("count after second up = %d, want 1", count)
	}
}

func TestMigrationRunnerDetectsVersionNameConflictIntegration(t *testing.T) {
	ctx := t.Context()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-migration-conflict"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.WithoutCancel(ctx), t, client) })

	runner := postgresx.NewMigrationRunner(client)
	if err := runner.Up(ctx, migrationSource{{Version: 1, Name: "original", UpSQL: `SELECT 1`}}); err != nil {
		t.Fatalf("Up() error = %v, want nil", err)
	}

	err := runner.Up(ctx, migrationSource{{Version: 1, Name: "renamed", UpSQL: `SELECT 1`}})
	if !foundationx.IsKind(err, foundationx.ErrorKindConflict) {
		t.Fatalf("Up() error = %v, want conflict", err)
	}
}

func TestMigrationRunnerRollsBackFailedMigrationIntegration(t *testing.T) {
	ctx := t.Context()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-migration-rollback"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.WithoutCancel(ctx), t, client) })

	runner := postgresx.NewMigrationRunner(client)
	err := runner.Up(ctx, migrationSource{{
		Version: 1,
		Name:    "failed_create",
		UpSQL:   `CREATE TABLE postgresx_rollback_items (id BIGINT PRIMARY KEY); INSERT INTO postgresx_missing_table (id) VALUES (1)`,
	}})
	if err == nil {
		t.Fatal("Up() error = nil, want failure")
	}

	var exists bool
	if err := client.QueryRow(ctx, `SELECT to_regclass('public.postgresx_rollback_items') IS NOT NULL`).Scan(&exists); err != nil {
		t.Fatalf("check rollback table: %v", err)
	}
	if exists {
		t.Fatal("rollback table exists, want failed migration to roll back")
	}

	applied, err := runner.Applied(ctx)
	if err != nil {
		t.Fatalf("Applied() error = %v, want nil", err)
	}
	if len(applied) != 0 {
		t.Fatalf("applied migrations = %d, want 0", len(applied))
	}
}

func currentIntegrationDSN(t *testing.T) string {
	t.Helper()
	if dsn := os.Getenv("POSTGRESX_INTEGRATION_DSN"); dsn != "" {
		return dsn
	}
	if dsn := os.Getenv("POSTGRES_TEST_DSN"); dsn != "" {
		return dsn
	}
	t.Skip("POSTGRESX_INTEGRATION_DSN or POSTGRES_TEST_DSN is not set")
	return ""
}

func assertIntegrationCount(ctx context.Context, t *testing.T, client *postgresx.Client, query string, want int) {
	t.Helper()
	var count int
	if err := client.QueryRow(ctx, query).Scan(&count); err != nil {
		t.Fatalf("count query: %v", err)
	}
	if count != want {
		t.Fatalf("count = %d, want %d", count, want)
	}
}

func assertAppliedMigrations(ctx context.Context, t *testing.T, runner *postgresx.MigrationRunner, names []string) {
	t.Helper()
	applied, err := runner.Applied(ctx)
	if err != nil {
		t.Fatalf("Applied() error = %v, want nil", err)
	}
	if len(applied) != len(names) {
		t.Fatalf("applied migrations = %d, want %d", len(applied), len(names))
	}
	for i, name := range names {
		if applied[i].Name != name {
			t.Fatalf("applied[%d].Name = %q, want %q", i, applied[i].Name, name)
		}
	}
}

func cleanupIntegrationSchema(ctx context.Context, t *testing.T, client *postgresx.Client) {
	t.Helper()
	for _, stmt := range []string{
		`DROP TABLE IF EXISTS postgresx_error_children`,
		`DROP TABLE IF EXISTS postgresx_error_parents`,
		`DROP TABLE IF EXISTS postgresx_tx_items`,
		`DROP TABLE IF EXISTS postgresx_client_items`,
		`DROP TABLE IF EXISTS postgresx_integration_items`,
		`DROP TABLE IF EXISTS postgresx_rollback_items`,
		`DROP TABLE IF EXISTS schema_migrations`,
	} {
		if _, err := client.Exec(ctx, stmt); err != nil {
			t.Fatalf("cleanup integration schema: %v", err)
		}
	}
}
