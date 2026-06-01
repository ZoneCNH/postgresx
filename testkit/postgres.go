package testkit

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

// Options controls StartPostgres.
type Options struct {
	DSN             string
	ApplicationName string
}

// StartPostgres opens a test PostgreSQL client from Options.DSN,
// POSTGRESX_INTEGRATION_DSN, or POSTGRES_TEST_DSN. It skips the test when no
// test DSN is configured.
func StartPostgres(ctx context.Context, t testing.TB, opts Options) *Fixture {
	t.Helper()
	dsn := opts.DSN
	if dsn == "" {
		dsn = integrationDSN()
	}
	if dsn == "" {
		t.Skip("POSTGRESX_INTEGRATION_DSN or POSTGRES_TEST_DSN is not set")
	}
	appName := opts.ApplicationName
	if appName == "" {
		appName = "postgresx-testkit"
	}
	cfg, err := ConfigFromDSN(dsn, appName)
	if err != nil {
		t.Fatalf("parse postgres fixture config: %v", err)
	}
	client, err := postgresx.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("open postgres fixture: %v", err)
	}
	fixture := &Fixture{client: client}
	t.Cleanup(fixture.Close)
	return fixture
}

// ConfigFromDSN converts a PostgreSQL URL into the explicit postgresx Config
// contract used by the public API.
func ConfigFromDSN(dsn string, applicationName string) (postgresx.Config, error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return postgresx.Config{}, err
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return postgresx.Config{}, fmt.Errorf("unsupported PostgreSQL scheme %q", parsed.Scheme)
	}

	port := 5432
	if raw := parsed.Port(); raw != "" {
		port, err = strconv.Atoi(raw)
		if err != nil {
			return postgresx.Config{}, fmt.Errorf("parse PostgreSQL port: %w", err)
		}
	}
	database := strings.TrimPrefix(parsed.Path, "/")
	sslmode := parsed.Query().Get("sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}
	secret, _ := parsed.User.Password()

	cfg := postgresx.DefaultConfig()
	cfg.Host = parsed.Hostname()
	cfg.Port = port
	cfg.Database = database
	cfg.User = parsed.User.Username()
	cfg.Password = foundationx.NewSecretString(secret)
	cfg.SSLMode = sslmode
	cfg.MaxOpenConns = 4
	cfg.MinIdleConns = 1
	cfg.HealthTimeout = 2 * time.Second
	cfg.ApplicationName = applicationName
	return cfg, nil
}

func integrationDSN() string {
	if dsn := os.Getenv("POSTGRESX_INTEGRATION_DSN"); dsn != "" {
		return dsn
	}
	return os.Getenv("POSTGRES_TEST_DSN")
}
