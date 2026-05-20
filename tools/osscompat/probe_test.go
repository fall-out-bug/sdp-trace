package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHasTool_Missing(t *testing.T) {
	if hasTool("this-tool-definitely-does-not-exist-017") {
		t.Error("expected missing tool to return false")
	}
}

func TestHasTool_Present(t *testing.T) {
	// `go` is guaranteed to be on PATH in any Go test environment.
	if !hasTool("go") {
		t.Error("expected `go` to be found on PATH")
	}
}

func TestRunProbe_NotAssessedWhenToolMissing(t *testing.T) {
	p := probe{
		Name:      "missing-tool-test",
		NeedsTool: "this-tool-definitely-does-not-exist-017",
		Run: func() (verifierState, string) {
			return statePass, "should not run"
		},
	}
	r := runProbe(p)
	if r.State != stateNotAssessed {
		t.Errorf("expected stateNotAssessed, got %s", r.State)
	}
	if r.Reason == "" {
		t.Error("expected reason for not_assessed")
	}
}

func TestRunProbe_RunsWhenToolPresent(t *testing.T) {
	p := probe{
		Name:      "go-version",
		NeedsTool: "go",
		Run: func() (verifierState, string) {
			return statePass, "ok"
		},
	}
	r := runProbe(p)
	if r.State != statePass {
		t.Errorf("expected statePass, got %s", r.State)
	}
	if r.Reason != "ok" {
		t.Errorf("expected reason 'ok', got %q", r.Reason)
	}
}

func TestPrintResults_Text(t *testing.T) {
	var buf bytes.Buffer
	results := []probeResult{
		{Name: "a", State: statePass, Reason: "ok"},
	}
	if err := printResults(&buf, results, false); err != nil {
		t.Fatalf("printResults: %v", err)
	}
	if !strings.Contains(buf.String(), "a") || !strings.Contains(buf.String(), "pass") {
		t.Errorf("unexpected text output: %s", buf.String())
	}
}

func TestPrintResults_JSON(t *testing.T) {
	var buf bytes.Buffer
	results := []probeResult{
		{Name: "a", State: statePass},
	}
	if err := printResults(&buf, results, true); err != nil {
		t.Fatalf("printResults: %v", err)
	}
	if !strings.Contains(buf.String(), `"name"`) {
		t.Errorf("unexpected JSON output: %s", buf.String())
	}
}
