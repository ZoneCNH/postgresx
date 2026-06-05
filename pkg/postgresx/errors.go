package postgresx

import (
	"context"
	"errors"
	"strings"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MapError normalizes context and PostgreSQL driver errors to foundationx.Error.
func MapError(op string, err error) error {
	if err == nil {
		return nil
	}
	if _, ok := foundationx.AsFoundationError(err); ok {
		return err
	}
	switch {
	case errors.Is(err, context.Canceled):
		return foundationx.WrapError(foundationx.ErrorKindCanceled, op, "operation canceled", err)
	case errors.Is(err, context.DeadlineExceeded):
		return foundationx.WrapError(foundationx.ErrorKindTimeout, op, "operation timed out", err).WithRetryable(true)
	case errors.Is(err, pgx.ErrNoRows):
		return foundationx.WrapError(foundationx.ErrorKindNotFound, op, "row not found", err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return mapPgError(op, pgErr, err)
	}

	text := strings.ToLower(err.Error())
	switch {
	case strings.Contains(text, "password authentication failed"):
		return foundationx.WrapError(foundationx.ErrorKindAuth, op, "authentication failed", err)
	case strings.Contains(text, "connect") || strings.Contains(text, "connection"):
		return foundationx.WrapError(foundationx.ErrorKindConnection, op, "connection failed", err).WithRetryable(true)
	default:
		return foundationx.WrapError(foundationx.ErrorKindInternal, op, "postgres operation failed", err)
	}
}

// IsRetryable reports whether a normalized error may be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if fxErr, ok := foundationx.AsFoundationError(err); ok {
		return fxErr.Retryable
	}
	return false
}

func mapPgError(op string, pgErr *pgconn.PgError, cause error) error {
	switch pgErr.Code {
	case "28P01":
		return foundationx.WrapError(foundationx.ErrorKindAuth, op, "authentication failed", cause)
	case "42601":
		return foundationx.WrapError(foundationx.ErrorKindValidation, op, "syntax error", cause)
	case "23505":
		return foundationx.WrapError(foundationx.ErrorKindAlreadyExist, op, "unique constraint violation", cause)
	case "23503":
		return foundationx.WrapError(foundationx.ErrorKindConflict, op, "foreign key constraint violation", cause)
	case "23502":
		return foundationx.WrapError(foundationx.ErrorKindValidation, op, "not null constraint violation", cause)
	case "23514":
		return foundationx.WrapError(foundationx.ErrorKindValidation, op, "check constraint violation", cause)
	case "40001":
		return foundationx.WrapError(foundationx.ErrorKindConflict, op, "serialization failure", cause).WithRetryable(true)
	case "40P01":
		return foundationx.WrapError(foundationx.ErrorKindConflict, op, "deadlock detected", cause).WithRetryable(true)
	case "55P03":
		return foundationx.WrapError(foundationx.ErrorKindConflict, op, "lock not available", cause).WithRetryable(true)
	case "57014":
		return foundationx.WrapError(foundationx.ErrorKindTimeout, op, "query canceled", cause).WithRetryable(true)
	}
	if strings.HasPrefix(pgErr.Code, "08") {
		return foundationx.WrapError(foundationx.ErrorKindConnection, op, "connection failed", cause).WithRetryable(true)
	}
	if strings.HasPrefix(pgErr.Code, "53") {
		return foundationx.WrapError(foundationx.ErrorKindUnavailable, op, "postgres resources unavailable", cause).WithRetryable(true)
	}
	if strings.HasPrefix(pgErr.Code, "57") {
		return foundationx.WrapError(foundationx.ErrorKindUnavailable, op, "postgres unavailable", cause).WithRetryable(true)
	}
	return foundationx.WrapError(foundationx.ErrorKindInternal, op, "postgres operation failed", cause)
}
