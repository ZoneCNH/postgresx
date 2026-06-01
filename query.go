package postgresx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// CommandTag is the stable command-result interface exposed by postgresx.
type CommandTag interface {
	RowsAffected() int64
}

// Row is the stable single-row scanner interface exposed by postgresx.
type Row interface {
	Scan(dest ...any) error
}

// Rows is the stable multi-row scanner interface exposed by postgresx.
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

// Queryer is the minimal query interface implemented by Client and Tx.
type Queryer interface {
	Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
}

type commandTag struct {
	tag pgconn.CommandTag
}

func (t commandTag) RowsAffected() int64 {
	return t.tag.RowsAffected()
}

type rowAdapter struct {
	row pgx.Row
	op  string
}

func (r rowAdapter) Scan(dest ...any) error {
	return MapError(r.op, r.row.Scan(dest...))
}

type errorRow struct {
	err error
}

func (r errorRow) Scan(dest ...any) error {
	return r.err
}

type rowsAdapter struct {
	rows pgx.Rows
	op   string
}

func (r rowsAdapter) Close() {
	r.rows.Close()
}

func (r rowsAdapter) Err() error {
	return MapError(r.op, r.rows.Err())
}

func (r rowsAdapter) Next() bool {
	return r.rows.Next()
}

func (r rowsAdapter) Scan(dest ...any) error {
	return MapError(r.op, r.rows.Scan(dest...))
}
