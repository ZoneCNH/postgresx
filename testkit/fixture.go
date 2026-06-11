package testkit

import (
	"context"

	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

// Fixture contains a live PostgreSQL client for integration tests.
type Fixture struct {
	client *postgresx.Client
}

// Client returns the postgresx client owned by the fixture.
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
