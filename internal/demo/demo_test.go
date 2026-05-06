package demo

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

func TestReportAndGatePassForObservedDemoRuns(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	runCommand(t, filepath.Join(root, "002-verification-run"), "verification-run", echo, "tests")
	contractPath := writeDemoContract(t, t.TempDir())

	reportDir := filepath.Join(t.TempDir(), "report")
	artifacts, err := WriteReport(root, reportDir, contractPath)
	if err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	if artifacts.Summary.RunCount != 2 {
		t.Fatalf("run count = %d", artifacts.Summary.RunCount)
	}
	if artifacts.Summary.AuditGrade {
		t.Fatalf("local report must not be audit grade")
	}
	for _, name := range []string{"summary.json", "evidence-table.json", "missing-telemetry.json", "timeline.md"} {
		if _, err := os.Stat(filepath.Join(reportDir, name)); err != nil {
			t.Fatalf("missing report artifact %s: %v", name, err)
		}
	}
	timeline, err := os.ReadFile(filepath.Join(reportDir, "timeline.md"))
	if err != nil {
		t.Fatalf("read timeline: %v", err)
	}
	if !strings.Contains(string(timeline), "agent_session_observed") || !strings.Contains(string(timeline), "verification_run_observed") {
		t.Fatalf("timeline missing classified runs:\n%s", string(timeline))
	}

	gatePath := filepath.Join(t.TempDir(), "gate-result.json")
	gate, err := WriteGate(root, gatePath, contractPath)
	if err != nil {
		t.Fatalf("WriteGate: %v", err)
	}
	if gate.LocalGate != GatePass {
		t.Fatalf("local gate = %s reasons=%v", gate.LocalGate, gate.Reasons)
	}
	if gate.AuditGradeGate != GateCannotVerify {
		t.Fatalf("audit gate = %s", gate.AuditGradeGate)
	}
	if !contains(gate.MissingAuditEvidence, "ci_oidc_witness") {
		t.Fatalf("missing audit evidence not reported: %v", gate.MissingAuditEvidence)
	}
	if contains(gate.RequiredEvidence, "all_runs_observed") {
		t.Fatalf("gate conditions leaked into required evidence: %v", gate.RequiredEvidence)
	}
	if !contains(gate.GateConditions, "all_runs_observed") {
		t.Fatalf("gate conditions missing: %v", gate.GateConditions)
	}
}

func TestReportAcceptsSingleRunDirectory(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	runDir := filepath.Join(t.TempDir(), "single-run")
	runCommand(t, runDir, "agent-session", echo, "hello")

	artifacts, err := WriteReport(runDir, filepath.Join(t.TempDir(), "report"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	if artifacts.Summary.RunCount != 1 {
		t.Fatalf("run count = %d", artifacts.Summary.RunCount)
	}
}

func TestGateFailsForTamperedRun(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	tampered := filepath.Join(root, "002-verification-run")
	runCommand(t, tampered, "verification-run", echo, "tests")

	eventPath := filepath.Join(tampered, "events", "000003-command_finished.json")
	var event map[string]any
	raw, err := os.ReadFile(eventPath)
	if err != nil {
		t.Fatalf("read event: %v", err)
	}
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	event["event_hash"] = "deadbeef"
	mutated, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	if err := os.WriteFile(eventPath, mutated, 0o644); err != nil {
		t.Fatalf("write event: %v", err)
	}

	gate, err := WriteGate(root, filepath.Join(t.TempDir(), "gate-result.json"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteGate: %v", err)
	}
	if gate.LocalGate != GateFail {
		t.Fatalf("expected fail gate, got %s", gate.LocalGate)
	}
}

func TestGateFailsForNotAssessedRun(t *testing.T) {
	contract := traceContractForTest()
	gate := EvaluateGate([]RunRow{
		{
			Name:         "agent-session",
			Kind:         "agent_session_observed",
			Result:       "observed",
			ClosureState: "completed",
		},
		{
			Name:         "verification-run",
			Kind:         "verification_run_observed",
			Result:       "observed",
			ClosureState: "completed",
		},
		{
			Name:         "not-assessed-run",
			Kind:         "unmatched",
			Result:       "not_assessed",
			ClosureState: "completed",
		},
	}, contract)
	if gate.LocalGate != GateFail {
		t.Fatalf("expected not_assessed run to fail local gate, got %s", gate.LocalGate)
	}
}

func TestGateUsesPassingCIWitness(t *testing.T) {
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "source": {"repository": "org/repo", "ref": "refs/heads/main", "commit_sha": "abc123"},
	  "ci": {"provider": "github-actions", "server_url": "https://github.com", "workflow": "sdp-trace", "job": "test", "run_id": "42", "run_attempt": "1", "actor": "octocat"},
	  "oidc": {"issuer": "https://token.actions.githubusercontent.com", "subject": "repo:org/repo:ref:refs/heads/main", "audience": "sdp-trace", "repository": "org/repo", "ref": "refs/heads/main", "sha": "abc123"},
	  "run_artifacts": [],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	gate := EvaluateGateWithWitness([]RunRow{
		{Name: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}, traceContractForTest(), witnessPath)
	if gate.CIWitnessGate != GatePass {
		t.Fatalf("ci witness gate = %s reasons=%v", gate.CIWitnessGate, gate.Reasons)
	}
	if contains(gate.MissingAuditEvidence, "ci_oidc_witness") {
		t.Fatalf("ci witness still listed as missing: %v", gate.MissingAuditEvidence)
	}
	if !contains(gate.MissingAuditEvidence, "external_witness_checkpoint") {
		t.Fatalf("external witness missing state lost: %v", gate.MissingAuditEvidence)
	}
	if gate.AuditGradeGate != GateCannotVerify {
		t.Fatalf("audit grade gate = %s", gate.AuditGradeGate)
	}
}

func TestGateUsesCannotVerifyCIWitness(t *testing.T) {
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "cannot_verify",
	  "trust_scope": "local_observed",
	  "reason": "missing_ci_oidc",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "missing_identity_fields": ["ACTIONS_ID_TOKEN_REQUEST_URL"],
	  "source": {"repository": "", "ref": "", "commit_sha": ""},
	  "ci": {"provider": "github-actions", "server_url": "", "workflow": "", "job": "", "run_id": "", "run_attempt": "", "actor": ""},
	  "run_artifacts": [],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	gate := EvaluateGateWithWitness([]RunRow{
		{Name: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}, traceContractForTest(), witnessPath)
	if gate.CIWitnessGate != GateCannotVerify {
		t.Fatalf("ci witness gate = %s", gate.CIWitnessGate)
	}
	if !contains(gate.MissingAuditEvidence, "ci_oidc_witness") {
		t.Fatalf("ci witness missing state lost: %v", gate.MissingAuditEvidence)
	}
	if !contains(gate.MissingAuditEvidence, "external_witness_checkpoint") {
		t.Fatalf("external witness missing state lost: %v", gate.MissingAuditEvidence)
	}
}

func TestMalformedRunAppearsAsCannotVerify(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	broken := filepath.Join(root, "002-verification-run")
	runCommand(t, broken, "verification-run", echo, "tests")
	if err := os.WriteFile(filepath.Join(broken, "run.json"), []byte("{"), 0o644); err != nil {
		t.Fatalf("break run manifest: %v", err)
	}

	gate, err := WriteGate(root, filepath.Join(t.TempDir(), "gate-result.json"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteGate: %v", err)
	}
	if gate.LocalGate != GateFail {
		t.Fatalf("expected fail gate, got %s", gate.LocalGate)
	}
	foundCannotVerify := false
	for _, row := range gate.Runs {
		if row.Name == "002-verification-run" && row.Result == "cannot_verify" {
			foundCannotVerify = true
		}
	}
	if !foundCannotVerify {
		t.Fatalf("malformed run was not reported as cannot_verify: %+v", gate.Runs)
	}
}

func TestVerifierArtifactWriteFailureAppearsAsCannotVerify(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	readOnlyRun := filepath.Join(root, "002-verification-run")
	runCommand(t, readOnlyRun, "verification-run", echo, "tests")
	verifierDir := filepath.Join(readOnlyRun, "verifier")
	if err := os.Chmod(verifierDir, 0o555); err != nil {
		t.Fatalf("chmod verifier dir: %v", err)
	}
	defer func() {
		_ = os.Chmod(verifierDir, 0o755)
	}()

	artifacts, err := WriteReport(root, filepath.Join(t.TempDir(), "report"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteReport should still produce report: %v", err)
	}
	if artifacts.Summary.CannotVerifyCount != 1 {
		t.Fatalf("cannot verify count = %d", artifacts.Summary.CannotVerifyCount)
	}
}

func TestReportRequiresOutputDirectory(t *testing.T) {
	if _, err := WriteReport(t.TempDir(), "", ""); err == nil {
		t.Fatalf("expected missing output directory error")
	}
}

func TestGateRequiresOutputFile(t *testing.T) {
	if _, err := WriteGate(t.TempDir(), "", ""); err == nil {
		t.Fatalf("expected missing output file error")
	}
}

func writeDemoContract(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "contract.json")
	contract := traceContractForTest()
	payload, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal contract: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}
	return path
}

func traceContractForTest() trace.Contract {
	return trace.Contract{
		ContractID: "demo-contract-driven-test",
		Version:    "sdp-trace-event.v1",
		RequiredEvents: []string{
			string(trace.EventRecorderAttached),
			string(trace.EventRunStarted),
			string(trace.EventCommandStarted),
			string(trace.EventCommandFinished),
			string(trace.EventRunClosed),
		},
		RequiredEvidence: []trace.EvidenceRequirement{
			{
				ID:            "agent_session_observed",
				EventType:     string(trace.EventCommandStarted),
				PayloadField:  "wrapper_name",
				PayloadEquals: "agent-session",
			},
			{
				ID:            "verification_run_observed",
				EventType:     string(trace.EventCommandStarted),
				PayloadField:  "wrapper_name",
				PayloadEquals: "verification-run",
			},
		},
	}
}

func runCommand(t *testing.T, runDir, wrapperName, command string, args ...string) {
	t.Helper()
	cmd := append([]string{command}, args...)
	result, err := recorder.Run(context.Background(), recorder.RecorderOptions{
		WrapperName:        wrapperName,
		OutputDir:          runDir,
		UseDefaultContract: true,
		Command:            cmd,
	})
	if err != nil {
		t.Fatalf("recorder run: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("recorder exit = %d", result.ExitCode)
	}
}

func mustFindCommand(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s not available", name)
	}
	return path
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
