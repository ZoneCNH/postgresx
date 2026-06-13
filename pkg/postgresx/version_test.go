package postgresx

import "testing"

func TestVersionContract(t *testing.T) {
	if ModuleName != "github.com/ZoneCNH/postgresx" {
		t.Fatalf("ModuleName = %q, want github.com/ZoneCNH/postgresx", ModuleName)
	}
	if Version != "v1.0.0" {
		t.Fatalf("Version = %q, want v1.0.0", Version)
	}
}
