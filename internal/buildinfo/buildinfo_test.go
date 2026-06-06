package buildinfo

import (
	"strings"
	"testing"
)

func TestResolvePrefersInjectedVersion(t *testing.T) {
	orig := Version
	t.Cleanup(func() { Version = orig })

	Version = "v1.2.3"
	if got := Resolve(); got != "v1.2.3" {
		t.Fatalf("Resolve() = %q, want injected version v1.2.3", got)
	}
}

func TestResolveFallsBackFromDev(t *testing.T) {
	orig := Version
	t.Cleanup(func() { Version = orig })

	// With the bare "dev" default, Resolve must not return "" -- it either
	// reports the module version (under `go test` that is usually empty/devel,
	// so it lands on "dev") or "dev". The contract is "never empty".
	Version = "dev"
	if got := Resolve(); got == "" {
		t.Fatal("Resolve() returned empty string for dev default")
	}
}

func TestStringIncludesCommitAndDate(t *testing.T) {
	origV, origC, origD := Version, Commit, Date
	t.Cleanup(func() { Version, Commit, Date = origV, origC, origD })

	Version, Commit, Date = "v0.1.0", "abc1234", "2026-06-05T00:00:00Z"
	got := String()
	for _, want := range []string{"v0.1.0", "abc1234", "2026-06-05T00:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Errorf("String() = %q, missing %q", got, want)
		}
	}
}

func TestStringBareVersion(t *testing.T) {
	origV, origC, origD := Version, Commit, Date
	t.Cleanup(func() { Version, Commit, Date = origV, origC, origD })

	Version, Commit, Date = "v0.1.0", "", ""
	if got := String(); got != "v0.1.0" {
		t.Fatalf("String() = %q, want bare v0.1.0", got)
	}
}
