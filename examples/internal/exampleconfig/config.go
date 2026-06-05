package exampleconfig

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

// RuntimeConfig is the example boundary around postgresx.Config. The core
// package never reads env; examples do so only to demonstrate caller-owned
// configuration loading.
type RuntimeConfig struct {
	Config postgresx.Config
	Live   bool
}

// FromEnv builds an example config. It returns Live=false unless the caller
// explicitly sets POSTGRESX_EXAMPLE_LIVE=1, allowing `go run ./examples/...`
// to act as a dry-run smoke check without requiring a database.
func FromEnv(applicationName string) (RuntimeConfig, error) {
	live := os.Getenv("POSTGRESX_EXAMPLE_LIVE") == "1"

	cfg := postgresx.DefaultConfig()
	cfg.Host = getenv("POSTGRES_HOST", "localhost")
	cfg.Database = getenv("POSTGRES_DATABASE", "postgres")
	cfg.User = getenv("POSTGRES_USER", "postgres")
	password, err := envPassword(live)
	if err != nil {
		return RuntimeConfig{}, err
	}
	cfg.Password = foundationx.NewSecretString(password)
	cfg.SSLMode = getenv("POSTGRES_SSLMODE", cfg.SSLMode)
	cfg.ApplicationName = applicationName

	port, err := strconv.Atoi(getenv("POSTGRES_PORT", "5432"))
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("parse POSTGRES_PORT: %w", err)
	}
	cfg.Port = port
	if live {
		if err := requireLiveEnv("POSTGRES_HOST", "POSTGRES_DATABASE", "POSTGRES_USER", "POSTGRES_PASSWORD"); err != nil {
			return RuntimeConfig{}, err
		}
	}
	return RuntimeConfig{Config: cfg, Live: live}, nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envPassword(live bool) (string, error) {
	if value := os.Getenv("POSTGRES_PASSWORD"); value != "" {
		return value, nil
	}
	if live {
		return "", fmt.Errorf("POSTGRES_PASSWORD must be set when POSTGRESX_EXAMPLE_LIVE=1")
	}
	return "postgres", nil
}

func requireLiveEnv(keys ...string) error {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			return fmt.Errorf("%s must be set when POSTGRESX_EXAMPLE_LIVE=1", key)
		}
	}
	return nil
}
