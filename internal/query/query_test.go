package query

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/recorder"
)

func TestMissingEvidenceQueryMatchesVerifier(t *testing.T) {
	tempDir := t.TempDir()
	contractPath := filepath.Join(tempDir, "contract.json")
	contract := map[string]any{
		"contract_id":     "query-missing",
		"version":         "v1",
		"required_events": []string{"recorder_attached", "run_started", "command_started", "command_finished", "run_closed", "file_mutation_observed"},
	}
	contractJSON, _ := json.Marshal(contract)
	if err := osWrite(contractPath, contractJSON); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	runDir := runAndCaptureQuery(t, []string{"q"}, contractPath)
	payload, err := MissingEvidence(runDir)
	if err != nil {
		t.Fatalf("query missing evidence: %v", err)
	}
	var table map[string]any
	if err := json.Unmarshal(payload, &table); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if _, ok := table["rows"]; !ok {
		t.Fatalf("expected rows in payload")
	}
}

func TestCaptureDepthQueryExposesReadOnlyFacts(t *testing.T) {
	run := adaptercapture.ValidTestInput().Run
	run.UnverifiedTaskExpanded = true
	run.TaskSupersessionCount = 2
	run.UnsupportedEventTypes = []string{"tool_call"}
	filtered := []adaptercapture.AdapterEvent{}
	for _, event := range run.AdapterEvents {
		if event.EventType != "tool_call" {
			filtered = append(filtered, event)
		}
	}
	run.AdapterEvents = filtered

	runDir := t.TempDir()
	payload, err := json.Marshal(run)
	if err != nil {
		t.Fatalf("marshal run: %v", err)
	}
	if err := osWrite(filepath.Join(runDir, "run.json"), payload); err != nil {
		t.Fatalf("write run: %v", err)
	}

	queryPayload, err := CaptureDepth(runDir)
	if err != nil {
		t.Fatalf("query capture depth: %v", err)
	}
	var summary CaptureDepthSummary
	if err := json.Unmarshal(queryPayload, &summary); err != nil {
		t.Fatalf("decode capture depth: %v", err)
	}
	if summary.Query != QueryCaptureDepth || summary.TopLevelAssessment != "not_emitted_for_query" {
		t.Fatalf("summary identity = %+v", summary)
	}
	if !summary.UnverifiedTaskExpanded || summary.TaskSupersessionCount != 2 {
		t.Fatalf("task facts = %+v", summary)
	}
	if len(summary.MissingAdapterEvents) != 1 || summary.MissingAdapterEvents[0] != "tool_call" {
		t.Fatalf("missing events = %+v", summary.MissingAdapterEvents)
	}
	if len(summary.UnsupportedObservers) != 1 || summary.UnsupportedObservers[0] != "tool_call" {
		t.Fatalf("unsupported observers = %+v", summary.UnsupportedObservers)
	}
	if len(summary.UnverifiedClaims) == 0 {
		t.Fatalf("expected unverified claims in summary")
	}
}

func runAndCaptureQuery(t *testing.T, command []string, contract string) string {
	t.Helper()
	echo := mustFindCommandQuery(t, "echo")
	ctx := context.Background()
	result, err := recorder.Run(ctx, recorder.RecorderOptions{
		WrapperName:        "query",
		UseDefaultContract: false,
		ContractPath:       contract,
		OutputDir:          t.TempDir(),
		Command:            append([]string{echo}, command...),
	})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	return result.RunDir
}

func mustFindCommandQuery(t *testing.T, name string) string {
	t.Helper()
	path, err := execLookPath(name)
	if err != nil {
		t.Skipf("%s not available", name)
	}
	return path
}

func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func osWrite(path string, value []byte) error {
	return os.WriteFile(path, value, 0o644)
}
