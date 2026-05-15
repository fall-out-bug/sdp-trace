package main

import (
	"sort"
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
		"go build ./cmd/sdp-trace",
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
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("quickstartCommands returned %d commands, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("quickstartCommands[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestQuickstartCommandsIgnoresNonGoRunLines(t *testing.T) {
	// quickstartCommands must ignore go test / go build lines and only
	// extract go run ./cmd/sdp-trace commands.
	doc := strings.Join([]string{
		"```text",
		"go test -count=1 ./...",
		"go build ./cmd/sdp-trace",
		"go run ./cmd/sdp-trace --help",
		"```",
	}, "\n")
	got := quickstartCommands(doc)
	want := []string{"go run ./cmd/sdp-trace --help"}
	if len(got) != len(want) {
		t.Fatalf("quickstartCommands returned %v, want %v", got, want)
	}
	if got[0] != want[0] {
		t.Fatalf("quickstartCommands[0] = %q, want %q", got[0], want[0])
	}
}

func dummyRegistry() []string {
	return []string{
		"sdp-trace --help",
		"sdp-trace doctor",
		"sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>",
		"sdp-trace verify <run-dir>",
		"sdp-trace explain <run-dir>",
	}
}

func TestCompareQuickstartWithRegistryMissingRequired(t *testing.T) {
	// Empty quickstart should fail because required commands are missing.
	err := compareQuickstartWithRegistry("", dummyRegistry())
	if err == nil {
		t.Fatal("compareQuickstartWithRegistry succeeded, want missing required error")
	}
	if !strings.Contains(err.Error(), "missing required commands") {
		t.Fatalf("error missing required commands: %v", err)
	}
}

func TestCompareQuickstartWithRegistryDetectsStaleSubcommand(t *testing.T) {
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
	err := compareQuickstartWithRegistry(doc, dummyRegistry())
	if err == nil {
		t.Fatal("compareQuickstartWithRegistry succeeded, want stale error")
	}
	if !strings.Contains(err.Error(), "stale") {
		t.Fatalf("error missing stale drift: %v", err)
	}
}

func TestCompareQuickstartWithRegistryAcceptsAnyFlagsForKnownSubcommand(t *testing.T) {
	// isKnownCommand matches by base subcommand, so a line with a known
	// subcommand but arbitrary flags is accepted. This is an intentional
	// design choice: the registry stores usage patterns with placeholders
	// (e.g. "<name>"), not concrete flag values. This test documents the
	// limitation so future maintainers understand the boundary.
	doc := strings.Join([]string{
		"```text",
		"go run ./cmd/sdp-trace --help",
		"go run ./cmd/sdp-trace doctor",
		"go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok",
		"go run ./cmd/sdp-trace verify .sdp-trace-runs/smoke",
		"go run ./cmd/sdp-trace explain .sdp-trace-runs/smoke",
		"go run ./cmd/sdp-trace wrap --bogus-flag --another-bogus",
		"```",
	}, "\n")
	// This should NOT report stale because "wrap" is a known subcommand,
	// even though the flags are not in the registry.
	t.Log("base-command matching is intentional; flag-level drift is not checked in Slice 1")
	err := compareQuickstartWithRegistry(doc, dummyRegistry())
	if err != nil {
		t.Fatalf("compareQuickstartWithRegistry rejected known subcommand with arbitrary flags: %v", err)
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

func TestIsQuickstartCommand(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"go run ./cmd/sdp-trace --help", true},
		{"go run ./cmd/sdp-trace doctor", true},
		{"  go run ./cmd/sdp-trace wrap --name smoke  ", true},
		{"go run ./cmd/sdp-trace", false},     // bare binary, no subcommand
		{"go test -count=1 ./...", false},     // unrelated go command
		{"go build ./cmd/sdp-trace", false},   // unrelated go command
		{"# go run ./cmd/sdp-trace --help", false}, // comment
	}
	for _, c := range cases {
		got := isQuickstartCommand(c.line)
		if got != c.want {
			t.Fatalf("isQuickstartCommand(%q) = %v, want %v", c.line, got, c.want)
		}
	}
}

func TestProcessQuickstartLineIgnoresNestedFenceWithInfoString(t *testing.T) {
	// A nested fence that carries an info string (e.g. ```bash) must not
	// close the current block. Only a bare ``` line ends the block.
	doc := strings.Join([]string{
		"```text",
		"go run ./cmd/sdp-trace --help",
		"```bash", // nested fence with info string (must not close)
		"go run ./cmd/sdp-trace doctor",
		"```",
	}, "\n")
	got := quickstartCommands(doc)
	want := []string{
		"go run ./cmd/sdp-trace --help",
		"go run ./cmd/sdp-trace doctor",
	}
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("quickstartCommands returned %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("quickstartCommands[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRegistryPrefix(t *testing.T) {
	cases := []struct {
		usage string
		want  string
	}{
		{"sdp-trace wrap --name <name> [--contract <file>]", "sdp-trace wrap --name"},
		{"sdp-trace doctor [--contract <file>]", "sdp-trace doctor"},
		{"sdp-trace verify <run-dir>", "sdp-trace verify"},
		{"sdp-trace version", "sdp-trace version"},
		{"sdp-trace query --query <missing-evidence|capture-depth> <run-dir>", "sdp-trace query --query"},
	}
	for _, c := range cases {
		got := registryPrefix(c.usage)
		if got != c.want {
			t.Fatalf("registryPrefix(%q) = %q, want %q", c.usage, got, c.want)
		}
	}
}

func TestRegistryHasBaseEmpty(t *testing.T) {
	// Empty base must not match every registry entry.
	set := map[string]bool{"sdp-trace version": true, "sdp-trace doctor": true}
	if registryHasBase(set, "") {
		t.Fatal("registryHasBase(empty) matched")
	}
}
