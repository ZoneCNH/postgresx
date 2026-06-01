package postgresx

import (
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
)

func TestValidateMigrations(t *testing.T) {
	migrations := []Migration{
		{Version: 1, Name: "create_accounts", UpSQL: "create table accounts(id bigint primary key)"},
		{Version: 2, Name: "add_index", UpSQL: "create index accounts_id_idx on accounts(id)"},
	}

	if err := validateMigrations(migrations); err != nil {
		t.Fatalf("validateMigrations() error = %v", err)
	}
}

func TestValidateMigrationsRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		migrations []Migration
		kind       foundationx.ErrorKind
	}{
		{
			name:       "non-positive version",
			migrations: []Migration{{Version: 0, Name: "bad", UpSQL: "select 1"}},
			kind:       foundationx.ErrorKindValidation,
		},
		{
			name:       "missing name",
			migrations: []Migration{{Version: 1, UpSQL: "select 1"}},
			kind:       foundationx.ErrorKindValidation,
		},
		{
			name:       "missing up sql",
			migrations: []Migration{{Version: 1, Name: "empty"}},
			kind:       foundationx.ErrorKindValidation,
		},
		{
			name: "duplicate version",
			migrations: []Migration{
				{Version: 1, Name: "one", UpSQL: "select 1"},
				{Version: 1, Name: "two", UpSQL: "select 2"},
			},
			kind: foundationx.ErrorKindConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMigrations(tt.migrations)
			if !foundationx.IsKind(err, tt.kind) {
				t.Fatalf("validateMigrations() error = %v, want kind %s", err, tt.kind)
			}
		})
	}
}
