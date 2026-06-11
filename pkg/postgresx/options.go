package postgresx

import (
	"context"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
)

// Option customizes Client behavior.
type Option func(*options)

type options struct {
	logger  Logger
	metrics Metrics
	clock   foundationx.Clock
}

// Field is a minimal structured log field.
type Field struct {
	Key   string
	Value any
}

// Logger is the minimal logging interface used by postgresx.
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
}

// Metrics is the minimal metrics interface used by postgresx.
type Metrics interface {
	IncCounter(name string, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
}

// WithLogger injects a logger. Passing nil keeps the noop logger.
func WithLogger(logger Logger) Option {
	return func(o *options) {
		if logger != nil {
			o.logger = logger
		}
	}
}

// WithMetrics injects a metrics recorder. Passing nil keeps noop metrics.
func WithMetrics(metrics Metrics) Option {
	return func(o *options) {
		if metrics != nil {
			o.metrics = metrics
		}
	}
}

// WithClock injects a clock for tests and deterministic health output.
func WithClock(clock foundationx.Clock) Option {
	return func(o *options) {
		if clock != nil {
			o.clock = clock
		}
	}
}

func defaultOptions() options {
	return options{
		logger:  noopLogger{},
		metrics: noopMetrics{},
		clock:   foundationx.NewRealClock(),
	}
}

func applyOptions(opts []Option) options {
	resolved := defaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&resolved)
		}
	}
	return resolved
}
