package postgresxtest

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

const (
	// EnvIntegrationDSN is the preferred integration DSN for postgresx tests.
	EnvIntegrationDSN = "POSTGRESX_INTEGRATION_DSN"
	// EnvPostgresTestDSN is the compatibility DSN used by shared CI fixtures.
	EnvPostgresTestDSN = "POSTGRES_TEST_DSN"
)

// IntegrationDSN returns the configured PostgreSQL test DSN, if any.
func IntegrationDSN() string {
	if dsn := os.Getenv(EnvIntegrationDSN); dsn != "" {
		return dsn
	}
	return os.Getenv(EnvPostgresTestDSN)
}

// ConfigFromDSN converts a PostgreSQL URL into an explicit postgresx Config.
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
	cfg.Password = postgresx.NewSecretString(secret)
	cfg.SSLMode = sslmode
	cfg.MaxOpenConns = 4
	cfg.MinIdleConns = 1
	cfg.ConnectTimeout = 15 * time.Second
	cfg.HealthTimeout = 15 * time.Second
	cfg.ApplicationName = applicationName
	return cfg, nil
}
