package main

import (
	"bytes"
	"encoding/json"
	"os"
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

func TestRunJSONSchemaFixtures(t *testing.T) {
	if !hasTool("check-jsonschema") {
		t.Skip("check-jsonschema not in PATH")
	}
	state, reason := runJSONSchemaFixtures()
	if state != statePass {
		t.Fatalf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunJSONSchemaWrapDrift(t *testing.T) {
	state, reason := runJSONSchemaWrapDrift()
	if state != stateCannotVerify {
		t.Errorf("expected stateCannotVerify, got %s", state)
	}
	if reason == "" {
		t.Error("expected reason for cannot_verify")
	}
}

func TestRunOPAPolicy(t *testing.T) {
	if !hasTool("opa") {
		t.Skip("opa not in PATH")
	}
	state, reason := runOPAPolicy()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunOPANegativeFixture(t *testing.T) {
	if !hasTool("opa") {
		t.Skip("opa not in PATH")
	}
	state, reason := runOPANegativeFixture()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunCUEImport(t *testing.T) {
	if !hasTool("cue") {
		t.Skip("cue not in PATH")
	}
	state, reason := runCUEImport()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunInTotoWrap(t *testing.T) {
	if !hasTool("in-toto-run") {
		t.Skip("in-toto-run not in PATH")
	}
	state, reason := runInTotoWrap()
	if state != stateCannotVerify {
		t.Errorf("expected stateCannotVerify, got %s: %s", state, reason)
	}
}

func TestRunCosignLocalSign(t *testing.T) {
	if !hasTool("cosign") {
		t.Skip("cosign not in PATH")
	}
	state, reason := runCosignLocalSign()
	if state != stateCannotVerify {
		t.Errorf("expected stateCannotVerify, got %s: %s", state, reason)
	}
}

func TestRunSLSANegative(t *testing.T) {
	if !hasTool("slsa-verifier") {
		t.Skip("slsa-verifier not in PATH")
	}
	state, reason := runSLSANegative()
	if state != stateCannotVerify {
		t.Errorf("expected stateCannotVerify, got %s: %s", state, reason)
	}
}

func TestWrapOutputIsNotJSONObject(t *testing.T) {
	// Structural evidence that live wrap output does not conform to
	// flight-recorder-run.schema.json.
	data, err := os.ReadFile("../../examples/flight-recorder/wrap-output-drift/run.json")
	if err != nil {
		t.Fatalf("read drift fixture: %v", err)
	}
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err == nil {
		t.Fatal("expected wrap output to not be a JSON object, but it parsed as one")
	}
	// Verify the fixture contains the expected plain-text prefix from the
	// frozen verbatim capture.
	if !strings.Contains(string(data), "run_dir: .sdp-trace-runs/run-") {
		t.Fatal("expected wrap output to contain the run_dir prefix")
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
