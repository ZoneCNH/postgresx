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
	cfg := postgresx.DefaultConfig()
	cfg.Host = getenv("POSTGRES_HOST", "localhost")
	cfg.Database = getenv("POSTGRES_DATABASE", "postgres")
	cfg.User = getenv("POSTGRES_USER", "postgres")
	cfg.Password = foundationx.NewSecretString(getenv("POSTGRES_PASSWORD", "postgres"))
	cfg.SSLMode = getenv("POSTGRES_SSLMODE", cfg.SSLMode)
	cfg.ApplicationName = applicationName

	port, err := strconv.Atoi(getenv("POSTGRES_PORT", "5432"))
	if err != nil {
		return RuntimeConfig{}, fmt.Errorf("parse POSTGRES_PORT: %w", err)
	}
	cfg.Port = port
	return RuntimeConfig{Config: cfg, Live: os.Getenv("POSTGRESX_EXAMPLE_LIVE") == "1"}, nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
