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
	if !strings.Contains(out.String(), `"writes_artifacts": false`) {
		t.Fatalf("expected no-write posture, got %s", out.String())
	}
	if !strings.Contains(out.String(), string(trace.RetentionModeDigestOnly)) {
		t.Fatalf("expected safe retention modes, got %s", out.String())
	}
	if strings.Contains(out.String(), `"command":`) || strings.Contains(out.String(), `"hi"`) {
		t.Fatalf("dry-run leaked raw command payload: %s", out.String())
	}
}

func TestPreviewOutputsNoWritePlan(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	workDir := t.TempDir()
	t.Chdir(workDir)

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"preview", "--", echo, "hi"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("preview exit: %d err=%s", exit, errOut.String())
	}
	var payload struct {
		Mode                string                      `json:"mode"`
		CommandDescriptor   trace.CommandDescriptor     `json:"command_descriptor"`
		Boundaries          []previewBoundary           `json:"boundaries"`
		OfflineImplications []previewOfflineImplication `json:"offline_implications"`
		WritesArtifacts     bool                        `json:"writes_artifacts"`
		SafeRetentionModes  []string                    `json:"safe_retention_modes"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("preview payload: %v", err)
	}
	if payload.Mode != "preview" {
		t.Fatalf("mode = %s", payload.Mode)
	}
	if payload.WritesArtifacts {
		t.Fatalf("preview must not write artifacts")
	}
	if payload.CommandDescriptor.Argc != 2 {
		t.Fatalf("preview command argc = %d", payload.CommandDescriptor.Argc)
	}
	if payload.CommandDescriptor.Retention.Mode != trace.RetentionModeDigestOnly {
		t.Fatalf("preview command retention = %s", payload.CommandDescriptor.Retention.Mode)
	}
	if len(payload.SafeRetentionModes) == 0 {
		t.Fatalf("missing safe retention modes")
	}
	if findPreviewBoundary(t, payload.Boundaries, string(trace.ObservationBoundaryAdapterSocket)).State != string(trace.ObservationStateNotIntegrated) {
		t.Fatalf("adapter boundary state missing")
	}
	if findPreviewImplication(t, payload.OfflineImplications, "ci_witnessed").State != string(trace.ObservationStateOfflineDev) {
		t.Fatalf("offline CI implication missing")
	}
	if strings.Contains(out.String(), `"hi"`) || strings.Contains(out.String(), echo) {
		t.Fatalf("preview leaked raw argv: %s", out.String())
	}
	if _, err := os.Stat(filepath.Join(workDir, ".sdp-trace-runs")); !os.IsNotExist(err) {
		t.Fatalf("preview wrote run artifacts or stat failed: %v", err)
	}
}

func TestDoctorReportsOfflineDevAndCannotVerifyCI(t *testing.T) {
	clearCIWitnessEnv(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"doctor"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("doctor exit: %d err=%s", exit, errOut.String())
	}
	var report doctorReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatalf("doctor payload: %v", err)
	}
	if report.Result != "offline_dev" {
		t.Fatalf("result = %s", report.Result)
	}
	check := findDoctorCheck(t, report.ControlPoints, "ci_witness_prerequisites")
	if check.State != string(trace.VerdictCannotVerify) {
		t.Fatalf("ci witness state = %s", check.State)
	}
	if len(check.Missing) == 0 {
		t.Fatalf("expected missing CI witness fields")
	}
	if findDoctorCheck(t, report.Environment, "offline_development").State != "offline_dev" {
		t.Fatalf("offline development state missing")
	}
	if findDoctorCheck(t, report.ControlPoints, "output_directory").State != "pass" {
		t.Fatalf("output directory check missing")
	}
	if findDoctorCheck(t, report.ControlPoints, "report_directory").State != "pass" {
		t.Fatalf("report directory check missing")
	}
	if findDoctorCheck(t, report.ControlPoints, "expected_evidence_references").State != "pass" {
		t.Fatalf("expected evidence check missing")
	}
}

func TestDoctorReportsContractLoadFailureCannotVerify(t *testing.T) {
	clearCIWitnessEnv(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	missingContract := filepath.Join(t.TempDir(), "missing-contract.json")
	exit := run([]string{"doctor", "--contract", missingContract}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("doctor exit: %d err=%s", exit, errOut.String())
	}
	var report doctorReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatalf("doctor payload: %v", err)
	}
	if report.Result != string(trace.VerdictCannotVerify) {
		t.Fatalf("result = %s", report.Result)
	}
	check := findDoctorCheck(t, report.ControlPoints, "contract")
	if check.State != string(trace.VerdictCannotVerify) {
		t.Fatalf("contract state = %s", check.State)
	}
}

func TestDoctorReportsUnwritableOutputDirectoryCannotVerify(t *testing.T) {
	clearCIWitnessEnv(t)
	filePath := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(filePath, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"doctor", "--output-dir", filePath}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("doctor exit: %d err=%s", exit, errOut.String())
	}
	var report doctorReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatalf("doctor payload: %v", err)
	}
	check := findDoctorCheck(t, report.ControlPoints, "output_directory")
	if check.State != string(trace.VerdictCannotVerify) {
		t.Fatalf("output directory state = %s", check.State)
	}
}

func TestDoctorReportsUnsupportedExpectedEvidenceCannotVerify(t *testing.T) {
	clearCIWitnessEnv(t)
	contractPath := filepath.Join(t.TempDir(), "contract.json")
	if err := os.WriteFile(contractPath, []byte(`{
	  "contract_id": "unsupported-contract",
	  "version": "test",
	  "required_events": ["recorder_attached", "model_call_observed"]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"doctor", "--contract", contractPath}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("doctor exit: %d err=%s", exit, errOut.String())
	}
	var report doctorReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatalf("doctor payload: %v", err)
	}
	check := findDoctorCheck(t, report.ControlPoints, "expected_evidence_references")
	if check.State != string(trace.VerdictCannotVerify) {
		t.Fatalf("expected evidence state = %s", check.State)
	}
	if len(check.Missing) == 0 {
		t.Fatalf("expected missing unsupported event references")
	}
}

func TestUsageMentionsDoctorAndPreview(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"help"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("help exit: %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "sdp-trace preview") || !strings.Contains(out.String(), "sdp-trace doctor") {
		t.Fatalf("usage missing new commands: %s", out.String())
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

func clearCIWitnessEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN",
		"ACTIONS_ID_TOKEN_REQUEST_URL",
		"GITHUB_ACTIONS",
		"GITHUB_ACTOR",
		"GITHUB_JOB",
		"GITHUB_REF",
		"GITHUB_REPOSITORY",
		"GITHUB_RUN_ATTEMPT",
		"GITHUB_RUN_ID",
		"GITHUB_SERVER_URL",
		"GITHUB_SHA",
		"GITHUB_WORKFLOW",
	} {
		t.Setenv(key, "")
	}
}

func findDoctorCheck(t *testing.T, checks []doctorCheck, id string) doctorCheck {
	t.Helper()
	for _, check := range checks {
		if check.ID == id {
			return check
		}
	}
	t.Fatalf("doctor check %s not found in %#v", id, checks)
	return doctorCheck{}
}

func findPreviewBoundary(t *testing.T, boundaries []previewBoundary, id string) previewBoundary {
	t.Helper()
	for _, boundary := range boundaries {
		if boundary.Boundary == id {
			return boundary
		}
	}
	t.Fatalf("preview boundary %s not found in %#v", id, boundaries)
	return previewBoundary{}
}

func findPreviewImplication(t *testing.T, implications []previewOfflineImplication, requirement string) previewOfflineImplication {
	t.Helper()
	for _, implication := range implications {
		if implication.Requirement == requirement {
			return implication
		}
	}
	t.Fatalf("preview implication %s not found in %#v", requirement, implications)
	return previewOfflineImplication{}
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
