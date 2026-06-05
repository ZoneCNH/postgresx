package postgresxtest

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestQueryAdapterPropagatesConfiguredErrors(t *testing.T) {
	ctx := context.Background()
	execErr := errors.New("exec failed")
	queryErr := errors.New("query failed")
	rowErr := errors.New("row failed")
	rowsErr := errors.New("rows failed")

	adapter := &QueryAdapter{ExecErr: execErr}
	if _, err := adapter.Exec(ctx, "insert into contract values($1)", 1); !errors.Is(err, execErr) {
		t.Fatalf("Exec() error = %v, want %v", err, execErr)
	}
	if len(adapter.Calls) != 1 || adapter.Calls[0].Operation != "exec" {
		t.Fatalf("Exec() calls = %#v, want one exec call", adapter.Calls)
	}

	adapter = &QueryAdapter{QueryErr: queryErr}
	if _, err := adapter.Query(ctx, "select 1"); !errors.Is(err, queryErr) {
		t.Fatalf("Query() error = %v, want %v", err, queryErr)
	}

	adapter = &QueryAdapter{Row: &Row{Err: rowErr}}
	if err := adapter.QueryRow(ctx, "select 1").Scan(new(int)); !errors.Is(err, rowErr) {
		t.Fatalf("QueryRow().Scan() error = %v, want %v", err, rowErr)
	}

	adapter = &QueryAdapter{QueryRows: &Rows{Rows: [][]any{{1}}, ErrValue: rowsErr}}
	rows, err := adapter.Query(ctx, "select 1")
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if !rows.Next() {
		t.Fatal("Rows.Next() = false, want one row")
	}
	var got int
	if err := rows.Scan(&got); err != nil {
		t.Fatalf("Rows.Scan() error = %v", err)
	}
	if err := rows.Err(); !errors.Is(err, rowsErr) {
		t.Fatalf("Rows.Err() = %v, want %v", err, rowsErr)
	}
}

func TestRowsScanRequiresCurrentRowAndDestShape(t *testing.T) {
	rows := &Rows{Rows: [][]any{{1, "one"}}}
	var id int
	if err := rows.Scan(&id); err == nil || !strings.Contains(err.Error(), "without current row") {
		t.Fatalf("Rows.Scan() before Next error = %v, want current-row guard", err)
	}
	if !rows.Next() {
		t.Fatal("Rows.Next() = false, want one row")
	}
	if err := rows.Scan(&id); err == nil || !strings.Contains(err.Error(), "value count") {
		t.Fatalf("Rows.Scan() dest count error = %v, want value count guard", err)
	}
	var name string
	if err := rows.Scan(&id, &name); err != nil {
		t.Fatalf("Rows.Scan() valid error = %v", err)
	}
	if id != 1 || name != "one" {
		t.Fatalf("Rows.Scan() got (%d, %q), want (1, one)", id, name)
	}
	if err := rows.Scan(nil, &name); err == nil || !strings.Contains(err.Error(), "non-nil pointer") {
		t.Fatalf("Rows.Scan() nil dest error = %v, want pointer guard", err)
	}
}

func TestRowsCloseRecordsClosed(t *testing.T) {
	rows := &Rows{}
	rows.Close()
	if !rows.Closed {
		t.Fatal("Rows.Close() did not record Closed=true")
	}
}

func TestScanValuesAssignsConvertibleAndNilValues(t *testing.T) {
	var small int32
	var text *string
	if err := scanValues([]any{int64(7), nil}, []any{&small, &text}); err != nil {
		t.Fatalf("scanValues() error = %v", err)
	}
	if small != 7 {
		t.Fatalf("converted int = %d, want 7", small)
	}
	if text != nil {
		t.Fatalf("nil destination = %v, want nil", text)
	}
}
