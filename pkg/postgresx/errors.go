package postgresx

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MapError normalizes context and PostgreSQL driver errors to a classified Error.
func MapError(op string, err error) error {
	if err == nil {
		return nil
	}
	if _, ok := AsError(err); ok {
		return err
	}
	switch {
	case errors.Is(err, context.Canceled):
		return WrapError(ErrorKindCanceled, op, "operation canceled", err)
	case errors.Is(err, context.DeadlineExceeded):
		return WrapError(ErrorKindTimeout, op, "operation timed out", err).WithRetryable(true)
	case errors.Is(err, pgx.ErrNoRows):
		return WrapError(ErrorKindNotFound, op, "row not found", err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return mapPgError(op, pgErr, err)
	}

	text := strings.ToLower(err.Error())
	switch {
	case strings.Contains(text, "password authentication failed"):
		return WrapError(ErrorKindAuth, op, "authentication failed", err)
	case strings.Contains(text, "connect") || strings.Contains(text, "connection"):
		return WrapError(ErrorKindConnection, op, "connection failed", err).WithRetryable(true)
	default:
		return WrapError(ErrorKindInternal, op, "postgres operation failed", err)
	}
}

// IsRetryable reports whether a normalized error may be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if fxErr, ok := AsError(err); ok {
		return fxErr.Retryable
	}
	return false
}

func mapPgError(op string, pgErr *pgconn.PgError, cause error) error {
	switch pgErr.Code {
	case "28P01":
		return WrapError(ErrorKindAuth, op, "authentication failed", cause)
	case "42601":
		return WrapError(ErrorKindValidation, op, "syntax error", cause)
	case "42P01":
		return WrapError(ErrorKindNotFound, op, "relation not found", cause)
	case "23505":
		return WrapError(ErrorKindAlreadyExist, op, "unique constraint violation", cause)
	case "23503":
		return WrapError(ErrorKindConflict, op, "foreign key constraint violation", cause)
	case "23502":
		return WrapError(ErrorKindValidation, op, "not null constraint violation", cause)
	case "23514":
		return WrapError(ErrorKindValidation, op, "check constraint violation", cause)
	case "40001":
		return WrapError(ErrorKindConflict, op, "serialization failure", cause).WithRetryable(true)
	case "40P01":
		return WrapError(ErrorKindConflict, op, "deadlock detected", cause).WithRetryable(true)
	case "55P03":
		return WrapError(ErrorKindConflict, op, "lock not available", cause).WithRetryable(true)
	case "57014":
		return WrapError(ErrorKindTimeout, op, "query canceled", cause).WithRetryable(true)
	}
	if strings.HasPrefix(pgErr.Code, "08") {
		return WrapError(ErrorKindConnection, op, "connection failed", cause).WithRetryable(true)
	}
	if strings.HasPrefix(pgErr.Code, "53") {
		return WrapError(ErrorKindUnavailable, op, "postgres resources unavailable", cause).WithRetryable(true)
	}
	if strings.HasPrefix(pgErr.Code, "57") {
		return WrapError(ErrorKindUnavailable, op, "postgres unavailable", cause).WithRetryable(true)
	}
	return WrapError(ErrorKindInternal, op, "postgres operation failed", cause)
}
