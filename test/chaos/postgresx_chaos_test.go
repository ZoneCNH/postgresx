package chaos_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ZoneCNH/postgresx/pkg/postgresx"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestChaosErrorMappingCoversRetryableTransientFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		input     error
		wantKind  postgresx.ErrorKind
		retryable bool
	}{
		{
			name:      "deadline timeout",
			input:     context.DeadlineExceeded,
			wantKind:  postgresx.ErrorKindTimeout,
			retryable: true,
		},
		{
			name:      "caller cancellation",
			input:     context.Canceled,
			wantKind:  postgresx.ErrorKindCanceled,
			retryable: false,
		},
		{
			name:      "resource unavailable",
			input:     &pgconn.PgError{Code: "53300", Message: "too many connections"},
			wantKind:  postgresx.ErrorKindUnavailable,
			retryable: true,
		},
		{
			name:      "connection exception",
			input:     &pgconn.PgError{Code: "08006", Message: "connection failure"},
			wantKind:  postgresx.ErrorKindConnection,
			retryable: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := postgresx.MapError("chaos."+tc.name, tc.input)
			if !postgresx.IsKind(err, tc.wantKind) {
				t.Fatalf("MapError kind mismatch: want %v, err=%v", tc.wantKind, err)
			}
			if errors.Is(tc.input, context.DeadlineExceeded) && !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("MapError should preserve deadline cause: %v", err)
			}
			if got := postgresx.IsRetryable(err); got != tc.retryable {
				t.Fatalf("IsRetryable()=%v, want %v for %v", got, tc.retryable, err)
			}
		})
	}
}

func TestChaosOpenFailureDoesNotLeakSecret(t *testing.T) {
	t.Parallel()

	const secret = "chaos-secret-value"
	cfg := postgresx.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = 1
	cfg.Database = "postgresx_chaos"
	cfg.User = "postgresx"
	cfg.Password = postgresx.NewSecretString(secret)
	cfg.ApplicationName = "postgresx-chaos-test"
	cfg.ConnectTimeout = 25 * time.Millisecond
	cfg.HealthTimeout = 25 * time.Millisecond

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	client, err := postgresx.Open(ctx, cfg)
	if client != nil {
		_ = client.Close(t.Context())
	}
	if err == nil {
		t.Fatal("Open with canceled context unexpectedly succeeded")
	}

	surfaces := map[string]string{
		"error":       err.Error(),
		"redactedDSN": cfg.RedactedDSN(),
		"sanitize":    fmt.Sprint(cfg.Sanitize()),
	}
	for name, surface := range surfaces {
		if strings.Contains(surface, secret) {
			t.Fatalf("%s leaked postgres password: %q", name, surface)
		}
	}
}

func TestChaosNilClientOperationsFailClosed(t *testing.T) {
	t.Parallel()

	var client *postgresx.Client
	err := client.Ping(t.Context())
	if !postgresx.IsKind(err, postgresx.ErrorKindConnection) {
		t.Fatalf("nil client Ping() want %v, err=%v", postgresx.ErrorKindConnection, err)
	}

	if err := client.Close(t.Context()); err != nil {
		t.Fatalf("nil client Close() should be a no-op: %v", err)
	}
}
