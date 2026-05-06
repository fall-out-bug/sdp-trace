package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestWrapVerifyExplainMissingEvidenceFlow(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	runDir := filepath.Join(t.TempDir(), "run")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"wrap",
		"--name", "query-check",
		"--output-dir", runDir,
		"--", echo, "hi",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("wrap exit: %d err=%s", exit, errOut.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"verify", runDir}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("verify exit: %d err=%s", exit, errOut.String())
	}
	var result struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("verify payload: %v", err)
	}
	if result.Result != string(trace.VerdictObserved) {
		t.Fatalf("result = %s, expected observed", result.Result)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"query", "--query", "missing-evidence", runDir}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("query exit: %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "\"rows\"") {
		t.Fatalf("unexpected query payload: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"explain", runDir}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("explain exit: %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "run_id:") {
		t.Fatalf("missing explain content: %s", out.String())
	}
}

func TestDryRunOutputsSimulation(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	contractPath := writeTestContract(t, context.Background(), t.TempDir())
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"dry-run",
		"--contract", contractPath,
		"--", echo, "hi",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("dry-run exit: %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "simulation") {
		t.Fatalf("expected simulation output, got %s", out.String())
	}
}

func TestRunRequiresTaskAndRecordsRun(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	runDir := filepath.Join(t.TempDir(), "run")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"run",
		"--task", "task-1",
		"--use-default-contract",
		"--output-dir", runDir,
		"--", echo, "ok",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run exit: %d err=%s", exit, errOut.String())
	}
	if _, err := os.Stat(filepath.Join(runDir, "run.json")); err != nil {
		t.Fatalf("run manifest missing: %v", err)
	}
}

func TestWrapPropagatesNonZeroExitCode(t *testing.T) {
	falseCmd := mustFindCommand(t, "false")
	runDir := filepath.Join(t.TempDir(), "run")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"wrap",
		"--name", "failing",
		"--output-dir", runDir,
		"--", falseCmd,
	}, &out, &errOut)
	if exit == 0 {
		t.Fatalf("expected non-zero exit for failing command")
	}
	var manifest trace.RunManifest
	raw, err := os.ReadFile(filepath.Join(runDir, "run.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if manifest.ClosureState != trace.ClosureStateCommandFailure {
		t.Fatalf("closure state = %s", manifest.ClosureState)
	}
}

func TestValidateFixtures(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrap(t, filepath.Join(root, "run-a"), echo, "a")
	runAndWrap(t, filepath.Join(root, "run-b"), echo, "b")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"validate-fixtures", root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate-fixtures exit: %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "run-a") || !strings.Contains(out.String(), "run-b") {
		t.Fatalf("validate-fixtures output missing runs: %s", out.String())
	}
}

func TestReportAndGateCommands(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "agent")
	runAndWrapNamed(t, filepath.Join(root, "002-verification-run"), "verification-run", echo, "test")
	contractPath := writeGateContract(t, t.TempDir())

	reportDir := filepath.Join(t.TempDir(), "demo-report")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"report", "--out", reportDir, "--contract", contractPath, root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("report exit %d err=%s", exit, errOut.String())
	}
	if _, err := os.Stat(filepath.Join(reportDir, "summary.json")); err != nil {
		t.Fatalf("summary missing: %v", err)
	}

	out.Reset()
	errOut.Reset()
	gatePath := filepath.Join(reportDir, "gate-result.json")
	exit = run([]string{"gate", "--out", gatePath, "--contract", contractPath, root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"audit_grade_gate": "cannot_verify"`) {
		t.Fatalf("gate output missing audit posture: %s", out.String())
	}
}

func TestWitnessCommandMissingCIIdentityCannotVerify(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "agent")

	outPath := filepath.Join(t.TempDir(), "ci-witness.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"witness", "--kind", "github-actions", "--out", outPath, root}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("expected cannot_verify exit, got %d stderr=%s", exit, errOut.String())
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	if !strings.Contains(string(raw), `"status": "cannot_verify"`) {
		t.Fatalf("witness did not record cannot_verify: %s", string(raw))
	}
}

func TestGateCommandAcceptsWitness(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "agent")
	runAndWrapNamed(t, filepath.Join(root, "002-verification-run"), "verification-run", echo, "test")
	contractPath := writeGateContract(t, t.TempDir())
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

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "--out", filepath.Join(t.TempDir(), "gate-result.json"), "--contract", contractPath, "--witness", witnessPath, root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"ci_witness_gate": "cannot_verify"`) {
		t.Fatalf("gate output missing ci witness posture: %s", out.String())
	}
}

func writeGateContract(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "gate-contract.json")
	contract := map[string]any{
		"contract_id": "cli-contract-driven-gate",
		"version":     "sdp-trace-event.v1",
		"required_events": []string{
			"recorder_attached",
			"run_started",
			"command_started",
			"command_finished",
			"run_closed",
		},
		"required_evidence": []map[string]any{
			{
				"id":             "agent_session_observed",
				"event_type":     "command_started",
				"payload_field":  "wrapper_name",
				"payload_equals": "agent-session",
			},
			{
				"id":             "verification_run_observed",
				"event_type":     "command_started",
				"payload_field":  "wrapper_name",
				"payload_equals": "verification-run",
			},
		},
	}
	payload, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal gate contract: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write gate contract: %v", err)
	}
	return path
}

func TestReportRequiresOut(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"report", t.TempDir()}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestFlagSetParsesEndOfFlags(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")
	if err := flags.parse([]string{"--name", "demo", "--", "echo", "hi"}); err != nil {
		t.Fatalf("parse returned error: %v", err)
	}
	if got, want := flags.rest(), []string{"echo", "hi"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected rest %v, got %v", want, got)
	}
}

func TestValidateFixturesHonorsExpectedFailure(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runDir := filepath.Join(root, "tamper-negative")
	runAndWrap(t, runDir, echo, "bad")
	eventPath := filepath.Join(runDir, "events", "000003-command_finished.json")
	var event trace.Event
	raw, err := os.ReadFile(eventPath)
	if err != nil {
		t.Fatalf("read event: %v", err)
	}
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	event.EventHash = "deadbeef"
	mutated, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("marshal mutated event: %v", err)
	}
	if err := os.WriteFile(eventPath, mutated, 0o644); err != nil {
		t.Fatalf("write mutated event: %v", err)
	}
	expectations := []byte(`{"tamper-negative":"fail"}` + "\n")
	if err := os.WriteFile(filepath.Join(root, "fixture-expectations.json"), expectations, 0o644); err != nil {
		t.Fatalf("write fixture expectation: %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"validate-fixtures", root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate-fixtures exit: %d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func TestValidateFixturesRejectsUnexpectedFailure(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runDir := filepath.Join(root, "tamper-negative")
	runAndWrap(t, runDir, echo, "bad")
	eventPath := filepath.Join(runDir, "events", "000003-command_finished.json")
	var event trace.Event
	raw, err := os.ReadFile(eventPath)
	if err != nil {
		t.Fatalf("read event: %v", err)
	}
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	event.EventHash = "deadbeef"
	mutated, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("marshal mutated event: %v", err)
	}
	if err := os.WriteFile(eventPath, mutated, 0o644); err != nil {
		t.Fatalf("write mutated event: %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"validate-fixtures", root}, &out, &errOut)
	if exit == 0 {
		t.Fatalf("validate-fixtures should reject unexpected failure; out=%s", out.String())
	}
}

func TestFlagSetRejectsMissingStringValueBeforeAnotherFlag(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")
	flags.setString("contract", "")
	if err := flags.parse([]string{"--name", "--contract", "contract.json", "--", "echo", "hi"}); err == nil {
		t.Fatalf("expected missing value error")
	}
}

func TestFlagSetRejectsUnknownFlags(t *testing.T) {
	flags := &flagSet{name: "wrap"}
	flags.setString("name", "")
	if err := flags.parse([]string{"--outpt-dir", "x", "--", "echo", "hi"}); err == nil {
		t.Fatalf("expected unknown flag error")
	}
}

func TestFlagSetParsesBooleanLiteral(t *testing.T) {
	flags := &flagSet{name: "dry-run"}
	flags.setBool("use-default-contract", true)
	if err := flags.parse([]string{"--use-default-contract", "false", "--", "echo", "hi"}); err != nil {
		t.Fatalf("parse returned error: %v", err)
	}
	if flags.boolValue("use-default-contract") {
		t.Fatalf("expected false")
	}
}

func runAndWrap(t *testing.T, runDir string, commandPath string, args ...string) {
	t.Helper()
	runAndWrapNamed(t, runDir, "fixture", commandPath, args...)
}

func runAndWrapNamed(t *testing.T, runDir string, wrapperName string, commandPath string, args ...string) {
	t.Helper()
	var out bytes.Buffer
	var errOut bytes.Buffer
	command := append([]string{
		"wrap",
		"--name", wrapperName,
		"--output-dir", runDir,
		"--", commandPath,
	}, args...)
	exit := run(command, &out, &errOut)
	if exit != 0 {
		t.Fatalf("wrap exit %d err=%s", exit, errOut.String())
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

func writeTestContract(t *testing.T, _ context.Context, dir string) string {
	t.Helper()
	contractPath := filepath.Join(dir, "contract.json")
	contract := trace.DefaultContract
	payload, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal contract: %v", err)
	}
	if err := os.WriteFile(contractPath, payload, 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}
	return contractPath
}
