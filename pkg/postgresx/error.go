package postgresx

import (
	"errors"
	"fmt"
)

// ErrorKind classifies an error without coupling callers to driver types.
type ErrorKind string

const (
	ErrorKindConfig       ErrorKind = "config"
	ErrorKindCanceled     ErrorKind = "canceled"
	ErrorKindTimeout      ErrorKind = "timeout"
	ErrorKindNotFound     ErrorKind = "not_found"
	ErrorKindAuth         ErrorKind = "auth"
	ErrorKindConnection   ErrorKind = "connection"
	ErrorKindInternal     ErrorKind = "internal"
	ErrorKindValidation   ErrorKind = "validation"
	ErrorKindConflict     ErrorKind = "conflict"
	ErrorKindAlreadyExist ErrorKind = "already_exists"
	ErrorKindUnavailable  ErrorKind = "unavailable"
	ErrorKindRateLimit    ErrorKind = "rate_limit"
)

// Error is a classified error with an optional cause.
type Error struct {
	Kind      ErrorKind
	Op        string
	Message   string
	Cause     error
	Retryable bool
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// Unwrap returns the wrapped cause for errors.Is / errors.As.
func (e *Error) Unwrap() error { return e.Cause }

// WithRetryable marks the error as safe to retry.
func (e *Error) WithRetryable(v bool) *Error {
	e.Retryable = v
	return e
}

// NewError creates a new Error without a cause.
func NewError(kind ErrorKind, op, message string) *Error {
	return &Error{Kind: kind, Op: op, Message: message}
}

// WrapError wraps an existing error with classification.
func WrapError(kind ErrorKind, op, message string, cause error) *Error {
	return &Error{Kind: kind, Op: op, Message: message, Cause: cause}
}

// AsError unwraps an Error from the error chain.
func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// IsKind checks whether the error chain contains an Error of the given kind.
func IsKind(err error, kind ErrorKind) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind == kind
	}
	return false
}
