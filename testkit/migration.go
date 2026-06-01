package testkit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/bytechainx/postgresx"
)

// FileMigrationSource loads migrations from files named
// <version>_<name>.sql, <version>_<name>.up.sql, or <version>_<name>.down.sql.
type FileMigrationSource struct {
	Dir string
}

// List implements postgresx.MigrationSource.
func (s FileMigrationSource) List(ctx context.Context) ([]postgresx.Migration, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}

	byVersion := map[int64]*postgresx.Migration{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		version, name, direction, err := parseMigrationFile(entry.Name())
		if err != nil {
			return nil, err
		}
		contents, err := os.ReadFile(filepath.Join(s.Dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		migration := byVersion[version]
		if migration == nil {
			migration = &postgresx.Migration{Version: version, Name: name}
			byVersion[version] = migration
		}
		switch direction {
		case "down":
			migration.DownSQL = string(contents)
		default:
			migration.UpSQL = string(contents)
		}
	}

	migrations := make([]postgresx.Migration, 0, len(byVersion))
	for _, migration := range byVersion {
		migrations = append(migrations, *migration)
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	return migrations, nil
}

// MigrateUp runs file migrations against a configured integration database.
func MigrateUp(ctx context.Context, t testing.TB, migrationsPath string) {
	t.Helper()
	fixture := StartPostgres(ctx, t, Options{ApplicationName: "postgresx-testkit-migration"})
	runner := postgresx.NewMigrationRunner(fixture.Client())
	if err := runner.Up(ctx, FileMigrationSource{Dir: migrationsPath}); err != nil {
		t.Fatalf("migrate up: %v", err)
	}
}

func parseMigrationFile(name string) (int64, string, string, error) {
	base := strings.TrimSuffix(name, ".sql")
	direction := "up"
	if strings.HasSuffix(base, ".up") {
		base = strings.TrimSuffix(base, ".up")
	} else if strings.HasSuffix(base, ".down") {
		base = strings.TrimSuffix(base, ".down")
		direction = "down"
	}

	separator := strings.IndexAny(base, "_-")
	if separator <= 0 || separator == len(base)-1 {
		return 0, "", "", fmt.Errorf("invalid migration file name %q", name)
	}
	version, err := strconv.ParseInt(base[:separator], 10, 64)
	if err != nil || version <= 0 {
		return 0, "", "", fmt.Errorf("invalid migration version in %q", name)
	}
	return version, base[separator+1:], direction, nil
}
