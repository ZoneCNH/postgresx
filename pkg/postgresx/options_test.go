package postgresx

import (
	"context"
	"testing"
	"time"
)

type fixedClock struct {
	t time.Time
}

func (c *fixedClock) Now() time.Time { return c.t }

type stubLogger struct{}

func (stubLogger) Debug(context.Context, string, ...Field) {}
func (stubLogger) Info(context.Context, string, ...Field)  {}
func (stubLogger) Warn(context.Context, string, ...Field)  {}
func (stubLogger) Error(context.Context, string, ...Field) {}

type stubMetrics struct{}

func (stubMetrics) IncCounter(string, map[string]string)                {}
func (stubMetrics) ObserveHistogram(string, float64, map[string]string) {}
func (stubMetrics) SetGauge(string, float64, map[string]string)         {}

func TestDefaultOptionsHasNoopDefaults(t *testing.T) {
	opts := defaultOptions()
	if _, ok := opts.logger.(noopLogger); !ok {
		t.Fatalf("logger = %T, want noopLogger", opts.logger)
	}
	if _, ok := opts.metrics.(noopMetrics); !ok {
		t.Fatalf("metrics = %T, want noopMetrics", opts.metrics)
	}
	if opts.clock == nil {
		t.Fatal("clock is nil")
	}
}

func TestApplyOptionsEmpty(t *testing.T) {
	opts := applyOptions(nil)
	if _, ok := opts.logger.(noopLogger); !ok {
		t.Fatalf("logger = %T, want noopLogger", opts.logger)
	}
}

func TestApplyOptionsSkipsNilOptions(t *testing.T) {
	opts := applyOptions([]Option{nil, nil})
	if _, ok := opts.logger.(noopLogger); !ok {
		t.Fatalf("logger = %T, want noopLogger with nil options", opts.logger)
	}
}

func TestApplyOptionsAppliesAllOptions(t *testing.T) {
	logger := stubLogger{}
	metrics := stubMetrics{}
	clock := &fixedClock{t: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}

	opts := applyOptions([]Option{
		WithLogger(logger),
		WithMetrics(metrics),
		WithClock(clock),
	})

	if opts.logger != logger {
		t.Fatal("logger not applied")
	}
	if opts.metrics != metrics {
		t.Fatal("metrics not applied")
	}
	if opts.clock != clock {
		t.Fatal("clock not applied")
	}
}

func TestWithLoggerNilKeepsDefault(t *testing.T) {
	opts := applyOptions([]Option{WithLogger(nil)})
	if _, ok := opts.logger.(noopLogger); !ok {
		t.Fatalf("logger = %T, want noopLogger after nil logger option", opts.logger)
	}
}

func TestWithMetricsNilKeepsDefault(t *testing.T) {
	opts := applyOptions([]Option{WithMetrics(nil)})
	if _, ok := opts.metrics.(noopMetrics); !ok {
		t.Fatalf("metrics = %T, want noopMetrics after nil metrics option", opts.metrics)
	}
}

func TestWithClockNilKeepsDefault(t *testing.T) {
	opts := applyOptions([]Option{WithClock(nil)})
	if opts.clock == nil {
		t.Fatal("clock is nil after WithClock(nil), want non-nil clock")
	}
}
