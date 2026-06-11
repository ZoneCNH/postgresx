package postgresx

import (
	"context"
	"time"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/jackc/pgx/v5"
)

// Tx is the transaction query interface exposed by postgresx.
type Tx interface {
	Queryer
}

// TxFunc is executed inside a transaction.
type TxFunc func(ctx context.Context, tx Tx) error

// TxOptions contains driver-neutral transaction options.
type TxOptions struct {
	IsolationLevel string
	ReadOnly       bool
}

type txAdapter struct {
	tx pgx.Tx
}

func (t txAdapter) Exec(ctx context.Context, sql string, args ...any) (CommandTag, error) {
	tag, err := t.tx.Exec(ctx, sql, args...)
	if err := MapError("postgresx.Tx.Exec", err); err != nil {
		return nil, err
	}
	return commandTag{tag: tag}, nil
}

func (t txAdapter) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err := MapError("postgresx.Tx.Query", err); err != nil {
		return nil, err
	}
	return rowsAdapter{rows: rows, op: "postgresx.Tx.Query.Rows"}, nil
}

func (t txAdapter) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return rowAdapter{row: t.tx.QueryRow(ctx, sql, args...), op: "postgresx.Tx.QueryRow"}
}

// WithTx runs fn in a transaction with default options.
func (c *Client) WithTx(ctx context.Context, fn TxFunc) error {
	return c.WithTxOptions(ctx, TxOptions{}, fn)
}

// WithTxOptions runs fn in a transaction and commits only when fn returns nil.
func (c *Client) WithTxOptions(ctx context.Context, opts TxOptions, fn TxFunc) (err error) {
	const op = "postgresx.Client.WithTxOptions"
	if err := c.ensureOpen(op); err != nil {
		return err
	}
	if fn == nil {
		return foundationx.NewError(foundationx.ErrorKindValidation, op, "transaction function is required")
	}

	start := time.Now()
	outcome := "commit"
	defer func() {
		if err != nil {
			outcome = "rollback"
		}
		c.opts.metrics.IncCounter(metricTxTotal, map[string]string{"outcome": outcome})
		c.opts.metrics.ObserveHistogram(metricTxDuration, time.Since(start).Seconds(), map[string]string{"outcome": outcome})
	}()

	pgxOpts := pgx.TxOptions{
		IsoLevel:   pgx.TxIsoLevel(opts.IsolationLevel),
		AccessMode: pgx.ReadWrite,
	}
	if opts.ReadOnly {
		pgxOpts.AccessMode = pgx.ReadOnly
	}

	tx, err := c.pool.BeginTx(ctx, pgxOpts)
	if err != nil {
		return MapError("postgresx.Client.WithTxOptions.Begin", err)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			outcome = "rollback"
			_ = tx.Rollback(ctx)
			panic(recovered)
		}
	}()

	if err := fn(ctx, txAdapter{tx: tx}); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			c.opts.logger.Warn(ctx, "postgresx transaction rollback failed", Field{Key: "error", Value: MapError("postgresx.Client.WithTxOptions.Rollback", rollbackErr)})
		}
		return MapError("postgresx.Client.WithTxOptions.Function", err)
	}
	if err := ctx.Err(); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			c.opts.logger.Warn(ctx, "postgresx transaction rollback failed", Field{Key: "error", Value: MapError("postgresx.Client.WithTxOptions.Rollback", rollbackErr)})
		}
		return MapError("postgresx.Client.WithTxOptions.Context", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return MapError("postgresx.Client.WithTxOptions.Commit", err)
	}
	return nil
}
