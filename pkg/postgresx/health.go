package postgresx

import (
	"context"
	"strconv"
	"time"
)

// HealthStatusEnum is the health state of a component.
type HealthStatusEnum string

const (
	HealthHealthy   HealthStatusEnum = "healthy"
	HealthDegraded  HealthStatusEnum = "degraded"
	HealthUnhealthy HealthStatusEnum = "unhealthy"
)

// HealthChecker is implemented by components that can report health.
type HealthChecker interface {
	Name() string
	Check(ctx context.Context) HealthStatus
}

// HealthStatus carries a health snapshot for a component.
type HealthStatus struct {
	Name      string
	Status    HealthStatusEnum
	Message   string
	CheckedAt time.Time
	LatencyMs int64
	Metadata  map[string]string
}

var _ HealthChecker = (*Client)(nil)

// Name returns the stable health component name.
func (c *Client) Name() string {
	return "postgresx"
}

// Check returns a health status without exposing secrets.
func (c *Client) Check(ctx context.Context) HealthStatus {
	now := time.Now()
	if c != nil && c.opts.clock != nil {
		now = c.opts.clock.Now()
	}
	start := time.Now()
	status := HealthHealthy
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
		if IsKind(err, ErrorKindTimeout) || IsKind(err, ErrorKindConnection) {
			status = HealthDegraded
		} else {
			status = HealthUnhealthy
		}
		if c == nil || c.closed.Load() {
			status = HealthUnhealthy
		}
	}

	result := NewHealthStatus(c.Name(), status, message, now, time.Since(start).Milliseconds())
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
func (c *Client) HealthCheck(ctx context.Context) HealthStatus {
	return c.Check(ctx)
}

// NewHealthStatus constructs a health status snapshot.
func NewHealthStatus(name string, status HealthStatusEnum, message string, timestamp time.Time, latencyMs int64) HealthStatus {
	return HealthStatus{
		Name:      name,
		Status:    status,
		Message:   message,
		CheckedAt: timestamp,
		LatencyMs: latencyMs,
		Metadata:  make(map[string]string),
	}
}

// WithMetadata adds a key-value pair to the health status metadata.
func (h HealthStatus) WithMetadata(key, value string) HealthStatus {
	h.Metadata[key] = value
	return h
}

// IsHealthy reports whether the component is healthy.
func (h HealthStatus) IsHealthy() bool {
	return h.Status == HealthHealthy
}
