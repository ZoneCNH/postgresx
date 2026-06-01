package testkit

import (
	"context"
	"testing"
)

// RequireReady fails the test when the fixture database is not ready.
func RequireReady(ctx context.Context, t testing.TB, fixture *Fixture) {
	t.Helper()
	if fixture == nil || fixture.Client() == nil {
		t.Fatal("fixture is not open")
	}
	status := fixture.Client().Check(ctx)
	if !status.IsHealthy() {
		t.Fatalf("database is not ready: %s", status.Message)
	}
}
