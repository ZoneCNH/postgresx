package postgresx

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client owns a PostgreSQL connection pool.
type Client struct {
	pool   *pgxpool.Pool
	cfg    Config
	opts   options
	closed atomic.Bool
}

// New validates cfg, opens a pgx pool, and verifies connectivity.
func New(ctx context.Context, cfg Config, opts ...Option) (*Client, error) {
	const op = "postgresx.New"
	cfg = cfg.withDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	resolved := applyOptions(opts)
	connectCtx := ctx
	cancel := func() {}
	if cfg.ConnectTimeout > 0 {
		connectCtx, cancel = context.WithTimeout(ctx, cfg.ConnectTimeout)
	}
	defer cancel()

	pgxCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, MapError(op, err)
	}
	pgxCfg.MaxConns = cfg.MaxOpenConns
	pgxCfg.MinConns = cfg.MinIdleConns
	pgxCfg.MaxConnLifetime = cfg.MaxConnLifetime
	pgxCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	pgxCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout
	if cfg.ApplicationName != "" {
		if pgxCfg.ConnConfig.RuntimeParams == nil {
			pgxCfg.ConnConfig.RuntimeParams = map[string]string{}
		}
		pgxCfg.ConnConfig.RuntimeParams["application_name"] = cfg.ApplicationName
	}

	pool, err := pgxpool.NewWithConfig(connectCtx, pgxCfg)
	if err != nil {
		return nil, MapError(op, err)
	}
	client := &Client{pool: pool, cfg: cfg, opts: resolved}
	healthCtx := ctx
	healthCancel := func() {}
	if cfg.HealthTimeout > 0 {
		healthCtx, healthCancel = context.WithTimeout(ctx, cfg.HealthTimeout)
	}
	defer healthCancel()
	if err := client.Ping(healthCtx); err != nil {
		pool.Close()
		client.closed.Store(true)
		return nil, MapError(op, err)
	}
	return client, nil
}

// Open is an alias for New.
func Open(ctx context.Context, cfg Config, opts ...Option) (*Client, error) {
	return New(ctx, cfg, opts...)
}

// Ping verifies pool connectivity.
func (c *Client) Ping(ctx context.Context) error {
	if err := c.ensureOpen("postgresx.Client.Ping"); err != nil {
		return err
	}
	start := time.Now()
	err := MapError("postgresx.Client.Ping", c.pool.Ping(ctx))
	labels := map[string]string{"outcome": "success"}
	if err != nil {
		labels["outcome"] = "error"
	}
	c.opts.metrics.IncCounter(metricQueryTotal, map[string]string{"operation": "ping", "outcome": labels["outcome"]})
	c.opts.metrics.ObserveHistogram(metricQueryDuration, time.Since(start).Seconds(), map[string]string{"operation": "ping", "outcome": labels["outcome"]})
	return err
}

// Close closes the pool. It is idempotent.
func (c *Client) Close(ctx context.Context) error {
	_ = ctx
	if c == nil || c.pool == nil {
		return nil
	}
	if c.closed.CompareAndSwap(false, true) {
		c.pool.Close()
	}
	return nil
}

// Stats returns driver-neutral pool statistics.
func (c *Client) Stats() PoolStats {
	if c == nil || c.pool == nil {
		return PoolStats{}
	}
	stat := c.pool.Stat()
	stats := PoolStats{
		TotalConns:        stat.TotalConns(),
		IdleConns:         stat.IdleConns(),
		AcquiredConns:     stat.AcquiredConns(),
		ConstructingConns: stat.ConstructingConns(),
		MaxConns:          stat.MaxConns(),
	}
	c.opts.metrics.SetGauge(metricPoolConns, float64(stats.TotalConns), map[string]string{"state": "total"})
	c.opts.metrics.SetGauge(metricPoolConns, float64(stats.IdleConns), map[string]string{"state": "idle"})
	c.opts.metrics.SetGauge(metricPoolConns, float64(stats.AcquiredConns), map[string]string{"state": "acquired"})
	c.opts.metrics.SetGauge(metricPoolConns, float64(stats.ConstructingConns), map[string]string{"state": "constructing"})
	c.opts.metrics.SetGauge(metricPoolConns, float64(stats.MaxConns), map[string]string{"state": "max"})
	return stats
}

// Queryer returns the client as the minimal query interface.
func (c *Client) Queryer() Queryer {
	return c
}

// Exec executes a SQL command.
func (c *Client) Exec(ctx context.Context, sql string, args ...any) (CommandTag, error) {
	const op = "postgresx.Client.Exec"
	if err := c.ensureOpen(op); err != nil {
		return nil, err
	}
	start := time.Now()
	tag, err := c.pool.Exec(ctx, sql, args...)
	mapped := MapError(op, err)
	recordQueryMetrics(c.opts.metrics, "exec", start, mapped)
	if mapped != nil {
		return nil, mapped
	}
	return commandTag{tag: tag}, nil
}

// Query executes a SQL query.
func (c *Client) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	const op = "postgresx.Client.Query"
	if err := c.ensureOpen(op); err != nil {
		return nil, err
	}
	start := time.Now()
	rows, err := c.pool.Query(ctx, sql, args...)
	mapped := MapError(op, err)
	recordQueryMetrics(c.opts.metrics, "query", start, mapped)
	if mapped != nil {
		return nil, mapped
	}
	return rowsAdapter{rows: rows, op: "postgresx.Client.Query.Rows"}, nil
}

// QueryRow executes a SQL query that is expected to return at most one row.
func (c *Client) QueryRow(ctx context.Context, sql string, args ...any) Row {
	const op = "postgresx.Client.QueryRow"
	if err := c.ensureOpen(op); err != nil {
		return errorRow{err: err}
	}
	return &metricRow{
		row:       rowAdapter{row: c.pool.QueryRow(ctx, sql, args...), op: "postgresx.Client.QueryRow"},
		metrics:   c.opts.metrics,
		operation: "query_row",
		start:     time.Now(),
	}
}

type metricRow struct {
	row       Row
	metrics   Metrics
	operation string
	start     time.Time
	once      sync.Once
}

func (r *metricRow) Scan(dest ...any) error {
	err := r.row.Scan(dest...)
	r.once.Do(func() {
		recordQueryMetrics(r.metrics, r.operation, r.start, err)
	})
	return err
}

func recordQueryMetrics(metrics Metrics, operation string, start time.Time, err error) {
	outcome := "success"
	if err != nil {
		outcome = "error"
	}
	labels := map[string]string{"operation": operation, "outcome": outcome}
	metrics.IncCounter(metricQueryTotal, labels)
	metrics.ObserveHistogram(metricQueryDuration, time.Since(start).Seconds(), labels)
}

func (c *Client) ensureOpen(op string) error {
	if c == nil || c.pool == nil || c.closed.Load() {
		return NewError(ErrorKindConnection, op, "client is closed")
	}
	return nil
}
