package postgresx

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
)

func validTestConfig() Config {
	cfg := DefaultConfig()
	cfg.Host = "localhost"
	cfg.Database = "app"
	cfg.User = "postgres"
	cfg.Password = foundationx.NewSecretString("secret")
	cfg.ApplicationName = "postgresx-test"
	return cfg
}

func TestConfigWithDefaultsAndValidate(t *testing.T) {
	cfg := Config{
		Host:     "localhost",
		Database: "app",
		User:     "postgres",
		Password: foundationx.NewSecretString("secret"),
	}.withDefaults()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if cfg.Port != 5432 {
		t.Fatalf("Port = %d, want 5432", cfg.Port)
	}
	if cfg.SSLMode != "disable" {
		t.Fatalf("SSLMode = %q, want disable", cfg.SSLMode)
	}
	if cfg.ConnectTimeout != 5*time.Second {
		t.Fatalf("ConnectTimeout = %s, want 5s", cfg.ConnectTimeout)
	}
}

func TestConfigDSNAndRedactedDSN(t *testing.T) {
	cfg := validTestConfig()
	dsn := cfg.DSN()
	redacted := cfg.RedactedDSN()

	if !strings.Contains(dsn, "secret") {
		t.Fatalf("DSN() did not include the driver secret: %q", dsn)
	}
	if strings.Contains(redacted, "secret") {
		t.Fatalf("RedactedDSN() leaked the driver secret: %q", redacted)
	}

	parsed, err := url.Parse(redacted)
	if err != nil {
		t.Fatalf("parse redacted DSN: %v", err)
	}
	gotSecret, ok := parsed.User.Password()
	if !ok || gotSecret != "***" {
		t.Fatalf("redacted password = %q, %v; want masked password", gotSecret, ok)
	}
	if parsed.Query().Get("sslmode") != "disable" {
		t.Fatalf("sslmode query = %q, want disable", parsed.Query().Get("sslmode"))
	}
	if parsed.Query().Get("application_name") != "postgresx-test" {
		t.Fatalf("application_name query = %q, want postgresx-test", parsed.Query().Get("application_name"))
	}
}

func TestConfigValidateRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Config)
	}{
		{
			name: "missing password",
			mutate: func(cfg *Config) {
				var zero foundationx.SecretString
				cfg.Password = zero
			},
		},
		{
			name: "negative connect timeout",
			mutate: func(cfg *Config) {
				cfg.ConnectTimeout = -time.Second
			},
		},
		{
			name: "min idle exceeds max open",
			mutate: func(cfg *Config) {
				cfg.MinIdleConns = 5
				cfg.MaxOpenConns = 2
			},
		},
		{
			name: "empty sslmode",
			mutate: func(cfg *Config) {
				cfg.SSLMode = ""
			},
		},
		{
			name: "negative max open conns",
			mutate: func(cfg *Config) {
				cfg.MaxOpenConns = -1
			},
		},
		{
			name: "negative min idle conns",
			mutate: func(cfg *Config) {
				cfg.MinIdleConns = -1
			},
		},
		{
			name: "negative max conn lifetime",
			mutate: func(cfg *Config) {
				cfg.MaxConnLifetime = -time.Minute
			},
		},
		{
			name: "negative max conn idle time",
			mutate: func(cfg *Config) {
				cfg.MaxConnIdleTime = -time.Minute
			},
		},
		{
			name: "negative health timeout",
			mutate: func(cfg *Config) {
				cfg.HealthTimeout = -time.Second
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validTestConfig()
			tt.mutate(&cfg)

			err := cfg.Validate()
			if !foundationx.IsKind(err, foundationx.ErrorKindConfig) {
				t.Fatalf("Validate() error = %v, want foundation config error", err)
			}
		})
	}
}

func TestConfigSanitizeMasksSecret(t *testing.T) {
	cfg := validTestConfig()
	sanitized := cfg.Sanitize()

	if sanitized.Password != "***" {
		t.Fatalf("Sanitize().Password = %q, want masked value", sanitized.Password)
	}
	if sanitized.Host != cfg.Host || sanitized.Database != cfg.Database || sanitized.User != cfg.User {
		t.Fatalf("Sanitize() lost non-sensitive fields: %+v", sanitized)
	}
}
