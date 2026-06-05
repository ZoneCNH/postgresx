package testkit

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

var _ postgresx.Queryer = (*FakeQueryer)(nil)
var _ postgresx.Row = (*FakeRow)(nil)
var _ postgresx.Rows = (*FakeRows)(nil)
var _ postgresx.CommandTag = FakeCommandTag{}

// QueryCall captures SQL and args observed by FakeQueryer.
type QueryCall struct {
	Operation string
	SQL       string
	Args      []any
}

// FakeCommandTag is a deterministic postgresx.CommandTag for unit tests.
type FakeCommandTag struct {
	Rows int64
}

// RowsAffected returns the configured affected row count.
func (t FakeCommandTag) RowsAffected() int64 {
	return t.Rows
}

// FakeQueryer is an in-memory postgresx.Queryer for downstream unit tests.
type FakeQueryer struct {
	Calls []QueryCall

	ExecTag FakeCommandTag
	ExecErr error

	QueryRows *FakeRows
	QueryErr  error

	Row postgresx.Row
}

// Exec records the call and returns the configured command tag or error.
func (q *FakeQueryer) Exec(_ context.Context, sql string, args ...any) (postgresx.CommandTag, error) {
	q.record("exec", sql, args)
	if q.ExecErr != nil {
		return nil, q.ExecErr
	}
	return q.ExecTag, nil
}

// Query records the call and returns the configured rows or error.
func (q *FakeQueryer) Query(_ context.Context, sql string, args ...any) (postgresx.Rows, error) {
	q.record("query", sql, args)
	if q.QueryErr != nil {
		return nil, q.QueryErr
	}
	if q.QueryRows == nil {
		return &FakeRows{}, nil
	}
	return q.QueryRows, nil
}

// QueryRow records the call and returns the configured row.
func (q *FakeQueryer) QueryRow(_ context.Context, sql string, args ...any) postgresx.Row {
	q.record("query_row", sql, args)
	if q.Row == nil {
		return &FakeRow{}
	}
	return q.Row
}

func (q *FakeQueryer) record(operation, sql string, args []any) {
	if q == nil {
		return
	}
	q.Calls = append(q.Calls, QueryCall{Operation: operation, SQL: sql, Args: append([]any(nil), args...)})
}

// FakeRow is a deterministic postgresx.Row backed by Values.
type FakeRow struct {
	Values []any
	Err    error
}

// Scan copies Values into dest or returns Err when configured.
func (r *FakeRow) Scan(dest ...any) error {
	if r == nil {
		return errors.New("fake row is nil")
	}
	if r.Err != nil {
		return r.Err
	}
	return scanValues(r.Values, dest)
}

// FakeRows is a deterministic postgresx.Rows backed by row value slices.
type FakeRows struct {
	Rows [][]any
	ErrValue error
	Closed bool

	index int
}

// Close records that the rows were closed.
func (r *FakeRows) Close() {
	if r != nil {
		r.Closed = true
	}
}

// Err returns the configured rows error.
func (r *FakeRows) Err() error {
	if r == nil {
		return errors.New("fake rows is nil")
	}
	return r.ErrValue
}

// Next advances to the next row.
func (r *FakeRows) Next() bool {
	if r == nil || r.index >= len(r.Rows) {
		return false
	}
	r.index++
	return true
}

// Scan copies the current row into dest.
func (r *FakeRows) Scan(dest ...any) error {
	if r == nil {
		return errors.New("fake rows is nil")
	}
	if r.index == 0 || r.index > len(r.Rows) {
		return errors.New("fake rows scan without current row")
	}
	return scanValues(r.Rows[r.index-1], dest)
}

func scanValues(values []any, dest []any) error {
	if len(values) != len(dest) {
		return fmt.Errorf("fake scan value count %d does not match dest count %d", len(values), len(dest))
	}
	for i, value := range values {
		target := reflect.ValueOf(dest[i])
		if target.Kind() != reflect.Ptr || target.IsNil() {
			return fmt.Errorf("fake scan dest %d must be a non-nil pointer", i)
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
		return fmt.Errorf("fake scan value %d type %T cannot assign to %s", i, value, target.Elem().Type())
	}
	return nil
}
