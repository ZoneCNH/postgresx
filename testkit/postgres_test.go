package testkit

import (
	"strings"
	"testing"
)

func TestConfigFromDSN(t *testing.T) {
	dsn := "postgres://" + "alice:" + "secret" + "@localhost:5433/app?sslmode=require"

	cfg, err := ConfigFromDSN(dsn, "postgresx-test")
	if err != nil {
		t.Fatalf("ConfigFromDSN() error = %v", err)
	}
	if cfg.Host != "localhost" || cfg.Port != 5433 || cfg.Database != "app" || cfg.User != "alice" {
		t.Fatalf("ConfigFromDSN() = %+v, want parsed host/port/database/user", cfg.Sanitize())
	}
	if cfg.SSLMode != "require" {
		t.Fatalf("SSLMode = %q, want require", cfg.SSLMode)
	}
	if cfg.ApplicationName != "postgresx-test" {
		t.Fatalf("ApplicationName = %q, want postgresx-test", cfg.ApplicationName)
	}
	if strings.Contains(cfg.RedactedDSN(), "secret") {
		t.Fatalf("RedactedDSN() leaked the driver secret: %q", cfg.RedactedDSN())
	}
}
