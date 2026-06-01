package postgresx_test

import (
	"context"
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/bytechainx/postgresx"
	"github.com/bytechainx/postgresx/testkit"
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

func TestMigrationRunnerUpIntegration(t *testing.T) {
	ctx := context.Background()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-migration-up"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.Background(), t, client) })

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
	ctx := context.Background()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-migration-conflict"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.Background(), t, client) })

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
	ctx := context.Background()
	fixture := testkit.StartPostgres(ctx, t, testkit.Options{ApplicationName: "postgresx-migration-rollback"})
	client := fixture.Client()
	cleanupIntegrationSchema(ctx, t, client)
	t.Cleanup(func() { cleanupIntegrationSchema(context.Background(), t, client) })

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
		`DROP TABLE IF EXISTS postgresx_integration_items`,
		`DROP TABLE IF EXISTS postgresx_rollback_items`,
		`DROP TABLE IF EXISTS schema_migrations`,
	} {
		if _, err := client.Exec(ctx, stmt); err != nil {
			t.Fatalf("cleanup integration schema: %v", err)
		}
	}
}
