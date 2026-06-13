package postgresx

import (
	"context"
	"strconv"
	"time"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
)

var _ foundationx.HealthChecker = (*Client)(nil)

// Name returns the stable health component name.
func (c *Client) Name() string {
	return "postgresx"
}

// Check returns a foundationx health status without exposing secrets.
func (c *Client) Check(ctx context.Context) foundationx.HealthStatus {
	now := time.Now()
	if c != nil && c.opts.clock != nil {
		now = c.opts.clock.Now()
	}
	start := time.Now()
	status := foundationx.HealthHealthy
	message := "ok"

	checkCtx := ctx
	cancel := func() {}
	if c != nil && c.cfg.HealthTimeout > 0 {
		checkCtx, cancel = context.WithTimeout(ctx, c.cfg.HealthTimeout)
	}
	defer cancel()

	err := c.Ping(checkCtx)
	if err != nil {
		message = err.Error()
		if foundationx.IsKind(err, foundationx.ErrorKindTimeout) || foundationx.IsKind(err, foundationx.ErrorKindConnection) {
			status = foundationx.HealthDegraded
		} else {
			status = foundationx.HealthUnhealthy
		}
		if c == nil || c.closed.Load() {
			status = foundationx.HealthUnhealthy
		}
	}

	result := foundationx.NewHealthStatus(c.Name(), status, message, now, time.Since(start).Milliseconds())
	if c != nil {
		result = result.WithMetadata("host", c.cfg.Host)
		result = result.WithMetadata("database", c.cfg.Database)
		result = result.WithMetadata("application_name", c.cfg.ApplicationName)
		stats := c.Stats()
		result = result.WithMetadata("pool_total_conns", strconv.FormatInt(int64(stats.TotalConns), 10))
		result = result.WithMetadata("pool_idle_conns", strconv.FormatInt(int64(stats.IdleConns), 10))
		result = result.WithMetadata("pool_acquired_conns", strconv.FormatInt(int64(stats.AcquiredConns), 10))
		result = result.WithMetadata("pool_constructing_conns", strconv.FormatInt(int64(stats.ConstructingConns), 10))
		result = result.WithMetadata("pool_max_conns", strconv.FormatInt(int64(stats.MaxConns), 10))
	}

	labels := map[string]string{"status": string(result.Status)}
	if c != nil {
		c.opts.metrics.IncCounter(metricHealthTotal, labels)
		c.opts.metrics.ObserveHistogram(metricHealthLatency, time.Since(start).Seconds(), labels)
	}
	return result
}

// HealthCheck implements the shared HealthCheck interface shape.
func (c *Client) HealthCheck(ctx context.Context) foundationx.HealthStatus {
	return c.Check(ctx)
}
