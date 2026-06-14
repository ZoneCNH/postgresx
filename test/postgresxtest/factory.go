package postgresxtest

import (
	"context"
	"testing"
	"time"

	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

// Fixture owns a live postgresx client for integration contract tests.
type Fixture struct {
	client *postgresx.Client
}

// Client returns the fixture client.
func (f *Fixture) Client() *postgresx.Client {
	if f == nil {
		return nil
	}
	return f.client
}

// Close releases the fixture client.
func (f *Fixture) Close() {
	if f == nil || f.client == nil {
		return
	}
	_ = f.client.Close(context.Background())
}

// Start opens a live PostgreSQL fixture or skips when no integration DSN is configured.
func Start(ctx context.Context, t testing.TB, applicationName string) *Fixture {
	t.Helper()
	dsn := IntegrationDSN()
	if dsn == "" {
		t.Skip("POSTGRESX_INTEGRATION_DSN or POSTGRES_TEST_DSN is not set")
	}
	if applicationName == "" {
		applicationName = "postgresx-l2-contract"
	}
	cfg, err := ConfigFromDSN(dsn, applicationName)
	if err != nil {
		t.Fatalf("parse postgres fixture config: %v", err)
	}
	client, err := openWithRetry(ctx, cfg, 15*time.Second)
	if err != nil {
		t.Fatalf("open postgres fixture: %v", err)
	}
	fixture := &Fixture{client: client}
	t.Cleanup(fixture.Close)
	return fixture
}

func openWithRetry(ctx context.Context, cfg postgresx.Config, timeout time.Duration) (*postgresx.Client, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for {
		client, err := postgresx.Open(ctx, cfg)
		if err == nil {
			return client, nil
		}
		lastErr = err
		if (!postgresx.IsKind(err, postgresx.ErrorKindTimeout) && !postgresx.IsKind(err, postgresx.ErrorKindConnection)) || time.Now().After(deadline) {
			return nil, lastErr
		}
		timer := time.NewTimer(250 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, lastErr
		case <-timer.C:
		}
	}
}
