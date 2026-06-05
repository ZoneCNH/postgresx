package postgresx

import (
	"context"
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
)

func TestClosedClientOperationsReturnConnectionError(t *testing.T) {
	ctx := t.Context()
	var client *Client

	if _, err := client.Exec(ctx, "select 1"); !foundationx.IsKind(err, foundationx.ErrorKindConnection) {
		t.Fatalf("Exec() error = %v, want connection error", err)
	}
	if _, err := client.Query(ctx, "select 1"); !foundationx.IsKind(err, foundationx.ErrorKindConnection) {
		t.Fatalf("Query() error = %v, want connection error", err)
	}
	if err := client.QueryRow(ctx, "select 1").Scan(new(int)); !foundationx.IsKind(err, foundationx.ErrorKindConnection) {
		t.Fatalf("QueryRow().Scan() error = %v, want connection error", err)
	}
	if err := client.WithTx(ctx, func(ctx context.Context, tx Tx) error { return nil }); !foundationx.IsKind(err, foundationx.ErrorKindConnection) {
		t.Fatalf("WithTx() error = %v, want connection error", err)
	}
	if err := client.Close(ctx); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}
