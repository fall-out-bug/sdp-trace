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

func skipUnlessIntegration(t *testing.T) {
	if os.Getenv("SDPTRACE_INTEGRATION") != "1" {
		t.Skip("set SDPTRACE_INTEGRATION=1 to run tests that invoke optional external CLIs")
	}
}

func TestRunJSONSchemaFixtures(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("check-jsonschema") {
		t.Skip("check-jsonschema not in PATH")
	}
	state, reason := runJSONSchemaFixtures()
	if state != statePass {
		t.Fatalf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunJSONSchemaWrapDrift(t *testing.T) {
	skipUnlessIntegration(t)
	state, reason := runJSONSchemaWrapDrift()
	// If check-jsonschema is missing we get cannot_verify; otherwise the drift
	// should still be present, so we expect fail (conformance failure). Only if
	// the drift is unexpectedly fixed do we get pass.
	if !hasTool("check-jsonschema") {
		if state != stateCannotVerify {
			t.Errorf("expected stateCannotVerify when check-jsonschema missing, got %s: %s", state, reason)
		}
		return
	}
	if state != stateFail && state != statePass {
		t.Errorf("unexpected state %s: %s", state, reason)
	}
	if reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestRunOPAPolicy(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("opa") {
		t.Skip("opa not in PATH")
	}
	state, reason := runOPAPolicy()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunOPANegativeFixture(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("opa") {
		t.Skip("opa not in PATH")
	}
	state, reason := runOPANegativeFixture()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunOPANegativeTraceID(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("opa") {
		t.Skip("opa not in PATH")
	}
	state, reason := runOPANegativeTraceID()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunOPANegativeProvenance(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("opa") {
		t.Skip("opa not in PATH")
	}
	state, reason := runOPANegativeProvenance()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunCUEImport(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("cue") {
		t.Skip("cue not in PATH")
	}
	state, reason := runCUEImport()
	if state != statePass {
		t.Errorf("expected statePass, got %s: %s", state, reason)
	}
}

func TestRunInTotoWrap(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("in-toto-run") {
		t.Skip("in-toto-run not in PATH")
	}
	state, reason := runInTotoWrap()
	if state != stateCannotVerify {
		t.Errorf("expected stateCannotVerify, got %s: %s", state, reason)
	}
}

func TestRunCosignLocalSign(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("cosign") {
		t.Skip("cosign not in PATH")
	}
	state, reason := runCosignLocalSign()
	if state != stateCannotVerify {
		t.Errorf("expected stateCannotVerify, got %s: %s", state, reason)
	}
}

func TestRunSLSANegative(t *testing.T) {
	skipUnlessIntegration(t)
	if !hasTool("slsa-verifier") {
		t.Skip("slsa-verifier not in PATH")
	}
	state, reason := runSLSANegative()
	if state != stateCannotVerify {
		t.Errorf("expected stateCannotVerify, got %s: %s", state, reason)
	}
}

func TestWrapRunJSONDriftFixture(t *testing.T) {
	// Structural evidence that the captured run.json does not conform to
	// flight-recorder-run.schema.json. This test checks the frozen snapshot
	// without invoking an external schema validator.
	data, err := os.ReadFile("../../examples/flight-recorder/wrap-output-drift/run.json")
	if err != nil {
		t.Fatalf("read drift fixture: %v", err)
	}
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatalf("expected run.json to be a JSON object: %v", err)
	}
	// The schema requires schema_version == "1.0.0" and many required fields.
	required := []string{
		"schema_version", "run_id", "profile", "trust_scope", "artifact_role",
		"created_at", "source_summary", "task_summary", "model_summary",
		"harness_summary", "evidence_retention_summary", "verifier_states",
	}
	missing := 0
	for _, key := range required {
		if _, ok := obj[key]; !ok {
			missing++
		}
	}
	if missing == 0 {
		t.Fatal("expected frozen run.json to be missing required schema fields, but all were present")
	}
	// Specifically, schema_version must be "1.0.0"; the captured manifest uses
	// a different version string.
	if sv, ok := obj["schema_version"].(string); ok && sv == "1.0.0" {
		t.Fatal("expected schema_version mismatch in frozen run.json")
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
