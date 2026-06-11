package postgresx

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
)

// Migration describes one caller-owned migration.
type Migration struct {
	Version int64
	Name    string
	UpSQL   string
	DownSQL string
}

// MigrationSource lists caller-owned migrations.
type MigrationSource interface {
	List(ctx context.Context) ([]Migration, error)
}

// AppliedMigration describes a migration recorded in schema_migrations.
type AppliedMigration struct {
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	AppliedAt time.Time `json:"applied_at"`
}

// MigrationRunner applies caller-provided migrations in version order.
type MigrationRunner struct {
	client *Client
}

// NewMigrationRunner creates a migration runner for client.
func NewMigrationRunner(client *Client) *MigrationRunner {
	return &MigrationRunner{client: client}
}

// Up applies pending migrations.
func (r *MigrationRunner) Up(ctx context.Context, source MigrationSource) error {
	const op = "postgresx.MigrationRunner.Up"
	if r == nil || r.client == nil {
		return foundationx.NewError(foundationx.ErrorKindConfig, op, "client is required")
	}
	if source == nil {
		return foundationx.NewError(foundationx.ErrorKindConfig, op, "migration source is required")
	}
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	migrations, err := source.List(ctx)
	if err != nil {
		return MapError(op, err)
	}
	slices.SortFunc(migrations, func(a, b Migration) int {
		if a.Version < b.Version {
			return -1
		}
		if a.Version > b.Version {
			return 1
		}
		return 0
	})
	if err := validateMigrations(migrations); err != nil {
		return err
	}

	applied, err := r.Applied(ctx)
	if err != nil {
		return err
	}
	appliedByVersion := make(map[int64]AppliedMigration, len(applied))
	for _, item := range applied {
		appliedByVersion[item.Version] = item
	}

	for _, migration := range migrations {
		if existing, ok := appliedByVersion[migration.Version]; ok {
			if existing.Name != migration.Name {
				return foundationx.WrapError(foundationx.ErrorKindConflict, op, "migration version already applied with different name", fmt.Errorf("version %d: applied %q, requested %q", migration.Version, existing.Name, migration.Name))
			}
			continue
		}
		if err := r.client.WithTx(ctx, func(ctx context.Context, tx Tx) error {
			if _, err := tx.Exec(ctx, migration.UpSQL); err != nil {
				return err
			}
			_, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`, migration.Version, migration.Name)
			return err
		}); err != nil {
			return MapError(op, err)
		}
	}
	return nil
}

// Applied returns applied migrations ordered by version.
func (r *MigrationRunner) Applied(ctx context.Context) ([]AppliedMigration, error) {
	if err := r.ensureTable(ctx); err != nil {
		return nil, err
	}
	rows, err := r.client.Query(ctx, `SELECT version, name, applied_at FROM schema_migrations ORDER BY version`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applied []AppliedMigration
	for rows.Next() {
		var item AppliedMigration
		if err := rows.Scan(&item.Version, &item.Name, &item.AppliedAt); err != nil {
			return nil, err
		}
		applied = append(applied, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return applied, nil
}

func (r *MigrationRunner) ensureTable(ctx context.Context) error {
	if r == nil || r.client == nil {
		return foundationx.NewError(foundationx.ErrorKindConfig, "postgresx.MigrationRunner.ensureTable", "client is required")
	}
	_, err := r.client.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	return err
}

func validateMigrations(migrations []Migration) error {
	seen := make(map[int64]string, len(migrations))
	for _, migration := range migrations {
		if migration.Version <= 0 {
			return foundationx.NewError(foundationx.ErrorKindValidation, "postgresx.validateMigrations", "migration version must be positive")
		}
		if migration.Name == "" {
			return foundationx.NewError(foundationx.ErrorKindValidation, "postgresx.validateMigrations", "migration name is required")
		}
		if migration.UpSQL == "" {
			return foundationx.NewError(foundationx.ErrorKindValidation, "postgresx.validateMigrations", "migration up sql is required")
		}
		if previous, ok := seen[migration.Version]; ok {
			return foundationx.WrapError(foundationx.ErrorKindConflict, "postgresx.validateMigrations", "duplicate migration version", fmt.Errorf("version %d: %q and %q", migration.Version, previous, migration.Name))
		}
		seen[migration.Version] = migration.Name
	}
	return nil
}
