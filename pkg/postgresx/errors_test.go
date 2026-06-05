package postgresx

import (
	"context"
	"errors"
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapErrorNormalizesKnownFailures(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		kind      foundationx.ErrorKind
		retryable bool
	}{
		{
			name:      "context canceled",
			err:       context.Canceled,
			kind:      foundationx.ErrorKindCanceled,
			retryable: false,
		},
		{
			name:      "context deadline",
			err:       context.DeadlineExceeded,
			kind:      foundationx.ErrorKindTimeout,
			retryable: true,
		},
		{
			name:      "no rows",
			err:       pgx.ErrNoRows,
			kind:      foundationx.ErrorKindNotFound,
			retryable: false,
		},
		{
			name:      "auth",
			err:       &pgconn.PgError{Code: "28P01"},
			kind:      foundationx.ErrorKindAuth,
			retryable: false,
		},
		{
			name:      "unique conflict",
			err:       &pgconn.PgError{Code: "23505"},
			kind:      foundationx.ErrorKindConflict,
			retryable: false,
		},
		{
			name:      "serialization retry",
			err:       &pgconn.PgError{Code: "40001"},
			kind:      foundationx.ErrorKindConflict,
			retryable: true,
		},
		{
			name:      "connection class",
			err:       &pgconn.PgError{Code: "08006"},
			kind:      foundationx.ErrorKindConnection,
			retryable: true,
		},
		{
			name:      "admin shutdown class",
			err:       &pgconn.PgError{Code: "57P01"},
			kind:      foundationx.ErrorKindUnavailable,
			retryable: true,
		},
		{
			name:      "unknown",
			err:       errors.New("driver failure"),
			kind:      foundationx.ErrorKindInternal,
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapError("postgresx.test", tt.err)
			if !foundationx.IsKind(err, tt.kind) {
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

func TestMapErrorPreservesFoundationError(t *testing.T) {
	original := foundationx.NewError(foundationx.ErrorKindValidation, "op", "bad input")
	mapped := MapError("postgresx.test", original)

	if !errors.Is(mapped, original) {
		t.Fatalf("MapError() = %v, want original %v", mapped, original)
	}
}
