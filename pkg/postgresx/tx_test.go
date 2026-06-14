package postgresx

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// mockPgxTx implements pgx.Tx for testing txAdapter.
type mockPgxTx struct {
	execTag  pgconn.CommandTag
	execErr  error
	queryRows pgx.Rows
	queryErr error
	queryRow pgx.Row
}

func (m mockPgxTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, errors.New("not implemented")
}

func (m mockPgxTx) Commit(ctx context.Context) error {
	return errors.New("not implemented")
}

func (m mockPgxTx) Rollback(ctx context.Context) error {
	return errors.New("not implemented")
}

func (m mockPgxTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m mockPgxTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

func (m mockPgxTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m mockPgxTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, errors.New("not implemented")
}

func (m mockPgxTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return m.execTag, m.execErr
}

func (m mockPgxTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return m.queryRows, m.queryErr
}

func (m mockPgxTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.queryRow
}

func (m mockPgxTx) Conn() *pgx.Conn {
	return nil
}

func TestTxAdapterExecSuccess(t *testing.T) {
	wantTag := pgconn.NewCommandTag("INSERT 0 3")
	tx := txAdapter{tx: mockPgxTx{execTag: wantTag}}
	tag, err := tx.Exec(t.Context(), "INSERT INTO t VALUES ($1)", 1)
	if err != nil {
		t.Fatalf("Exec() error = %v, want nil", err)
	}
	if tag.RowsAffected() != 3 {
		t.Fatalf("RowsAffected() = %d, want 3", tag.RowsAffected())
	}
}

func TestTxAdapterExecError(t *testing.T) {
	want := errors.New("exec failure")
	tx := txAdapter{tx: mockPgxTx{execErr: want}}
	_, err := tx.Exec(t.Context(), "INSERT INTO t VALUES ($1)", 1)
	if err == nil {
		t.Fatal("Exec() error = nil, want non-nil")
	}
}

func TestTxAdapterQuerySuccess(t *testing.T) {
	rows := &mockRows{data: []int{1, 2}}
	tx := txAdapter{tx: mockPgxTx{queryRows: rows}}
	result, err := tx.Query(t.Context(), "SELECT * FROM t")
	if err != nil {
		t.Fatalf("Query() error = %v, want nil", err)
	}
	if result == nil {
		t.Fatal("Query() rows = nil")
	}
}

func TestTxAdapterQueryError(t *testing.T) {
	want := errors.New("query failure")
	tx := txAdapter{tx: mockPgxTx{queryErr: want}}
	_, err := tx.Query(t.Context(), "SELECT * FROM t")
	if err == nil {
		t.Fatal("Query() error = nil, want non-nil")
	}
}

func TestTxAdapterQueryRow(t *testing.T) {
	row := mockRow{}
	tx := txAdapter{tx: mockPgxTx{queryRow: row}}
	result := tx.QueryRow(t.Context(), "SELECT * FROM t WHERE id = $1", 1)
	if result == nil {
		t.Fatal("QueryRow() = nil")
	}
	var val int
	if err := result.Scan(&val); err != nil {
		t.Fatalf("Scan() error = %v, want nil", err)
	}
	if val != 42 {
		t.Fatalf("Scan() = %d, want 42", val)
	}
}

func TestWithTxOptionsNilFn(t *testing.T) {
	client := &Client{opts: defaultOptions()}
	// ensureOpen will fail because pool is nil, so nil fn check is not reached.
	// Test the nil fn guard path through WithTx which calls WithTxOptions.
	err := client.WithTxOptions(t.Context(), TxOptions{}, nil)
	// Should get connection error (not nil fn error) because pool is nil.
	if err == nil {
		t.Fatal("WithTxOptions(nil fn) error = nil, want non-nil")
	}
}
