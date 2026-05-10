package verifier

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/recorder"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestVerifyReturnsObservedForUntamperedRun(t *testing.T) {
	runDir := runAndCapture(t, []string{"ok"}, true, "")

	result, table, audit, err := VerifyRun(runDir)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if result.Result != trace.VerdictObserved {
		t.Fatalf("expected observed, got %s", result.Result)
	}
	if len(table.Rows) != 0 {
		t.Fatalf("expected empty missing-evidence table, got %d", len(table.Rows))
	}
	if audit != nil {
		t.Fatalf("expected no integrity audit")
	}
}

func TestVerifyDetectsTamperedChain(t *testing.T) {
	runDir := runAndCapture(t, []string{"tamper"}, true, "")
	eventPath := filepath.Join(runDir, "events", "000003-command_finished.json")
	var event trace.Event
	data, err := os.ReadFile(eventPath)
	if err != nil {
		t.Fatalf("read event: %v", err)
	}
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	event.EventHash = "sha256:deadbeef"
	mutated, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("marshal mutated event: %v", err)
	}
	if err := os.WriteFile(eventPath, mutated, 0o644); err != nil {
		t.Fatalf("write mutated event: %v", err)
	}
	result, _, audit, err := VerifyRun(runDir)
	if err != nil {
		t.Fatalf("verify unexpected error: %v", err)
	}
	if result.Result != trace.VerdictFail {
		t.Fatalf("expected fail, got %s", result.Result)
	}
	if audit == nil || audit.Issue == "" {
		t.Fatalf("expected integrity audit for tamper")
	}
}

func TestVerifyReportsMissingEvidenceRows(t *testing.T) {
	tempDir := t.TempDir()
	contractPath := filepath.Join(tempDir, "contract.json")
	contract := trace.Contract{
		ContractID:     "missing-evidence-test",
		Version:        "v1",
		RequiredEvents: []string{"recorder_attached", "run_started", "command_started", "command_finished", "run_closed", "test_observed"},
	}
	if err := writeJSONFile(contractPath, contract); err != nil {
		t.Fatalf("write contract: %v", err)
	}
	runDir := runAndCapture(t, []string{"missing"}, true, contractPath)
	result, table, _, err := VerifyRun(runDir)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if result.Result != trace.VerdictNotAssessed {
		t.Fatalf("expected not_assessed, got %s", result.Result)
	}
	if len(table.Rows) != 1 {
		t.Fatalf("expected one missing row, got %d", len(table.Rows))
	}
	if table.Rows[0].ExpectedEvent != "test_observed" {
		t.Fatalf("expected missing event test_observed, got %s", table.Rows[0].ExpectedEvent)
	}
}

func TestVerifyCannotVerifyMissingManifest(t *testing.T) {
	runDir := t.TempDir()
	result, _, audit, err := VerifyRun(runDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result != trace.VerdictCannotVerify {
		t.Fatalf("expected cannot_verify, got %s", result.Result)
	}
	if audit == nil || audit.Issue != "run_manifest_missing" {
		t.Fatalf("expected run_manifest_missing audit, got %#v", audit)
	}
}

func TestVerifyCannotVerifyMissingEventsDirectory(t *testing.T) {
	runDir := t.TempDir()
	manifest := trace.RunManifest{
		SchemaVersion:   trace.SchemaVersion,
		RunID:           "missing-events",
		RecorderVersion: trace.RecorderVersion,
		ContractID:      trace.DefaultContract.ContractID,
		EventCount:      1,
	}
	if err := writeJSONFile(filepath.Join(runDir, "run.json"), manifest); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	result, _, audit, err := VerifyRun(runDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result != trace.VerdictCannotVerify {
		t.Fatalf("expected cannot_verify, got %s", result.Result)
	}
	if audit == nil || audit.Issue != "event_load_failed" {
		t.Fatalf("expected event_load_failed audit, got %#v", audit)
	}
}

func TestExplainRunIncludesMissingEvidence(t *testing.T) {
	tempDir := t.TempDir()
	contractPath := filepath.Join(tempDir, "contract.json")
	contract := trace.Contract{
		ContractID:     "explain-missing-evidence",
		Version:        "v1",
		RequiredEvents: []string{"recorder_attached", "run_started", "command_started", "command_finished", "run_closed", "test_observed"},
	}
	if err := writeJSONFile(contractPath, contract); err != nil {
		t.Fatalf("write contract: %v", err)
	}
	runDir := runAndCapture(t, []string{"explain"}, true, contractPath)

	explanation, err := ExplainRun(runDir)
	if err != nil {
		t.Fatalf("ExplainRun() error = %v", err)
	}
	for _, want := range []string{
		"result: not_assessed",
		"missing_evidence:",
		" - test_observed: missing (required_by_contract)",
		"contract_path: " + contractPath,
	} {
		if !strings.Contains(explanation, want) {
			t.Fatalf("explanation missing %q:\n%s", want, explanation)
		}
	}
}

func TestWriteVerifierArtifactsWritesAuditConditionally(t *testing.T) {
	runDir := t.TempDir()
	result := trace.VerifierResult{
		RunID:  "run-a",
		RunDir: runDir,
		Result: trace.VerdictCannotVerify,
		Reason: "tampered",
	}
	table := trace.MissingEvidenceTable{
		ContractID: "contract-a",
		Rows: []trace.MissingEvidenceRow{{
			ExpectedEvent: "test_observed",
			ObservedState: "missing",
			Reason:        "required event not found",
		}},
	}
	audit := &trace.IntegrityAudit{
		RunID:  "run-a",
		Issue:  "tampered_chain",
		Reason: "hash mismatch",
	}

	if err := WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
		t.Fatalf("WriteVerifierArtifacts() error = %v", err)
	}
	for _, name := range []string{"verifier-result.json", "missing-evidence-table.json", "integrity-audit.json"} {
		if _, err := os.Stat(filepath.Join(runDir, "verifier", name)); err != nil {
			t.Fatalf("expected %s to be written: %v", name, err)
		}
	}

	runDirNoAudit := t.TempDir()
	if err := WriteVerifierArtifacts(runDirNoAudit, result, table, nil); err != nil {
		t.Fatalf("WriteVerifierArtifacts() without audit error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(runDirNoAudit, "verifier", "integrity-audit.json")); !os.IsNotExist(err) {
		t.Fatalf("integrity audit should not be written when nil: %v", err)
	}
}

func runAndCapture(t *testing.T, command []string, useDefault bool, contract string) string {
	t.Helper()
	echo := mustFindCommand(t, "echo")
	ctx := context.Background()
	opts := recorder.RecorderOptions{
		WrapperName:        "test",
		UseDefaultContract: useDefault,
		OutputDir:          t.TempDir(),
		Command:            append([]string{echo}, command...),
	}
	if contract != "" {
		opts.ContractPath = contract
	}
	result, err := recorder.Run(ctx, opts)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	return result.RunDir
}

func mustFindCommand(t *testing.T, name string) string {
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

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
