package postgresx

import "context"

const (
	metricQueryTotal    = "postgresx.query.total"
	metricQueryDuration = "postgresx.query.duration_seconds"
	metricTxTotal       = "postgresx.tx.total"
	metricTxDuration    = "postgresx.tx.duration_seconds"
	metricHealthTotal   = "postgresx.health.total"
	metricHealthLatency = "postgresx.health.latency_seconds"
	metricPoolConns     = "postgresx.pool.connections"
)

type noopLogger struct{}

func (noopLogger) Debug(context.Context, string, ...Field) {}
func (noopLogger) Info(context.Context, string, ...Field)  {}
func (noopLogger) Warn(context.Context, string, ...Field)  {}
func (noopLogger) Error(context.Context, string, ...Field) {}

type noopMetrics struct{}

func (noopMetrics) IncCounter(string, map[string]string)                {}
func (noopMetrics) ObserveHistogram(string, float64, map[string]string) {}
func (noopMetrics) SetGauge(string, float64, map[string]string)         {}
