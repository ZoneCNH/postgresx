package postgresx

import (
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func benchmarkConfig() Config {
	cfg := DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = 5432
	cfg.Database = "postgresx_benchmark"
	cfg.User = "postgresx"
	cfg.Password = NewSecretString("benchmark-secret-value")
	cfg.ApplicationName = "postgresx-benchmark"
	return cfg
}

func BenchmarkConfigRedactedDSNSmoke(b *testing.B) {
	cfg := benchmarkConfig()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		dsn := cfg.RedactedDSN()
		if dsn == "" || strings.Contains(dsn, "benchmark-secret-value") {
			b.Fatalf("RedactedDSN produced unsafe DSN: %q", dsn)
		}
	}
}

func BenchmarkConfigSanitizeSmoke(b *testing.B) {
	cfg := benchmarkConfig()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		safe := cfg.Sanitize()
		if strings.Contains(safe.Password, "benchmark-secret-value") {
			b.Fatalf("Sanitize leaked password: %q", safe.Password)
		}
	}
}

func BenchmarkMapErrorPostgresCodeSmoke(b *testing.B) {
	pgErr := &pgconn.PgError{Code: "40001", Message: "serialization failure"}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := MapError("benchmark.serialization", pgErr)
		if !IsKind(err, ErrorKindConflict) || !IsRetryable(err) {
			b.Fatalf("MapError produced unexpected classification: %v", err)
		}
	}
}
