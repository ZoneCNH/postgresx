package exampleconfig

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/bytechainx/postgresx"
)

// FromEnv builds a postgresx Config from non-secret connection settings and a
// secret supplied by the caller's environment.
func FromEnv(applicationName string) (postgresx.Config, error) {
	cfg := postgresx.DefaultConfig()
	cfg.Host = getenv("POSTGRES_HOST", "localhost")
	cfg.Database = getenv("POSTGRES_DATABASE", "postgres")
	cfg.User = getenv("POSTGRES_USER", "postgres")
	cfg.Password = foundationx.NewSecretString(firstEnv("POSTGRES_SECRET", "POSTGRES_PASSWORD"))
	cfg.SSLMode = getenv("POSTGRES_SSLMODE", cfg.SSLMode)
	cfg.ApplicationName = applicationName

	port, err := strconv.Atoi(getenv("POSTGRES_PORT", "5432"))
	if err != nil {
		return postgresx.Config{}, fmt.Errorf("parse POSTGRES_PORT: %w", err)
	}
	cfg.Port = port
	return cfg, nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}
