package testkit

import (
	"context"
	"errors"
	"testing"
)

func TestFakeQueryerCapturesCallsAndReturnsConfiguredResults(t *testing.T) {
	ctx := context.Background()
	queryer := &FakeQueryer{
		ExecTag:   FakeCommandTag{Rows: 2},
		QueryRows: &FakeRows{Rows: [][]any{{1, "one"}, {2, "two"}}},
		Row:       &FakeRow{Values: []any{3, "three"}},
	}

	tag, err := queryer.Exec(ctx, "update things set name=$1", "one")
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}
	if got := tag.RowsAffected(); got != 2 {
		t.Fatalf("RowsAffected() = %d, want 2", got)
	}

	rows, err := queryer.Query(ctx, "select id, name from things")
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	var seen []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			t.Fatalf("Rows.Scan() error = %v", err)
		}
		seen = append(seen, name)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows.Err() error = %v", err)
	}

	var id int
	var name string
	if err := queryer.QueryRow(ctx, "select id, name from things where id=$1", 3).Scan(&id, &name); err != nil {
		t.Fatalf("QueryRow().Scan() error = %v", err)
	}

	if len(seen) != 2 || seen[0] != "one" || seen[1] != "two" {
		t.Fatalf("seen rows = %+v, want [one two]", seen)
	}
	if id != 3 || name != "three" {
		t.Fatalf("row = (%d, %q), want (3, three)", id, name)
	}
	if len(queryer.Calls) != 3 {
		t.Fatalf("calls = %+v, want 3 calls", queryer.Calls)
	}
	if queryer.Calls[0].Operation != "exec" || queryer.Calls[0].SQL != "update things set name=$1" {
		t.Fatalf("first call = %+v, want exec call", queryer.Calls[0])
	}
}

func TestFakeQueryerReturnsConfiguredErrors(t *testing.T) {
	ctx := context.Background()
	execErr := errors.New("exec failed")
	queryErr := errors.New("query failed")
	rowErr := errors.New("row failed")
	queryer := &FakeQueryer{
		ExecErr:  execErr,
		QueryErr: queryErr,
		Row:      &FakeRow{Err: rowErr},
	}

	if _, err := queryer.Exec(ctx, "delete from things"); !errors.Is(err, execErr) {
		t.Fatalf("Exec() error = %v, want %v", err, execErr)
	}
	if _, err := queryer.Query(ctx, "select * from things"); !errors.Is(err, queryErr) {
		t.Fatalf("Query() error = %v, want %v", err, queryErr)
	}
	if err := queryer.QueryRow(ctx, "select 1").Scan(new(int)); !errors.Is(err, rowErr) {
		t.Fatalf("QueryRow().Scan() error = %v, want %v", err, rowErr)
	}
}
