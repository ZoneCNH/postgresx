package postgresx

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestClosedClientOperationsReturnConnectionError(t *testing.T) {
	ctx := t.Context()
	var client *Client

	if _, err := client.Exec(ctx, "select 1"); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("Exec() error = %v, want connection error", err)
	}
	if _, err := client.Query(ctx, "select 1"); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("Query() error = %v, want connection error", err)
	}
	if err := client.QueryRow(ctx, "select 1").Scan(new(int)); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("QueryRow().Scan() error = %v, want connection error", err)
	}
	if err := client.WithTx(ctx, func(ctx context.Context, tx Tx) error { return nil }); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("WithTx() error = %v, want connection error", err)
	}
	if err := client.Close(ctx); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}

func TestNonNilClosedClientOperationsReturnConnectionError(t *testing.T) {
	ctx := t.Context()
	client := &Client{opts: defaultOptions()}
	client.closed.Store(true)

	if err := client.Close(ctx); err != nil {
		t.Fatalf("first Close() error = %v, want nil", err)
	}
	if err := client.Close(ctx); err != nil {
		t.Fatalf("second Close() error = %v, want nil", err)
	}
	if err := client.Ping(ctx); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("Ping() error = %v, want connection error", err)
	}
	if _, err := client.Exec(ctx, "select 1"); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("Exec() error = %v, want connection error", err)
	}
	if _, err := client.Query(ctx, "select 1"); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("Query() error = %v, want connection error", err)
	}
	if err := client.QueryRow(ctx, "select 1").Scan(new(int)); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("QueryRow().Scan() error = %v, want connection error", err)
	}
	if err := client.WithTx(ctx, func(ctx context.Context, tx Tx) error { return nil }); !IsKind(err, ErrorKindConnection) {
		t.Fatalf("WithTx() error = %v, want connection error", err)
	}
	if stats := client.Stats(); stats != (PoolStats{}) {
		t.Fatalf("Stats() = %+v, want zero stats", stats)
	}
}

func TestMetricNamesMatchContract(t *testing.T) {
	tests := map[string]string{
		"query total":    metricQueryTotal,
		"query duration": metricQueryDuration,
		"tx total":       metricTxTotal,
		"tx duration":    metricTxDuration,
		"health total":   metricHealthTotal,
		"health latency": metricHealthLatency,
		"pool conns":     metricPoolConns,
	}

	for name, got := range tests {
		t.Run(name, func(t *testing.T) {
			if got == "" {
				t.Fatal("metric name is empty")
			}
			if !strings.HasPrefix(got, "postgresx.") {
				t.Fatalf("metric name = %q, want postgresx dotted namespace", got)
			}
		})
	}
}

func TestMetricRowRecordsOnceOnScan(t *testing.T) {
	metrics := &captureMetrics{}
	row := &metricRow{
		row:       staticRow{},
		metrics:   metrics,
		operation: "query_row",
		start:     time.Now(),
	}

	if err := row.Scan(new(int)); err != nil {
		t.Fatalf("first Scan() error = %v", err)
	}
	if err := row.Scan(new(int)); err != nil {
		t.Fatalf("second Scan() error = %v", err)
	}

	if got := metrics.counterCalls; got != 1 {
		t.Fatalf("counter calls = %d, want 1", got)
	}
	if got := metrics.histogramCalls; got != 1 {
		t.Fatalf("histogram calls = %d, want 1", got)
	}
	if metrics.lastCounterName != metricQueryTotal {
		t.Fatalf("counter metric = %q, want %q", metrics.lastCounterName, metricQueryTotal)
	}
	if metrics.lastCounterLabels["operation"] != "query_row" || metrics.lastCounterLabels["outcome"] != "success" {
		t.Fatalf("counter labels = %+v, want query_row success", metrics.lastCounterLabels)
	}
}

func TestRecordQueryMetricsErrorPath(t *testing.T) {
	metrics := &captureMetrics{}
	recordQueryMetrics(metrics, "exec", time.Now(), context.Canceled)

	if metrics.counterCalls != 1 {
		t.Fatalf("counter calls = %d, want 1", metrics.counterCalls)
	}
	if metrics.histogramCalls != 1 {
		t.Fatalf("histogram calls = %d, want 1", metrics.histogramCalls)
	}
	if metrics.lastCounterLabels["outcome"] != "error" {
		t.Fatalf("outcome label = %q, want error", metrics.lastCounterLabels["outcome"])
	}
	if metrics.lastCounterLabels["operation"] != "exec" {
		t.Fatalf("operation label = %q, want exec", metrics.lastCounterLabels["operation"])
	}
}

func TestEnsureOpenNilPool(t *testing.T) {
	client := &Client{opts: defaultOptions()}
	err := client.ensureOpen("test.EnsureOpen")
	if !IsKind(err, ErrorKindConnection) {
		t.Fatalf("ensureOpen() error = %v, want connection error", err)
	}
}

func TestQueryerReturnsClient(t *testing.T) {
	client := &Client{opts: defaultOptions()}
	q := client.Queryer()
	if q != client {
		t.Fatal("Queryer() did not return the client")
	}
}

func TestHealthCheckDelegatesToCheck(t *testing.T) {
	var client *Client
	status := client.HealthCheck(t.Context())
	if status.Name != "postgresx" {
		t.Fatalf("Name = %q, want postgresx", status.Name)
	}
	if status.Status != HealthUnhealthy {
		t.Fatalf("Status = %q, want unhealthy", status.Status)
	}
}

type staticRow struct{}

func (staticRow) Scan(dest ...any) error {
	if len(dest) == 1 {
		if value, ok := dest[0].(*int); ok {
			*value = 1
		}
	}
	return nil
}

type captureMetrics struct {
	counterCalls      int
	histogramCalls    int
	lastCounterName   string
	lastCounterLabels map[string]string
}

func (m *captureMetrics) IncCounter(name string, labels map[string]string) {
	m.counterCalls++
	m.lastCounterName = name
	m.lastCounterLabels = labels
}

func (m *captureMetrics) ObserveHistogram(string, float64, map[string]string) {
	m.histogramCalls++
}

func (m *captureMetrics) SetGauge(string, float64, map[string]string) {}
