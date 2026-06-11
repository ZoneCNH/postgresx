// Package postgresxtest provides reusable contract-test adapters for postgresx.
package postgresxtest

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

var _ postgresx.Queryer = (*QueryAdapter)(nil)
var _ postgresx.Tx = (*QueryAdapter)(nil)
var _ postgresx.Row = (*Row)(nil)
var _ postgresx.Rows = (*Rows)(nil)
var _ postgresx.CommandTag = CommandTag{}

// QueryCall captures a Queryer operation for SQL contract assertions.
type QueryCall struct {
	Operation string
	SQL       string
	Args      []any
}

// CommandTag is a deterministic postgresx.CommandTag.
type CommandTag struct {
	Rows int64
}

// RowsAffected returns the configured affected-row count.
func (t CommandTag) RowsAffected() int64 { return t.Rows }

// QueryAdapter is an in-memory postgresx.Queryer/Tx for contract tests.
type QueryAdapter struct {
	Calls []QueryCall

	ExecTag CommandTag
	ExecErr error

	QueryRows *Rows
	QueryErr  error

	Row postgresx.Row
}

// Exec records the SQL command and returns the configured command tag or error.
func (q *QueryAdapter) Exec(_ context.Context, sql string, args ...any) (postgresx.CommandTag, error) {
	q.record("exec", sql, args)
	if q.ExecErr != nil {
		return nil, q.ExecErr
	}
	return q.ExecTag, nil
}

// Query records the SQL query and returns configured rows or error.
func (q *QueryAdapter) Query(_ context.Context, sql string, args ...any) (postgresx.Rows, error) {
	q.record("query", sql, args)
	if q.QueryErr != nil {
		return nil, q.QueryErr
	}
	if q.QueryRows == nil {
		return &Rows{}, nil
	}
	return q.QueryRows, nil
}

// QueryRow records the SQL query and returns the configured row.
func (q *QueryAdapter) QueryRow(_ context.Context, sql string, args ...any) postgresx.Row {
	q.record("query_row", sql, args)
	if q.Row == nil {
		return &Row{}
	}
	return q.Row
}

func (q *QueryAdapter) record(operation, sql string, args []any) {
	if q == nil {
		return
	}
	q.Calls = append(q.Calls, QueryCall{Operation: operation, SQL: sql, Args: append([]any(nil), args...)})
}

// Row is a deterministic postgresx.Row backed by Values.
type Row struct {
	Values []any
	Err    error
}

// Scan copies Values into dest or returns Err when configured.
func (r *Row) Scan(dest ...any) error {
	if r == nil {
		return errors.New("postgresxtest row is nil")
	}
	if r.Err != nil {
		return r.Err
	}
	return scanValues(r.Values, dest)
}

// Rows is a deterministic postgresx.Rows backed by row value slices.
type Rows struct {
	Rows     [][]any
	ErrValue error
	Closed   bool

	index int
}

// Close records that the rows were closed.
func (r *Rows) Close() {
	if r != nil {
		r.Closed = true
	}
}

// Err returns the configured rows error.
func (r *Rows) Err() error {
	if r == nil {
		return errors.New("postgresxtest rows is nil")
	}
	return r.ErrValue
}

// Next advances to the next row.
func (r *Rows) Next() bool {
	if r == nil || r.index >= len(r.Rows) {
		return false
	}
	r.index++
	return true
}

// Scan copies the current row into dest.
func (r *Rows) Scan(dest ...any) error {
	if r == nil {
		return errors.New("postgresxtest rows is nil")
	}
	if r.index == 0 || r.index > len(r.Rows) {
		return errors.New("postgresxtest rows scan without current row")
	}
	return scanValues(r.Rows[r.index-1], dest)
}

func scanValues(values []any, dest []any) error {
	if len(values) != len(dest) {
		return fmt.Errorf("postgresxtest scan value count %d does not match dest count %d", len(values), len(dest))
	}
	for i, value := range values {
		target := reflect.ValueOf(dest[i])
		if target.Kind() != reflect.Ptr || target.IsNil() {
			return fmt.Errorf("postgresxtest scan dest %d must be a non-nil pointer", i)
		}
		source := reflect.ValueOf(value)
		if !source.IsValid() {
			target.Elem().Set(reflect.Zero(target.Elem().Type()))
			continue
		}
		if source.Type().AssignableTo(target.Elem().Type()) {
			target.Elem().Set(source)
			continue
		}
		if source.Type().ConvertibleTo(target.Elem().Type()) {
			target.Elem().Set(source.Convert(target.Elem().Type()))
			continue
		}
		return fmt.Errorf("postgresxtest scan value %d type %T cannot assign to %s", i, value, target.Elem().Type())
	}
	return nil
}
