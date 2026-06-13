package contracts

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/ZoneCNH/foundationx/pkg/foundationx"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

type jsonObject map[string]any

func TestErrorSchemaMatchesPublicContract(t *testing.T) {
	schema := readSchema(t, "error.schema.json")

	assertSetEqual(t, "required", stringList(t, schema["required"]), []string{
		"kind",
		"message",
		"retryable",
	})

	expectedKinds := []string{
		string(foundationx.ErrorKindConfig),
		string(foundationx.ErrorKindValidation),
		string(foundationx.ErrorKindConnection),
		string(foundationx.ErrorKindUnavailable),
		string(foundationx.ErrorKindTimeout),
		string(foundationx.ErrorKindAuth),
		string(foundationx.ErrorKindConflict),
		string(foundationx.ErrorKindRateLimit),
		string(foundationx.ErrorKindCanceled),
		string(foundationx.ErrorKindNotFound),
		string(foundationx.ErrorKindAlreadyExist),
		string(foundationx.ErrorKindInternal),
	}
	assertSetEqual(t, "error kinds", enumValues(t, schema, "kind"), expectedKinds)

	if err := postgresx.MapError("contract", context.DeadlineExceeded); !postgresx.IsRetryable(err) {
		t.Fatalf("deadline exceeded should be retryable")
	}
	if err := postgresx.MapError("contract", context.Canceled); !foundationx.IsKind(err, foundationx.ErrorKindCanceled) {
		t.Fatalf("context canceled should map to canceled, got %v", err)
	}
}

func TestHealthSchemaMatchesPublicContract(t *testing.T) {
	schema := readSchema(t, "health.schema.json")
	properties := object(t, schema["properties"])

	assertSetEqual(t, "required", stringList(t, schema["required"]), []string{
		"name",
		"status",
		"message",
		"checked_at",
		"latency_ms",
		"metadata",
	})
	assertSetEqual(t, "health statuses", enumValues(t, schema, "status"), []string{
		string(foundationx.HealthHealthy),
		string(foundationx.HealthDegraded),
		string(foundationx.HealthUnhealthy),
	})
	assertStructProperties(t, reflect.TypeOf(foundationx.HealthStatus{}), properties, map[string]string{
		"Name":      "name",
		"Status":    "status",
		"Message":   "message",
		"CheckedAt": "checked_at",
		"LatencyMs": "latency_ms",
		"Metadata":  "metadata",
	})

	metadata := object(t, properties["metadata"])
	if metadata["type"] != "object" {
		t.Fatalf("metadata type = %v, want object", metadata["type"])
	}

	statsProperties := object(t, object(t, properties["pool_stats"])["properties"])
	assertStructProperties(t, reflect.TypeOf(postgresx.PoolStats{}), statsProperties, map[string]string{
		"TotalConns":        "total_conns",
		"IdleConns":         "idle_conns",
		"AcquiredConns":     "acquired_conns",
		"ConstructingConns": "constructing_conns",
		"MaxConns":          "max_conns",
	})
	assertSetEqual(t, "pool stats required", stringList(t, object(t, properties["pool_stats"])["required"]), keys(statsProperties))
}

func TestConfigSchemaMatchesPublicContract(t *testing.T) {
	schema := readSchema(t, "config.schema.json")
	properties := object(t, schema["properties"])

	assertSetEqual(t, "required", stringList(t, schema["required"]), []string{
		"host",
		"port",
		"database",
		"user",
		"password",
		"sslmode",
	})
	assertStructProperties(t, reflect.TypeOf(postgresx.Config{}), properties, map[string]string{
		"Host":            "host",
		"Port":            "port",
		"Database":        "database",
		"User":            "user",
		"Password":        "password",
		"SSLMode":         "sslmode",
		"MaxOpenConns":    "max_open_conns",
		"MinIdleConns":    "min_idle_conns",
		"ConnectTimeout":  "connect_timeout_ms",
		"HealthTimeout":   "health_timeout_ms",
		"MaxConnLifetime": "max_conn_lifetime_ms",
		"MaxConnIdleTime": "max_conn_idle_time_ms",
		"ApplicationName": "application_name",
	})

	for _, name := range []string{
		"connect_timeout_ms",
		"health_timeout_ms",
		"max_conn_lifetime_ms",
		"max_conn_idle_time_ms",
	} {
		property := object(t, properties[name])
		if property["type"] != "integer" {
			t.Fatalf("%s type = %v, want integer", name, property["type"])
		}
		if property["minimum"] != float64(0) {
			t.Fatalf("%s minimum = %v, want 0", name, property["minimum"])
		}
	}

	secret := object(t, properties["password"])
	if secret["x-secret"] != true {
		t.Fatalf("password x-secret = %v, want true", secret["x-secret"])
	}
}

func TestMetricsContractDocumentsPublicHooks(t *testing.T) {
	contents, err := os.ReadFile("metrics.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(contents)
	for _, snippet := range []string{
		"Metrics",
		"IncCounter",
		"ObserveHistogram",
		"SetGauge",
		"WithMetrics",
		"postgresx.query.total",
		"postgresx.query.duration_seconds",
		"postgresx.tx.total",
		"postgresx.tx.duration_seconds",
		"postgresx.health.total",
		"postgresx.health.latency_seconds",
		"postgresx.pool.connections",
	} {
		if !strings.Contains(text, snippet) {
			t.Fatalf("metrics contract missing %q", snippet)
		}
	}
}

func TestVersionContractDocumentsPublicAPIBaseline(t *testing.T) {
	contents, err := os.ReadFile("../VERSION")
	if err != nil {
		t.Fatal(err)
	}
	version := strings.TrimSpace(string(contents))
	if version == "" {
		t.Fatal("VERSION is empty")
	}
	if postgresx.Version != version {
		t.Fatalf("postgresx.Version = %q, want VERSION %q", postgresx.Version, version)
	}

	for _, path := range []string{"public_api.md", "../docs/api.md"} {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		want := "Stable " + version + " surfaces:"
		if !strings.Contains(string(contents), want) {
			t.Fatalf("%s missing %q", path, want)
		}
	}
}

func readSchema(t *testing.T, name string) jsonObject {
	t.Helper()
	contents, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	var schema jsonObject
	if err := json.Unmarshal(contents, &schema); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"$schema", "title", "type", "properties"} {
		if schema[field] == nil {
			t.Fatalf("%s missing %s", name, field)
		}
	}
	return schema
}

func enumValues(t *testing.T, schema jsonObject, propertyName string) []string {
	t.Helper()
	properties := object(t, schema["properties"])
	property := object(t, properties[propertyName])
	return stringList(t, property["enum"])
}

func assertStructProperties(t *testing.T, typ reflect.Type, properties jsonObject, mapping map[string]string) {
	t.Helper()
	for fieldName, propertyName := range mapping {
		if _, ok := typ.FieldByName(fieldName); !ok {
			t.Fatalf("%s missing public field %s", typ.Name(), fieldName)
		}
		if _, ok := properties[propertyName]; !ok {
			t.Fatalf("%s schema missing property %s", typ.Name(), propertyName)
		}
	}
}

func assertSetEqual(t *testing.T, label string, got, want []string) {
	t.Helper()
	got = slices.Clone(got)
	want = slices.Clone(want)
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
}

func stringList(t *testing.T, value any) []string {
	t.Helper()
	raw, ok := value.([]any)
	if !ok {
		t.Fatalf("value has type %T, want []any", value)
	}
	values := make([]string, 0, len(raw))
	for _, item := range raw {
		text, ok := item.(string)
		if !ok {
			t.Fatalf("list item has type %T, want string", item)
		}
		values = append(values, text)
	}
	return values
}

func object(t *testing.T, value any) jsonObject {
	t.Helper()
	obj, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("value has type %T, want object", value)
	}
	return obj
}

func keys(obj jsonObject) []string {
	values := make([]string, 0, len(obj))
	for key := range obj {
		values = append(values, key)
	}
	return values
}
