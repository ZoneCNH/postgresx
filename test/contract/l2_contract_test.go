package contract_test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
	"github.com/ZoneCNH/postgresx/test/postgresxtest"
)

func TestP0SQLContract(t *testing.T) {
	ctx := context.Background()
	queryer := &postgresxtest.QueryAdapter{
		ExecTag:   postgresxtest.CommandTag{Rows: 2},
		QueryRows: &postgresxtest.Rows{Rows: [][]any{{1, "one"}, {2, "two"}}},
		Row:       &postgresxtest.Row{Values: []any{3, "three"}},
	}

	tag, err := queryer.Exec(ctx, "insert into l2_contract(id) values($1)", 1)
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}
	if got := tag.RowsAffected(); got != 2 {
		t.Fatalf("RowsAffected() = %d, want 2", got)
	}

	rows, err := queryer.Query(ctx, "select id, name from l2_contract where id > $1", 0)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	defer rows.Close()
	var seen []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			t.Fatalf("Rows.Scan() error = %v", err)
		}
		seen = append(seen, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows.Err() = %v", err)
	}
	if !reflect.DeepEqual(seen, []string{"one", "two"}) {
		t.Fatalf("rows = %v, want [one two]", seen)
	}

	var id int
	var name string
	if err := queryer.QueryRow(ctx, "select id, name from l2_contract where id=$1", 3).Scan(&id, &name); err != nil {
		t.Fatalf("QueryRow().Scan() error = %v", err)
	}
	if id != 3 || name != "three" {
		t.Fatalf("row = (%d, %q), want (3, three)", id, name)
	}

	wantCalls := []postgresxtest.QueryCall{
		{Operation: "exec", SQL: "insert into l2_contract(id) values($1)", Args: []any{1}},
		{Operation: "query", SQL: "select id, name from l2_contract where id > $1", Args: []any{0}},
		{Operation: "query_row", SQL: "select id, name from l2_contract where id=$1", Args: []any{3}},
	}
	if !reflect.DeepEqual(queryer.Calls, wantCalls) {
		t.Fatalf("calls = %#v, want %#v", queryer.Calls, wantCalls)
	}
}

func TestP0SQLContractErrorPropagation(t *testing.T) {
	ctx := context.Background()
	execErr := errors.New("exec failed")
	queryErr := errors.New("query failed")
	rowErr := errors.New("row failed")
	rowsErr := errors.New("rows failed")

	queryer := &postgresxtest.QueryAdapter{ExecErr: execErr}
	if _, err := queryer.Exec(ctx, "delete from l2_contract where id=$1", 10); !errors.Is(err, execErr) {
		t.Fatalf("Exec() error = %v, want %v", err, execErr)
	}
	if len(queryer.Calls) != 1 || queryer.Calls[0].Operation != "exec" {
		t.Fatalf("Exec() calls = %#v, want one exec call", queryer.Calls)
	}

	queryer = &postgresxtest.QueryAdapter{QueryErr: queryErr}
	if _, err := queryer.Query(ctx, "select id from l2_contract"); !errors.Is(err, queryErr) {
		t.Fatalf("Query() error = %v, want %v", err, queryErr)
	}

	queryer = &postgresxtest.QueryAdapter{Row: &postgresxtest.Row{Err: rowErr}}
	if err := queryer.QueryRow(ctx, "select id from l2_contract where id=$1", 11).Scan(new(int)); !errors.Is(err, rowErr) {
		t.Fatalf("QueryRow().Scan() error = %v, want %v", err, rowErr)
	}

	queryer = &postgresxtest.QueryAdapter{QueryRows: &postgresxtest.Rows{Rows: [][]any{{12}}, ErrValue: rowsErr}}
	rows, err := queryer.Query(ctx, "select id from l2_contract")
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if !rows.Next() {
		t.Fatal("Rows.Next() = false, want one row before Err()")
	}
	var id int
	if err := rows.Scan(&id); err != nil {
		t.Fatalf("Rows.Scan() error = %v", err)
	}
	if err := rows.Err(); !errors.Is(err, rowsErr) {
		t.Fatalf("Rows.Err() = %v, want %v", err, rowsErr)
	}
}

func TestP0TxContract(t *testing.T) {
	var _ postgresx.Tx = (*postgresxtest.QueryAdapter)(nil)

	ctx := context.Background()
	tx := &postgresxtest.QueryAdapter{ExecTag: postgresxtest.CommandTag{Rows: 1}}
	txFn := postgresx.TxFunc(func(ctx context.Context, tx postgresx.Tx) error {
		tag, err := tx.Exec(ctx, "update l2_contract set seen=true where id=$1", 7)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			t.Fatalf("RowsAffected() = %d, want 1", tag.RowsAffected())
		}
		return nil
	})
	if err := txFn(ctx, tx); err != nil {
		t.Fatalf("TxFunc() error = %v", err)
	}
	if len(tx.Calls) != 1 || tx.Calls[0].Operation != "exec" {
		t.Fatalf("tx calls = %#v, want one exec", tx.Calls)
	}
}

func TestP0TxContractPropagatesExecutorErrors(t *testing.T) {
	ctx := context.Background()
	execErr := errors.New("tx exec failed")
	tx := &postgresxtest.QueryAdapter{ExecErr: execErr}
	txFn := postgresx.TxFunc(func(ctx context.Context, tx postgresx.Tx) error {
		_, err := tx.Exec(ctx, "update l2_contract set seen=true where id=$1", 17)
		return err
	})

	if err := txFn(ctx, tx); !errors.Is(err, execErr) {
		t.Fatalf("TxFunc() error = %v, want %v", err, execErr)
	}
	if len(tx.Calls) != 1 || tx.Calls[0].Operation != "exec" {
		t.Fatalf("tx calls = %#v, want one exec", tx.Calls)
	}
}

func TestP0LivePostgresTxRollbackContract(t *testing.T) {
	ctx := t.Context()
	fixture := postgresxtest.Start(ctx, t, "postgresx-l2-rollback-contract")
	client := fixture.Client()
	sentinel := errors.New("rollback contract sentinel")

	if _, err := client.Exec(ctx, "create temporary table l2_contract_rollback(id integer) on commit preserve rows"); err != nil {
		t.Fatalf("create temporary table: %v", err)
	}
	err := client.WithTx(ctx, func(ctx context.Context, tx postgresx.Tx) error {
		if _, err := tx.Exec(ctx, "insert into l2_contract_rollback(id) values($1)", 1); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("WithTx() error = %v, want sentinel rollback error", err)
	}

	var count int
	if err := client.QueryRow(ctx, "select count(*) from l2_contract_rollback").Scan(&count); err != nil {
		t.Fatalf("count rollback table: %v", err)
	}
	if count != 0 {
		t.Fatalf("rollback table count = %d, want 0", count)
	}
}

func TestP0PoolContract(t *testing.T) {
	cfg := postgresx.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Database = "postgres"
	cfg.User = "postgres"
	cfg.Password = foundationx.NewSecretString("contract-secret")
	cfg.MaxOpenConns = 3
	cfg.MinIdleConns = 1
	cfg.ApplicationName = "postgresx-l2-contract"

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	dsn := cfg.RedactedDSN()
	if strings.Contains(dsn, "contract-secret") {
		t.Fatalf("RedactedDSN() = %q, leaked password", dsn)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse RedactedDSN(): %v", err)
	}
	if masked, ok := parsed.User.Password(); !ok || masked != "***" {
		t.Fatalf("RedactedDSN() password = %q, %v; want masked password", masked, ok)
	}
	sanitized := cfg.Sanitize()
	if sanitized.Password != "***" {
		t.Fatalf("Sanitize().Password = %q, want mask", sanitized.Password)
	}

	statsType := reflect.TypeOf(postgresx.PoolStats{})
	for _, field := range []string{"TotalConns", "IdleConns", "AcquiredConns", "ConstructingConns", "MaxConns"} {
		if _, ok := statsType.FieldByName(field); !ok {
			t.Fatalf("PoolStats missing field %s", field)
		}
	}
}

func TestP0LivePostgresContract(t *testing.T) {
	fixture := postgresxtest.Start(t.Context(), t, "postgresx-l2-contract")
	client := fixture.Client()

	if err := client.Ping(t.Context()); err != nil {
		t.Fatalf("Ping() error = %v", err)
	}
	if stats := client.Stats(); stats.MaxConns <= 0 {
		t.Fatalf("Stats().MaxConns = %d, want positive", stats.MaxConns)
	}
	if err := client.WithTx(t.Context(), func(ctx context.Context, tx postgresx.Tx) error {
		var got int
		if err := tx.QueryRow(ctx, "select 1").Scan(&got); err != nil {
			return err
		}
		if got != 1 {
			t.Fatalf("select 1 = %d, want 1", got)
		}
		return nil
	}); err != nil {
		t.Fatalf("WithTx() error = %v", err)
	}
}

func TestFormalPackagesDoNotDependOnL2TestTooling(t *testing.T) {
	repoRoot := findRepoRoot(t)
	for _, dir := range []string{"pkg", "internal"} {
		root := filepath.Join(repoRoot, dir)
		if _, err := os.Stat(root); err != nil {
			continue
		}
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			text := string(contents)
			for _, forbidden := range []string{"testkitx", "xlibgate", "xlib-standard"} {
				if strings.Contains(text, forbidden) {
					t.Fatalf("formal package file %s contains forbidden L2 tooling dependency marker %q", path, forbidden)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
