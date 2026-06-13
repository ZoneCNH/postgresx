#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GO="${GO:-go}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

cd "$tmpdir"

cat > go.mod <<EOF
module postgresx-downstream-smoke

go 1.23

require github.com/ZoneCNH/postgresx v0.0.0

replace github.com/ZoneCNH/postgresx => $ROOT_DIR
EOF

mkdir -p smoke
cat > smoke/postgresx_downstream_test.go <<'EOF'
package smoke_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	foundationx "github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
	"github.com/ZoneCNH/postgresx/test/postgresxtest"
)

func TestConfigRedactionUsesPublicAPI(t *testing.T) {
	const secret = "downstream-secret-value"
	cfg := postgresx.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = 5432
	cfg.Database = "postgresx_downstream"
	cfg.User = "postgresx"
	cfg.Password = foundationx.NewSecretString(secret)
	cfg.ApplicationName = "postgresx-downstream-smoke"

	if cfg.DSN() == "" {
		t.Fatal("DSN should be available through the public config API")
	}

	surfaces := map[string]string{
		"redactedDSN": cfg.RedactedDSN(),
		"sanitize":    fmt.Sprint(cfg.Sanitize()),
	}
	for name, surface := range surfaces {
		if strings.Contains(surface, secret) {
			t.Fatalf("%s leaked password: %q", name, surface)
		}
	}
}

func TestQueryerContractWithPostgresxTestkit(t *testing.T) {
	var _ postgresx.Queryer = (*postgresxtest.QueryAdapter)(nil)

	adapter := postgresxtest.QueryAdapter{
		Row: &postgresxtest.Row{Values: []any{"alice"}},
	}

	var got string
	if err := adapter.QueryRow(context.Background(), "select name").Scan(&got); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if got != "alice" {
		t.Fatalf("Scan got %q, want alice", got)
	}
}

func TestOptionalLiveDSNFromEnvironment(t *testing.T) {
	dsn := os.Getenv("POSTGRESX_INTEGRATION_DSN")
	if dsn == "" {
		t.Skip("POSTGRESX_INTEGRATION_DSN not set")
	}

	cfg, err := postgresxtest.ConfigFromDSN(dsn, "postgresx-downstream-smoke")
	if err != nil {
		t.Fatalf("ConfigFromDSN failed: %v", err)
	}
	cfg.ConnectTimeout = 5 * time.Second
	cfg.HealthTimeout = 5 * time.Second

	client, err := postgresx.Open(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer client.Close(context.Background())

	if err := client.Ping(context.Background()); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}
EOF

GOWORK=off "$GO" test -mod=mod ./...

echo "downstream smoke passed"
