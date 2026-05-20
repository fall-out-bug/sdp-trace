package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func skipUnlessIntegration(t *testing.T) {
	if os.Getenv("SDPTRACE_INTEGRATION") != "1" {
		t.Skip("set SDPTRACE_INTEGRATION=1 to run tests that build the product binary")
	}
}

func TestRun_List(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-list"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	for _, b := range builtIns {
		if !strings.Contains(out.String(), b.Name) {
			t.Errorf("output missing benchmark %q", b.Name)
		}
	}
}

func TestRun_CustomCommand(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-n", "3", "go", "env", "GOOS"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), "go env GOOS") {
		t.Fatalf("expected output for custom command, got: %s", out.String())
	}
}

func TestRun_JSON(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-json", "-n", "2", "go", "env", "GOOS"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.HasPrefix(strings.TrimSpace(out.String()), "[") {
		t.Fatalf("expected JSON array output, got: %s", out.String())
	}
}

func TestRun_BuiltinBench(t *testing.T) {
	skipUnlessIntegration(t)
	var out bytes.Buffer
	code := run([]string{"-bench", "sdp-trace-version", "-n", "1"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), "sdp-trace-version") {
		t.Fatalf("expected output for builtin bench, got: %s", out.String())
	}
}

func TestRun_UnknownBench(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-bench", "nonexistent"}, &out, &out)
	if code != 2 {
		t.Fatalf("expected exit 2, got %d", code)
	}
}

func TestRun_MissingCommand(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"this-command-does-not-exist-017"}, &out, &out)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name    string
		results []benchmarkResult
		want    int
	}{
		{"all ok", []benchmarkResult{{Name: "a", MedianMs: 1}}, 0},
		{"one error", []benchmarkResult{{Name: "a", Error: "fail"}}, 1},
		{"mixed", []benchmarkResult{{Name: "a", MedianMs: 1}, {Name: "b", Error: "fail"}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exitCode(tt.results); got != tt.want {
				t.Errorf("exitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}
