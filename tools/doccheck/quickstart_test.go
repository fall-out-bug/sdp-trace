package main

import (
	"strings"
	"testing"
)

func TestQuickstartCommandsExtractsGoRunLines(t *testing.T) {
	doc := strings.Join([]string{
		"## Smoke Path",
		"",
		"```text",
		"go test -count=1 ./...",
		"go run ./cmd/sdp-trace --help",
		"go run ./cmd/sdp-trace doctor",
		"```",
		"",
		"Some prose.",
		"",
		"```bash",
		"go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok",
		"```",
	}, "\n")
	got := quickstartCommands(doc)
	want := []string{
		"go run ./cmd/sdp-trace --help",
		"go run ./cmd/sdp-trace doctor",
		"go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok",
	}
	if len(got) != len(want) {
		t.Fatalf("quickstartCommands returned %d commands, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("quickstartCommands[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCompareQuickstartWithRegistryMissingRequired(t *testing.T) {
	// Empty quickstart should fail because required commands are missing.
	err := compareQuickstartWithRegistry("")
	if err == nil {
		t.Fatal("compareQuickstartWithRegistry succeeded, want missing required error")
	}
	if !strings.Contains(err.Error(), "missing required commands") {
		t.Fatalf("error missing required commands: %v", err)
	}
}

func TestCompareQuickstartWithRegistryDetectsStale(t *testing.T) {
	doc := strings.Join([]string{
		"```text",
		"go run ./cmd/sdp-trace --help",
		"go run ./cmd/sdp-trace doctor",
		"go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok",
		"go run ./cmd/sdp-trace verify .sdp-trace-runs/smoke",
		"go run ./cmd/sdp-trace explain .sdp-trace-runs/smoke",
		"go run ./cmd/sdp-trace stale-command-that-does-not-exist",
		"```",
	}, "\n")
	err := compareQuickstartWithRegistry(doc)
	if err == nil {
		t.Fatal("compareQuickstartWithRegistry succeeded, want stale error")
	}
	if !strings.Contains(err.Error(), "stale") {
		t.Fatalf("error missing stale drift: %v", err)
	}
}

func TestIsKnownCommandExactMatch(t *testing.T) {
	set := map[string]bool{"sdp-trace version": true}
	if !isKnownCommand("sdp-trace version", set) {
		t.Fatal("isKnownCommand exact match failed")
	}
}

func TestIsKnownCommandPrefixMatch(t *testing.T) {
	set := map[string]bool{"sdp-trace wrap --name <name>": true}
	if !isKnownCommand("go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok", set) {
		t.Fatal("isKnownCommand prefix match failed")
	}
}

func TestIsKnownCommandUnknown(t *testing.T) {
	set := map[string]bool{"sdp-trace version": true}
	if isKnownCommand("sdp-trace unknown", set) {
		t.Fatal("isKnownCommand unknown command matched")
	}
}
