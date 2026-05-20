package main

import (
	"bytes"
	"strings"
	"testing"
)

// testRegistry is a minimal registry used for CLI-level tests so they do not
// depend on optional external tools being present on PATH.
var testRegistry = []probe{
	{
		Name:        "always-pass",
		Description: "test probe that always passes",
		Run:         func() (verifierState, string) { return statePass, "ok" },
	},
	{
		Name:        "always-not-assessed",
		Description: "test probe that is not assessed",
		Run:         func() (verifierState, string) { return stateNotAssessed, "na" },
	},
	{
		Name:        "always-cannot-verify",
		Description: "test probe that cannot verify",
		Run:         func() (verifierState, string) { return stateCannotVerify, "cant" },
	},
}

func TestRun_AllProbes_Text(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{}, &out, &out, testRegistry)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	got := out.String()
	for _, p := range testRegistry {
		if !strings.Contains(got, p.Name) {
			t.Errorf("output missing probe %q", p.Name)
		}
	}
}

func TestRun_AllProbes_JSON(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-json"}, &out, &out, testRegistry)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.HasPrefix(strings.TrimSpace(out.String()), "[") {
		t.Fatalf("expected JSON array output, got: %s", out.String())
	}
}

func TestRun_SingleProbe(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-probe", testRegistry[0].Name}, &out, &out, testRegistry)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), testRegistry[0].Name) {
		t.Fatalf("expected output for probe %q, got: %s", testRegistry[0].Name, out.String())
	}
}

func TestRun_ListProbes(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-list"}, &out, &out, testRegistry)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	for _, p := range testRegistry {
		if !strings.Contains(out.String(), p.Name) {
			t.Errorf("output missing probe %q", p.Name)
		}
	}
}

func TestRun_UnknownProbe(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-probe", "nonexistent"}, &out, &out, testRegistry)
	if code != 2 {
		t.Fatalf("expected exit 2, got %d", code)
	}
}

func TestRun_SingleProbeJSON(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-json", "-probe", testRegistry[0].Name}, &out, &out, testRegistry)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.HasPrefix(strings.TrimSpace(out.String()), "[") {
		t.Fatalf("expected JSON array output, got: %s", out.String())
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name    string
		results []probeResult
		want    int
	}{
		{"all pass", []probeResult{{State: statePass}}, 0},
		{"one fail", []probeResult{{State: statePass}, {State: stateFail}}, 1},
		{"not assessed", []probeResult{{State: stateNotAssessed}}, 0},
		{"cannot verify", []probeResult{{State: stateCannotVerify}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exitCode(tt.results); got != tt.want {
				t.Errorf("exitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}
