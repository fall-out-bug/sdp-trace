package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

func TestAdapterCaptureAssessRequiresInputsWithoutWriting(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"assess", "--profile", "adapter-capture", "--out", outPath}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("adapter capture missing inputs exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Fatalf("adapter capture wrote artifact despite usage error")
	}
}

func TestAdapterCaptureAssessPassesExplainsAndQueries(t *testing.T) {
	root := t.TempDir()
	runDir := writeAdapterCaptureFixtureInputs(t, root, adaptercapture.ValidTestInput().Run)
	outPath := filepath.Join(root, "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "adapter-capture",
		"--out", outPath,
		"--run", runDir,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("adapter capture assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	assertNoAdapterLeak(t, out.String())
	var result adaptercapture.AssessmentResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("assessment payload: %v", err)
	}
	if result.AdapterCaptureAssessment != adaptercapture.StatePass || result.SchemaVersion != adaptercapture.SchemaVersion {
		t.Fatalf("assessment result = %+v", result)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"assess", "explain", "--assessment-result", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("adapter capture explain exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "Adapter capture assessment: pass") ||
		!strings.Contains(out.String(), "Adapter condition run_binding_established: pass") {
		t.Fatalf("adapter capture explain missing fields: %s", out.String())
	}
	assertNoAdapterLeak(t, out.String())

	out.Reset()
	errOut.Reset()
	exit = run([]string{"query", "--query", "capture-depth", runDir}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("capture-depth query exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"query": "capture-depth"`) ||
		!strings.Contains(out.String(), `"top_level_assessment": "not_emitted_for_query"`) ||
		!strings.Contains(out.String(), `"task_supersession_count": 0`) ||
		!strings.Contains(out.String(), `"unverified_task_expanded": false`) {
		t.Fatalf("capture-depth query missing summary: %s", out.String())
	}
	assertNoAdapterLeak(t, out.String())
}

func TestAdapterCaptureAssessRejectsAgentReportedExecutedTestAndLeakyRefs(t *testing.T) {
	root := t.TempDir()
	input := adaptercapture.ValidTestInput()
	for i := range input.Run.AdapterEvents {
		if input.Run.AdapterEvents[i].EventType == "test_observed" {
			input.Run.AdapterEvents[i].TestProvenance = "agent_reported"
			input.Run.AdapterEvents[i].ExecutedEvidenceClaimed = true
		}
	}
	input.Run.ProviderRefs = []adaptercapture.ProviderRef{{ReviewRef: "https://review.invalid/7?token=secret-token"}}
	runDir := writeAdapterCaptureFixtureInputs(t, root, input.Run)
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "adapter-capture",
		"--out", filepath.Join(root, "assessment.json"),
		"--run", runDir,
	}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("adapter capture assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"reason_code": "agent_reported_test_not_executed"`) ||
		!strings.Contains(out.String(), `"reason_code": "provider_ref_contains_secret"`) {
		t.Fatalf("adapter capture missing fail reasons: %s", out.String())
	}
	assertNoAdapterLeak(t, out.String())
}

func TestAdapterCapturePreviewDoesNotWriteOrLeak(t *testing.T) {
	root := t.TempDir()
	outPath := filepath.Join(root, "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"assess", "preview", "--profile", "adapter-capture", "--out", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("adapter capture preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"selected_profile": "adapter_capture"`) ||
		!strings.Contains(out.String(), `"claim": "preview is read-only and does not emit an adapter capture verdict"`) {
		t.Fatalf("adapter capture preview missing fields: %s", out.String())
	}
	assertNoAdapterLeak(t, out.String())
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Fatalf("adapter capture preview wrote artifact")
	}
}

func writeAdapterCaptureFixtureInputs(t *testing.T, dir string, run adaptercapture.RunEvidence) string {
	t.Helper()
	runDir := filepath.Join(dir, "run")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run: %v", err)
	}
	writeTestJSON(t, filepath.Join(runDir, "run.json"), run)
	return runDir
}

func assertNoAdapterLeak(t *testing.T, text string) {
	t.Helper()
	for _, marker := range []string{
		"secret-token",
		"raw prompt",
		"raw response",
		"raw_review_body",
		"tool_input_body",
		"tool_output_body",
		"model_request_payload",
		"model_response_payload",
		"gateway_evidence_ref",
		"--password",
		"Bearer ",
		"access_token=",
		"credential=",
		"oidc_token",
		"session_id=",
		"adapter_config_raw",
	} {
		if strings.Contains(text, marker) {
			t.Fatalf("output leaked sensitive marker %q: %s", marker, text)
		}
	}
}
