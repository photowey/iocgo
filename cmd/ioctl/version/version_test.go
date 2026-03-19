package version

import "testing"

func TestNow(t *testing.T) {
	originalVersion := Version
	defer func() {
		Version = originalVersion
	}()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{name: "plain semver", version: "0.1.0", want: "v0.1.0"},
		{name: "tag semver", version: "v0.1.0", want: "v0.1.0"},
		{name: "dev build", version: "dev", want: "dev"},
		{name: "unknown", version: "", want: "(unknown)"},
	}

	for _, tt := range tests {
		Version = tt.version
		if got := Now(); got != tt.want {
			t.Fatalf("%s: Now() = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestSummary(t *testing.T) {
	originalVersion := Version
	originalCommit := Commit
	originalBuildTime := BuildTime
	defer func() {
		Version = originalVersion
		Commit = originalCommit
		BuildTime = originalBuildTime
	}()

	Version = "v0.2.0"
	Commit = "abc1234"
	BuildTime = "2026-03-20T00:00:00Z"

	if got, want := Summary(), "v0.2.0 (commit abc1234, built 2026-03-20T00:00:00Z)"; got != want {
		t.Fatalf("Summary() = %q, want %q", got, want)
	}

	Commit = ""
	BuildTime = ""

	if got, want := Summary(), "v0.2.0"; got != want {
		t.Fatalf("Summary() without metadata = %q, want %q", got, want)
	}
}
