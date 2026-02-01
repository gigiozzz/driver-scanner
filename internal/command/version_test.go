package command

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

func TestVersionOptions_Run_Full(t *testing.T) {
	// Override package-level vars for deterministic output.
	origVersion, origCommit, origDate := Version, CommitHash, BuildDate
	Version = "1.2.3"
	CommitHash = "abc1234"
	BuildDate = "2025-01-01T00:00:00Z"
	t.Cleanup(func() {
		Version, CommitHash, BuildDate = origVersion, origCommit, origDate
	})

	var buf bytes.Buffer
	o := &VersionOptions{Out: &buf}

	if err := o.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	expected := fmt.Sprintf(
		"driver-scanner version: 1.2.3\nCommit: abc1234\nBuilt: 2025-01-01T00:00:00Z\nGo version: %s\nOS/Arch: %s/%s\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH,
	)

	if got != expected {
		t.Errorf("unexpected output:\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestVersionOptions_Run_Short(t *testing.T) {
	origVersion := Version
	Version = "1.2.3"
	t.Cleanup(func() { Version = origVersion })

	var buf bytes.Buffer
	o := &VersionOptions{Short: true, Out: &buf}

	if err := o.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if got != "1.2.3\n" {
		t.Errorf("unexpected output: got %q, want %q", got, "1.2.3\n")
	}
}
