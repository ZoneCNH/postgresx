package postgresx

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCommandTagRowsAffected(t *testing.T) {
	tag := commandTag{tag: pgconn.NewCommandTag("INSERT 0 5")}
	if got := tag.RowsAffected(); got != 5 {
		t.Fatalf("RowsAffected() = %d, want 5", got)
	}
}

func TestCommandTagRowsAffectedZero(t *testing.T) {
	tag := commandTag{tag: pgconn.NewCommandTag("UPDATE 0")}
	if got := tag.RowsAffected(); got != 0 {
		t.Fatalf("RowsAffected() = %d, want 0", got)
	}
}

func TestErrorRowScan(t *testing.T) {
	want := errors.New("test error")
	row := errorRow{err: want}
	if got := row.Scan(new(int)); got != want {
		t.Fatalf("Scan() = %v, want %v", got, want)
	}
}

// mockRow implements pgx.Row.
type mockRow struct {
	err error
}

func (r mockRow) Scan(dest ...any) error {
	if len(dest) > 0 {
		if v, ok := dest[0].(*int); ok {
			*v = 42
		}
	}
	return r.err
}

func TestRowAdapterScanSuccess(t *testing.T) {
	row := rowAdapter{row: mockRow{}, op: "test.Scan"}
	var val int
	if err := row.Scan(&val); err != nil {
		t.Fatalf("Scan() error = %v, want nil", err)
	}
	if val != 42 {
		t.Fatalf("Scan() set value to %d, want 42", val)
	}
}

func TestRowAdapterScanError(t *testing.T) {
	want := errors.New("scan failure")
	row := rowAdapter{row: mockRow{err: want}, op: "test.Scan"}
	if err := row.Scan(new(int)); err == nil {
		t.Fatal("Scan() error = nil, want non-nil")
	}
}

// mockRows implements pgx.Rows.
type mockRows struct {
	data        []int
	pos         int
	closeCalled bool
	scanErr     error
	errVal      error
}

func (r *mockRows) Close()                        { r.closeCalled = true }
func (r *mockRows) Err() error                    { return r.errVal }
func (r *mockRows) CommandTag() pgconn.CommandTag { return pgconn.NewCommandTag("SELECT 0") }
func (r *mockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}
func (r *mockRows) Next() bool {
	if r.pos < len(r.data) {
		r.pos++
		return true
	}
	return false
}
func (r *mockRows) Scan(dest ...any) error {
	if r.scanErr != nil {
		return r.scanErr
	}
	if r.pos > 0 && r.pos <= len(r.data) {
		if v, ok := dest[0].(*int); ok {
			*v = r.data[r.pos-1]
		}
	}
	return nil
}
func (r *mockRows) Values() ([]any, error) { return nil, nil }
func (r *mockRows) RawValues() [][]byte    { return nil }
func (r *mockRows) Conn() *pgx.Conn        { return nil }

func TestRowsAdapterClose(t *testing.T) {
	rows := &mockRows{}
	adapter := rowsAdapter{rows: rows}
	adapter.Close()
	if !rows.closeCalled {
		t.Fatal("Close() was not called on underlying rows")
	}
}

func TestRowsAdapterErr(t *testing.T) {
	want := errors.New("rows error")
	rows := &mockRows{errVal: want}
	adapter := rowsAdapter{rows: rows, op: "test.Err"}
	err := adapter.Err()
	if err == nil {
		t.Fatal("Err() = nil, want error")
	}
}

func TestRowsAdapterErrNil(t *testing.T) {
	rows := &mockRows{}
	adapter := rowsAdapter{rows: rows, op: "test.Err"}
	if err := adapter.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
}

func TestRowsAdapterNext(t *testing.T) {
	rows := &mockRows{data: []int{1, 2}}
	adapter := rowsAdapter{rows: rows}
	if !adapter.Next() {
		t.Fatal("first Next() = false, want true")
	}
	if !adapter.Next() {
		t.Fatal("second Next() = false, want true")
	}
	if adapter.Next() {
		t.Fatal("third Next() = true, want false")
	}
}

func TestRowsAdapterScan(t *testing.T) {
	rows := &mockRows{data: []int{99}}
	adapter := rowsAdapter{rows: rows, op: "test.Scan"}
	if !adapter.Next() {
		t.Fatal("Next() = false, want true")
	}
	var val int
	if err := adapter.Scan(&val); err != nil {
		t.Fatalf("Scan() error = %v, want nil", err)
	}
	if val != 99 {
		t.Fatalf("Scan() = %d, want 99", val)
	}
}
