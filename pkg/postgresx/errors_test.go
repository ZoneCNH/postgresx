package postgresx

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapErrorNormalizesKnownFailures(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		kind      ErrorKind
		retryable bool
	}{
		{
			name:      "context canceled",
			err:       context.Canceled,
			kind:      ErrorKindCanceled,
			retryable: false,
		},
		{
			name:      "context deadline",
			err:       context.DeadlineExceeded,
			kind:      ErrorKindTimeout,
			retryable: true,
		},
		{
			name:      "no rows",
			err:       pgx.ErrNoRows,
			kind:      ErrorKindNotFound,
			retryable: false,
		},
		{
			name:      "auth",
			err:       &pgconn.PgError{Code: "28P01"},
			kind:      ErrorKindAuth,
			retryable: false,
		},
		{
			name:      "syntax validation",
			err:       &pgconn.PgError{Code: "42601"},
			kind:      ErrorKindValidation,
			retryable: false,
		},
		{
			name:      "undefined table not found",
			err:       &pgconn.PgError{Code: "42P01"},
			kind:      ErrorKindNotFound,
			retryable: false,
		},
		{
			name:      "unique already exists",
			err:       &pgconn.PgError{Code: "23505"},
			kind:      ErrorKindAlreadyExist,
			retryable: false,
		},
		{
			name:      "foreign key conflict",
			err:       &pgconn.PgError{Code: "23503"},
			kind:      ErrorKindConflict,
			retryable: false,
		},
		{
			name:      "not null validation",
			err:       &pgconn.PgError{Code: "23502"},
			kind:      ErrorKindValidation,
			retryable: false,
		},
		{
			name:      "check validation",
			err:       &pgconn.PgError{Code: "23514"},
			kind:      ErrorKindValidation,
			retryable: false,
		},
		{
			name:      "serialization retry",
			err:       &pgconn.PgError{Code: "40001"},
			kind:      ErrorKindConflict,
			retryable: true,
		},
		{
			name:      "deadlock retry",
			err:       &pgconn.PgError{Code: "40P01"},
			kind:      ErrorKindConflict,
			retryable: true,
		},
		{
			name:      "lock not available retry",
			err:       &pgconn.PgError{Code: "55P03"},
			kind:      ErrorKindConflict,
			retryable: true,
		},
		{
			name:      "query canceled retry",
			err:       &pgconn.PgError{Code: "57014"},
			kind:      ErrorKindTimeout,
			retryable: true,
		},
		{
			name:      "connection class",
			err:       &pgconn.PgError{Code: "08006"},
			kind:      ErrorKindConnection,
			retryable: true,
		},
		{
			name:      "resource class",
			err:       &pgconn.PgError{Code: "53300"},
			kind:      ErrorKindUnavailable,
			retryable: true,
		},
		{
			name:      "admin shutdown class",
			err:       &pgconn.PgError{Code: "57P01"},
			kind:      ErrorKindUnavailable,
			retryable: true,
		},
		{
			name:      "unknown",
			err:       errors.New("driver failure"),
			kind:      ErrorKindInternal,
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapError("postgresx.test", tt.err)
			if !IsKind(err, tt.kind) {
				t.Fatalf("MapError() = %v, want kind %s", err, tt.kind)
			}
			if got := IsRetryable(err); got != tt.retryable {
				t.Fatalf("IsRetryable() = %v, want %v", got, tt.retryable)
			}
			if !errors.Is(err, tt.err) {
				t.Fatalf("MapError() does not unwrap original error %v", tt.err)
			}
		})
	}
}

func TestMapErrorNormalizesHeuristicDriverFailures(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		kind      ErrorKind
		retryable bool
	}{
		{
			name:      "password auth text",
			err:       errors.New("password authentication failed for user postgres"),
			kind:      ErrorKindAuth,
			retryable: false,
		},
		{
			name:      "connection text",
			err:       errors.New("dial tcp: connection refused"),
			kind:      ErrorKindConnection,
			retryable: true,
		},
		{
			name:      "wrapped context",
			err:       errors.Join(errors.New("driver wrapper"), context.DeadlineExceeded),
			kind:      ErrorKindTimeout,
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapError("postgresx.test", tt.err)
			if !IsKind(err, tt.kind) {
				t.Fatalf("MapError() = %v, want kind %s", err, tt.kind)
			}
			if got := IsRetryable(err); got != tt.retryable {
				t.Fatalf("IsRetryable() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestMapErrorNilAndRetryableNil(t *testing.T) {
	if err := MapError("postgresx.test", nil); err != nil {
		t.Fatalf("MapError(nil) = %v, want nil", err)
	}
	if IsRetryable(nil) {
		t.Fatal("IsRetryable(nil) = true, want false")
	}
}

func TestIsRetryableNonFoundationError(t *testing.T) {
	if IsRetryable(errors.New("plain error")) {
		t.Fatal("IsRetryable(plain error) = true, want false")
	}
}

func TestMapErrorPreservesFoundationError(t *testing.T) {
	original := NewError(ErrorKindValidation, "op", "bad input")
	mapped := MapError("postgresx.test", original)

	if !errors.Is(mapped, original) {
		t.Fatalf("MapError() = %v, want original %v", mapped, original)
	}
}
