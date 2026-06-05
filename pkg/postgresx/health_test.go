package postgresx

import (
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
)

func TestNilClientCheckReturnsUnhealthyStatus(t *testing.T) {
	var client *Client

	status := client.Check(t.Context())
	if status.Name != "postgresx" {
		t.Fatalf("Name = %q, want postgresx", status.Name)
	}
	if status.Status != foundationx.HealthUnhealthy {
		t.Fatalf("Status = %q, want unhealthy", status.Status)
	}
	if status.Message == "" {
		t.Fatal("Message is empty")
	}
	if len(status.Metadata) != 0 {
		t.Fatalf("Metadata = %+v, want empty metadata for nil client", status.Metadata)
	}
}
