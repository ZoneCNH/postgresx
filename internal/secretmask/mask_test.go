package secretmask

import (
	"strings"
	"testing"
)

func TestMaskPostgresURL(t *testing.T) {
	raw := strings.Join([]string{"postgres://user:", "secret-value", "@localhost/db"}, "")
	masked := Mask(raw)
	if strings.Contains(masked, "secret-value") {
		t.Fatalf("secret leaked: %s", masked)
	}
	if !strings.Contains(masked, "<masked>") {
		t.Fatalf("missing mask marker: %s", masked)
	}
}

func TestMaskKeyValuePasswords(t *testing.T) {
	raw := strings.Join([]string{
		"host=localhost user=postgres ",
		"password=",
		"secret-value ",
		"PGPASSWORD=",
		"other",
	}, "")
	masked := Mask(raw)
	if strings.Contains(masked, "secret-value") || strings.Contains(masked, "other") {
		t.Fatalf("secret leaked: %s", masked)
	}
}
