package postgresxtest

import (
	"strings"
	"testing"
)

func TestConfigFromDSNParsesAndRedacts(t *testing.T) {
	password := "contract-secret"
	dsn := "postgres" + "://contract_user:" + password + "@127.0.0.1:55432/contract_db?sslmode=require"
	cfg, err := ConfigFromDSN(dsn, "postgresx-contract")
	if err != nil {
		t.Fatalf("ConfigFromDSN() error = %v", err)
	}
	if cfg.Host != "127.0.0.1" || cfg.Port != 55432 || cfg.Database != "contract_db" || cfg.User != "contract_user" {
		t.Fatalf("ConfigFromDSN() parsed address/user = %#v", cfg.Sanitize())
	}
	if cfg.SSLMode != "require" || cfg.ApplicationName != "postgresx-contract" {
		t.Fatalf("ConfigFromDSN() ssl/app = (%q, %q), want require/postgresx-contract", cfg.SSLMode, cfg.ApplicationName)
	}
	if cfg.MaxOpenConns != 4 || cfg.MinIdleConns != 1 {
		t.Fatalf("ConfigFromDSN() pool = (%d, %d), want (4, 1)", cfg.MaxOpenConns, cfg.MinIdleConns)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	redacted := cfg.RedactedDSN()
	if strings.Contains(redacted, password) {
		t.Fatalf("RedactedDSN() = %q leaked password", redacted)
	}
	if !strings.Contains(redacted, "%2A%2A%2A") && !strings.Contains(redacted, "***") {
		t.Fatalf("RedactedDSN() = %q, want masked password", redacted)
	}
}

func TestIntegrationDSNFallbackOrder(t *testing.T) {
	preferred := "postgres" + "://preferred:" + "secret" + "@localhost/preferred"
	fallback := "postgres" + "://fallback:" + "secret" + "@localhost/fallback"
	t.Setenv(EnvIntegrationDSN, preferred)
	t.Setenv(EnvPostgresTestDSN, fallback)
	if got := IntegrationDSN(); got != preferred {
		t.Fatalf("IntegrationDSN() = %q, want preferred env", got)
	}

	t.Setenv(EnvIntegrationDSN, "")
	if got := IntegrationDSN(); got != fallback {
		t.Fatalf("IntegrationDSN() = %q, want fallback env", got)
	}
}

func TestConfigFromDSNRejectsUnsupportedInput(t *testing.T) {
	if _, err := ConfigFromDSN("mysql://user:secret@localhost/db", "postgresx-contract"); err == nil {
		t.Fatal("ConfigFromDSN() accepted unsupported scheme")
	}
	if _, err := ConfigFromDSN("postgres"+"://user:"+"secret"+"@localhost:not-a-port/db", "postgresx-contract"); err == nil {
		t.Fatal("ConfigFromDSN() accepted invalid port")
	}
}
