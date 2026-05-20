package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_AllProbes_Text(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	got := out.String()
	for _, p := range registry {
		if !strings.Contains(got, p.Name) {
			t.Errorf("output missing probe %q", p.Name)
		}
	}
}

func TestRun_AllProbes_JSON(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-json"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.HasPrefix(strings.TrimSpace(out.String()), "[") {
		t.Fatalf("expected JSON array output, got: %s", out.String())
	}
}

func TestRun_SingleProbe(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-probe", registry[0].Name}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), registry[0].Name) {
		t.Fatalf("expected output for probe %q, got: %s", registry[0].Name, out.String())
	}
}

func TestRun_UnknownProbe(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-probe", "nonexistent"}, &out, &out)
	if code != 2 {
		t.Fatalf("expected exit 2, got %d", code)
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
