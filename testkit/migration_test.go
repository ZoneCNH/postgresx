package testkit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMigrationFile(t *testing.T) {
	tests := []struct {
		name          string
		version       int64
		migrationName string
		direction     string
	}{
		{name: "001_create_accounts.sql", version: 1, migrationName: "create_accounts", direction: "up"},
		{name: "002-add-index.up.sql", version: 2, migrationName: "add-index", direction: "up"},
		{name: "002-add-index.down.sql", version: 2, migrationName: "add-index", direction: "down"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, migrationName, direction, err := parseMigrationFile(tt.name)
			if err != nil {
				t.Fatalf("parseMigrationFile() error = %v", err)
			}
			if version != tt.version || migrationName != tt.migrationName || direction != tt.direction {
				t.Fatalf("parseMigrationFile() = (%d, %q, %q), want (%d, %q, %q)", version, migrationName, direction, tt.version, tt.migrationName, tt.direction)
			}
		})
	}
}

func TestFileMigrationSourceList(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"002_add_index.up.sql":    "create index accounts_id_idx on accounts(id)",
		"002_add_index.down.sql":  "drop index accounts_id_idx",
		"001_create_accounts.sql": "create table accounts(id bigint primary key)",
		"ignored_not_sql.txt":     "ignored",
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatalf("write migration fixture %s: %v", name, err)
		}
	}

	migrations, err := (FileMigrationSource{Dir: dir}).List(t.Context())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(migrations) != 2 {
		t.Fatalf("len(List()) = %d, want 2", len(migrations))
	}
	if migrations[0].Version != 1 || migrations[0].Name != "create_accounts" || migrations[0].UpSQL == "" {
		t.Fatalf("first migration = %+v, want version 1 create_accounts", migrations[0])
	}
	if migrations[1].Version != 2 || migrations[1].Name != "add_index" || migrations[1].UpSQL == "" || migrations[1].DownSQL == "" {
		t.Fatalf("second migration = %+v, want version 2 add_index with up/down SQL", migrations[1])
	}
}
