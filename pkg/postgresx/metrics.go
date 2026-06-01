package postgresx

import "context"

const (
	metricQueryTotal    = "postgresx_query_total"
	metricQueryDuration = "postgresx_query_duration_seconds"
	metricTxTotal       = "postgresx_tx_total"
	metricTxDuration    = "postgresx_tx_duration_seconds"
	metricHealthTotal   = "postgresx_health_total"
	metricHealthLatency = "postgresx_health_latency_seconds"
	metricPoolConns     = "postgresx_pool_connections"
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
