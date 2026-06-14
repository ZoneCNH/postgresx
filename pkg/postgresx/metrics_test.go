package postgresx

import (
	"context"
	"testing"
)

func TestNoopLoggerMethodsDoNotPanic(t *testing.T) {
	logger := noopLogger{}
	ctx := context.Background()

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message", Field{Key: "k", Value: "v"})
	logger.Error(ctx, "error message", Field{Key: "err", Value: "something"})
}

func TestNoopMetricsMethodsDoNotPanic(t *testing.T) {
	metrics := noopMetrics{}

	metrics.IncCounter("test.counter", map[string]string{"label": "value"})
	metrics.ObserveHistogram("test.histogram", 1.5, map[string]string{"label": "value"})
	metrics.SetGauge("test.gauge", 42.0, map[string]string{"label": "value"})
}

func TestNoopLoggerSatisfiesLoggerInterface(t *testing.T) {
	var _ Logger = noopLogger{}
}

func TestNoopMetricsSatisfiesMetricsInterface(t *testing.T) {
	var _ Metrics = noopMetrics{}
}
