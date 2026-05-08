package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
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

func TestReleaseProofWritesFailForMissingManifestArtifact(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.invalid")
	runGit(t, repo, "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, "source.txt"), "source\n")
	runGit(t, repo, "add", "source.txt")
	runGit(t, repo, "commit", "-m", "source")
	sourceCommit := strings.TrimSpace(gitOutput(t, repo, "rev-parse", "HEAD"))
	writeFile(t, filepath.Join(repo, "manifest.json"), `{
  "id": "test-manifest",
  "signing_profile": "sdp-trace-signature/sigstore-dsse-keyless-v1",
  "trusted_identity_policy_ref": "policy.json",
  "source_commit": "`+sourceCommit+`",
  "artifacts": [
    {"path": "missing.txt", "kind": "doc", "sha256": "1111111111111111111111111111111111111111111111111111111111111111"}
  ],
  "accountability": {
    "dri": {"identity_ref": "role:dri", "actor_type": "human_role"},
    "approver": {"identity_ref": "role:approver", "actor_type": "human_role"},
    "escalation": {"identity_ref": "role:cto", "actor_type": "human_role"},
    "authority_scope": "contract_release",
    "accountability_claim": "release_approval",
    "approval_ref": "approval",
    "risk_owner": {"identity_ref": "role:risk", "actor_type": "human_role"},
    "line_of_defense": "second"
	  }
	}`)
	runGit(t, repo, "add", "manifest.json")
	runGit(t, repo, "commit", "-m", "manifest")
	chdir(t, repo)
	outPath := filepath.Join(t.TempDir(), "release-proof.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"release-proof",
		"--manifest", "manifest.json",
		"--out", outPath,
	}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("release-proof exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"release_verification_state": "fail"`) ||
		!strings.Contains(out.String(), `"trusted_contract_release": false`) {
		t.Fatalf("release-proof output missing fail boundary: %s", out.String())
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("release proof artifact missing: %v", err)
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
	chdir(t, workDir)

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

func TestInstallRepoObserverDryRunDoesNotRequirePromptCooperation(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init")
	chdir(t, repo)

	outPath := filepath.Join(t.TempDir(), "repo-observer-status.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"install", "repo-observer",
		"--profile", repoobserver.ProfileGithubActionsGitHooksV1,
		"--repository-id", "demo_repo",
		"--out", outPath,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("install dry-run exit: %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "agent_prompt") || !strings.Contains(out.String(), "agent_reported") {
		t.Fatalf("human table must expose agent prompt as non-proof surface: %s", out.String())
	}
	if _, err := os.Stat(filepath.Join(repo, ".githooks", "pre-commit")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote hook file or stat failed: %v", err)
	}
	var status repoobserver.Status
	if err := readJSONFile(outPath, &status); err != nil {
		t.Fatalf("status json: %v", err)
	}
	if status.SchemaVersion != repoobserver.SchemaVersion {
		t.Fatalf("schema version = %s", status.SchemaVersion)
	}
	if status.RepositoryID != "demo_repo" {
		t.Fatalf("repository id = %s", status.RepositoryID)
	}
	if status.RepositoryRootRef != "current_repository" {
		t.Fatalf("repository root ref leaked path: %s", status.RepositoryRootRef)
	}
	if surfaceByID(t, status, repoobserver.SurfaceAgentPrompt).ReasonCode != repoobserver.ReasonAgentReportedNotProof {
		t.Fatalf("agent prompt surface did not mark prompt cooperation as non-proof")
	}
}

func TestInstallRepoObserverWriteAndDoctorProfile(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init")
	chdir(t, repo)

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"install", "repo-observer",
		"--profile", repoobserver.ProfileGithubActionsGitHooksV1,
		"--repository-id", "demo_repo",
		"--write",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("install write exit: %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if got := strings.TrimSpace(gitOutput(t, repo, "config", "--get", "core.hooksPath")); got != ".githooks" {
		t.Fatalf("core.hooksPath = %s", got)
	}
	for _, rel := range []string{
		".sdp-trace/config.json",
		".githooks/pre-commit",
		".githooks/post-commit",
		".githooks/pre-push",
		".github/workflows/sdp-trace-observe.yml",
	} {
		if _, err := os.Stat(filepath.Join(repo, rel)); err != nil {
			t.Fatalf("missing installed file %s: %v", rel, err)
		}
	}
	hookData, err := os.ReadFile(filepath.Join(repo, ".githooks", "pre-commit"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(hookData), "printenv") || strings.Contains(string(hookData), "env >") {
		t.Fatalf("hook must not persist raw environment: %s", string(hookData))
	}

	out.Reset()
	errOut.Reset()
	outPath := filepath.Join(t.TempDir(), "doctor.json")
	exit = run([]string{
		"doctor",
		"--profile", repoobserver.ProfileGithubActionsGitHooksV1,
		"--out", outPath,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("doctor profile exit: %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	var status repoobserver.Status
	if err := readJSONFile(outPath, &status); err != nil {
		t.Fatalf("doctor json: %v", err)
	}
	if status.InstallState != repoobserver.StatePass {
		t.Fatalf("install state = %s", status.InstallState)
	}
	if status.RepositoryID != "demo_repo" {
		t.Fatalf("doctor repository id = %s", status.RepositoryID)
	}
	if status.ProofState != repoobserver.StateNotAssessed {
		t.Fatalf("proof state = %s", status.ProofState)
	}
	if surfaceByID(t, status, repoobserver.SurfaceCIArtifactUpload).TrustScope != repoobserver.ScopeCIUploaded {
		t.Fatalf("ci artifact upload scope not represented")
	}
}

func TestUsageMentionsDoctorPreviewAndRepoObserverInstall(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"help"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("help exit: %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "sdp-trace preview") ||
		!strings.Contains(out.String(), "sdp-trace doctor --profile github-actions-git-hooks-v1") ||
		!strings.Contains(out.String(), "sdp-trace install repo-observer") {
		t.Fatalf("usage missing new commands: %s", out.String())
	}
}

func TestCheckpointCreateAndVerifyCLI(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	runDir := filepath.Join(t.TempDir(), "run")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"run",
		"--task", "task-1",
		"--use-default-contract",
		"--output-dir", runDir,
		"--", echo, "SECRET_TOKEN_SHOULD_NOT_APPEAR",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run exit: %d err=%s", exit, errOut.String())
	}

	key, err := checkpoint.GenerateKeyPair("local-dev")
	if err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(t.TempDir(), "key.json")
	writeJSONFileForTest(t, keyPath, key)
	checkpointPath := filepath.Join(t.TempDir(), "checkpoint.json")

	out.Reset()
	errOut.Reset()
	exit = run([]string{
		"checkpoint", "create",
		"--run", runDir,
		"--out", checkpointPath,
		"--private-key", keyPath,
		"--signer-id", "local-dev",
		"--id", "checkpoint-001",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("checkpoint create exit: %d err=%s", exit, errOut.String())
	}
	if strings.Contains(out.String(), "SECRET_TOKEN_SHOULD_NOT_APPEAR") {
		t.Fatalf("checkpoint create leaked sensitive output: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{
		"checkpoint", "verify",
		"--run", runDir,
		"--checkpoint", checkpointPath,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("checkpoint verify exit: %d err=%s", exit, errOut.String())
	}
	var result checkpoint.VerificationResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("verify result: %v", err)
	}
	if result.Result != checkpoint.StatePass || result.TrustScope != checkpoint.TrustScopeLocalSigned {
		t.Fatalf("unexpected verify result: %+v", result)
	}
	if strings.Contains(out.String(), "SECRET_TOKEN_SHOULD_NOT_APPEAR") {
		t.Fatalf("checkpoint verify leaked sensitive output: %s", out.String())
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

func writeJSONFileForTest(t *testing.T, path string, value any) {
	t.Helper()
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func readJSONFileForTest(t *testing.T, path string, value any) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, value); err != nil {
		t.Fatal(err)
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
	if exit != exitCannotVerify {
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

func TestWitnessCommandRejectsUnknownKind(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"witness", "--kind", "jenkins", "--out", filepath.Join(t.TempDir(), "witness.json"), root}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("expected usage exit, got %d stdout=%s stderr=%s", exit, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "github-actions, gitlab-ci, buildkite, or customer-pki") {
		t.Fatalf("stderr missing allowed kinds: %s", errOut.String())
	}
}

func TestWitnessCommandBuildkiteRequiresExplicitEnvelope(t *testing.T) {
	for _, key := range []string{"BUILDKITE", "BUILDKITE_BUILD_ID", "BUILDKITE_JOB_ID", "BUILDKITE_COMMIT", "GITLAB_CI", "CI_PIPELINE_ID", "CI_JOB_ID", "CI_COMMIT_SHA"} {
		t.Setenv(key, "")
	}
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "agent")

	outPath := filepath.Join(t.TempDir(), "buildkite-witness.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"witness", "--kind", "buildkite", "--out", outPath, root}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("expected cannot_verify exit, got %d stderr=%s out=%s", exit, errOut.String(), out.String())
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	if !strings.Contains(string(raw), `"reason": "witness_identity_missing"`) {
		t.Fatalf("witness upgraded without envelope: %s", string(raw))
	}
}

func TestWitnessCommandBuildkitePassesWithExplicitEnvelope(t *testing.T) {
	root := t.TempDir()
	runDir := filepath.Join(root, "001-agent-session")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run: %v", err)
	}
	runJSON := []byte(`{"run_id":"pipeline-42"}`)
	if err := os.WriteFile(filepath.Join(runDir, "run.json"), runJSON, 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}
	sum := sha256.Sum256(runJSON)
	dir := t.TempDir()
	envelopePath := filepath.Join(dir, "buildkite-envelope.json")
	writeJSONFileForTest(t, envelopePath, map[string]any{
		"profile_id":            "buildkite-v1",
		"profile_version":       "1.0",
		"provider_kind":         "buildkite",
		"requested_trust_scope": "ci_witnessed",
		"source": map[string]string{
			"repository": "org/repo",
			"ref":        "refs/heads/main",
			"commit_sha": "abc123",
		},
		"ci": map[string]string{
			"provider": "buildkite",
			"run_id":   "pipeline-42",
			"job":      "verify",
		},
		"run_artifacts": []map[string]string{
			{"path": "001-agent-session/run.json", "sha256": hex.EncodeToString(sum[:])},
		},
		"profile_states": map[string]string{
			"identity_state":         "pass",
			"signer_authority_state": "pass",
			"freshness_state":        "pass",
			"artifact_binding_state": "pass",
			"source_binding_state":   "pass",
			"run_binding_state":      "pass",
			"policy_binding_state":   "pass",
			"independence_state":     "ci_isolated_job",
		},
	})
	outPath := filepath.Join(dir, "buildkite-witness.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"witness", "--kind", "buildkite", "--witness-envelope", envelopePath, "--out", outPath, root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("expected pass exit, got %d stderr=%s out=%s", exit, errOut.String(), out.String())
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	if !strings.Contains(string(raw), `"status": "pass"`) ||
		!strings.Contains(string(raw), `"established_trust_scope": "ci_witnessed"`) {
		t.Fatalf("witness did not record buildkite pass: %s", string(raw))
	}
}

func TestWitnessCommandCustomerPKIMissingFlagsUsage(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"witness", "--kind", "customer-pki", "--out", filepath.Join(t.TempDir(), "customer-pki-witness.json"), root}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("expected usage exit, got %d stdout=%s stderr=%s", exit, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "--customer-pki-authority-policy") {
		t.Fatalf("stderr missing required customer-pki flags: %s", errOut.String())
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
	if exit != exitCannotVerify {
		t.Fatalf("gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"ci_witness_gate": "cannot_verify"`) {
		t.Fatalf("gate output missing ci witness posture: %s", out.String())
	}
}

func TestGateCommandFailsForWitnessArtifactMismatch(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "agent")
	runAndWrapNamed(t, filepath.Join(root, "002-verification-run"), "verification-run", echo, "test")
	contractPath := writeGateContract(t, t.TempDir())
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "source": {"repository": "org/repo", "ref": "refs/heads/main", "commit_sha": "abc123"},
	  "ci": {"provider": "github-actions", "server_url": "https://github.com", "workflow": "sdp-trace", "job": "test", "run_id": "42", "run_attempt": "1", "actor": "octocat"},
	  "run_artifacts": [{"path": "001-agent-session/run.json", "sha256": "0000000000000000000000000000000000000000000000000000000000000000"}],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "--out", filepath.Join(t.TempDir(), "gate-result.json"), "--contract", contractPath, "--witness", witnessPath, root}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"ci_witness_gate": "fail"`) {
		t.Fatalf("gate output missing witness failure: %s", out.String())
	}
}

func TestGateCommandCannotVerifyWhenWitnessOmitsRunArtifact(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "agent")
	runAndWrapNamed(t, filepath.Join(root, "002-verification-run"), "verification-run", echo, "test")
	contractPath := writeGateContract(t, t.TempDir())
	digest := sha256FileForTest(t, filepath.Join(root, "001-agent-session", "run.json"))
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "source": {"repository": "org/repo", "ref": "refs/heads/main", "commit_sha": "abc123"},
	  "ci": {"provider": "github-actions", "server_url": "https://github.com", "workflow": "sdp-trace", "job": "test", "run_id": "42", "run_attempt": "1", "actor": "octocat"},
	  "run_artifacts": [{"path": "001-agent-session/run.json", "sha256": "`+digest+`"}],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "--out", filepath.Join(t.TempDir(), "gate-result.json"), "--contract", contractPath, "--witness", witnessPath, root}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "ci witness artifact 002-verification-run/run.json is missing from witness") {
		t.Fatalf("gate output missing omitted artifact reason: %s", out.String())
	}
}

func TestProtectedGateRequiresCheckpointPolicyAndWitnessFlags(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	gatePath := filepath.Join(t.TempDir(), "gate-result.json")
	exit := run([]string{"gate", "--profile", "protected", "--out", gatePath, t.TempDir()}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("protected gate missing inputs exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if _, err := os.Stat(gatePath); !os.IsNotExist(err) {
		t.Fatalf("protected gate wrote artifact despite usage error")
	}
}

func TestManagedAssessRequiresInputsWithoutWriting(t *testing.T) {
	gatePath := filepath.Join(t.TempDir(), "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"assess", "--profile", "managed-harness", "--out", gatePath}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("managed assess missing inputs exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if _, err := os.Stat(gatePath); !os.IsNotExist(err) {
		t.Fatalf("managed assess wrote artifact despite usage error")
	}
}

func TestManagedAssessPassesAndExplains(t *testing.T) {
	root := t.TempDir()
	paths := writeManagedFixtureInputs(t, root)
	outPath := filepath.Join(root, "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "managed-harness",
		"--out", outPath,
		"--contract", paths.contract,
		"--run", paths.run,
		"--adapter-registry", paths.registry,
		"--managed-policy", paths.policy,
		"--managed-witness", paths.witness,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("managed assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if strings.Contains(out.String(), "secret-token") || strings.Contains(out.String(), "raw prompt") {
		t.Fatalf("managed assess leaked sensitive marker: %s", out.String())
	}
	var result managed.AssessmentResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("assessment payload: %v", err)
	}
	if result.ManagedHarnessAssessment != managed.StatePass || result.SchemaVersion != managed.SchemaVersion {
		t.Fatalf("assessment result = %+v", result)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"assess", "explain", "--assessment-result", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("managed explain exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "Managed harness assessment: pass") ||
		!strings.Contains(out.String(), "Managed condition managed_witness_bound: pass") {
		t.Fatalf("managed explain missing fields: %s", out.String())
	}
}

func TestManagedAssessRejectsPostHocPolicyAndWitnessMismatch(t *testing.T) {
	root := t.TempDir()
	paths := writeManagedFixtureInputs(t, root)
	var policy managed.Policy
	readTestJSON(t, paths.policy, &policy)
	policy.PolicyProvenance.Source = "run_local"
	writeTestJSON(t, paths.policy, policy)
	var witness managed.Witness
	readTestJSON(t, paths.witness, &witness)
	witness.RunNonce = "wrong-nonce"
	writeTestJSON(t, paths.witness, witness)

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "managed-harness",
		"--out", filepath.Join(root, "assessment.json"),
		"--contract", paths.contract,
		"--run", paths.run,
		"--adapter-registry", paths.registry,
		"--managed-policy", paths.policy,
		"--managed-witness", paths.witness,
	}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("managed assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"reason_code": "post_hoc_policy"`) ||
		!strings.Contains(out.String(), `"reason_code": "managed_witness_mismatch"`) {
		t.Fatalf("managed assess missing fail reasons: %s", out.String())
	}
}

func TestManagedAssessPreviewDoesNotWriteOrVerify(t *testing.T) {
	root := t.TempDir()
	outPath := filepath.Join(root, "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"assess", "preview", "--profile", "managed-harness", "--out", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("managed preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"selected_profile": "managed_harness"`) ||
		!strings.Contains(out.String(), `"managed_policy": "absent"`) ||
		!strings.Contains(out.String(), `"claim": "preview is read-only and does not emit a managed verdict"`) {
		t.Fatalf("managed preview missing fields: %s", out.String())
	}
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Fatalf("managed preview wrote artifact")
	}
}

func TestProtectedGateMalformedNamedInputIsUsageError(t *testing.T) {
	for _, malformed := range []string{"checkpoint", "policy", "witness"} {
		t.Run(malformed, func(t *testing.T) {
			dir := t.TempDir()
			checkpointPath := filepath.Join(dir, "checkpoint.json")
			policyPath := filepath.Join(dir, "policy.json")
			witnessPath := filepath.Join(dir, "witness.json")
			gatePath := filepath.Join(dir, "gate-result.json")
			writeJSONFileForTest(t, checkpointPath, checkpoint.SignedCheckpoint{})
			writeJSONFileForTest(t, policyPath, checkpoint.TrustedCheckpointPolicy{SchemaVersion: checkpoint.PolicySchemaVersion})
			writeJSONFileForTest(t, witnessPath, demo.WitnessSummary{Kind: "github-actions", Status: demo.GatePass, TrustScope: "ci_witnessed", Reason: "ci_identity_present"})
			switch malformed {
			case "checkpoint":
				if err := os.WriteFile(checkpointPath, []byte(`{not-json`), 0o644); err != nil {
					t.Fatal(err)
				}
			case "policy":
				if err := os.WriteFile(policyPath, []byte(`{not-json`), 0o644); err != nil {
					t.Fatal(err)
				}
			case "witness":
				if err := os.WriteFile(witnessPath, []byte(`{not-json`), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run([]string{
				"gate",
				"--profile", "protected",
				"--out", gatePath,
				"--checkpoint", checkpointPath,
				"--checkpoint-policy", policyPath,
				"--witness", witnessPath,
				dir,
			}, &out, &errOut)
			if exit != exitUsage {
				t.Fatalf("protected gate malformed input exit %d err=%s out=%s", exit, errOut.String(), out.String())
			}
			if _, err := os.Stat(gatePath); !os.IsNotExist(err) {
				t.Fatalf("protected gate wrote artifact despite malformed input")
			}
		})
	}
}

func TestProtectedGatePreviewRendersAbsentInputsWithoutWriting(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "SECRET_TOKEN_BLOCK16")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "preview", "--profile", "protected", root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("protected preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"selected_profile": "protected"`) ||
		!strings.Contains(out.String(), `"checkpoint": "absent"`) ||
		!strings.Contains(out.String(), `"checkpoint_policy": "absent"`) ||
		!strings.Contains(out.String(), `"witness": "absent"`) {
		t.Fatalf("protected preview missing inspectability statuses: %s", out.String())
	}
	if strings.Contains(out.String(), "SECRET_TOKEN_BLOCK16") {
		t.Fatalf("protected preview leaked secret-like value: %s", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, "gate-result.json")); !os.IsNotExist(err) {
		t.Fatalf("protected preview wrote gate artifact")
	}
}

func TestDefaultGateDoesNotEmitProtectedFields(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "ok")
	runAndWrapNamed(t, filepath.Join(root, "002-verification-run"), "verification-run", echo, "ok")
	contractPath := writeGateContract(t, t.TempDir())
	gatePath := filepath.Join(t.TempDir(), "gate-result.json")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "--out", gatePath, "--contract", contractPath, root}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if strings.Contains(out.String(), `"protected_gate"`) || strings.Contains(out.String(), `"selected_profile"`) || strings.Contains(out.String(), `"protected_conditions"`) {
		t.Fatalf("default gate emitted protected fields: %s", out.String())
	}
}

func TestProtectedGateRejectsLocalSignedCheckpointCLI(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	runDir := filepath.Join(t.TempDir(), "run")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"run",
		"--task", "task-1",
		"--use-default-contract",
		"--output-dir", runDir,
		"--", echo, "SECRET_TOKEN_BLOCK16",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run exit %d err=%s", exit, errOut.String())
	}
	key, err := checkpoint.GenerateKeyPair("local-dev")
	if err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(t.TempDir(), "key.json")
	writeJSONFileForTest(t, keyPath, key)
	checkpointPath := filepath.Join(t.TempDir(), "checkpoint.json")
	out.Reset()
	errOut.Reset()
	exit = run([]string{"checkpoint", "create", "--run", runDir, "--out", checkpointPath, "--private-key", keyPath, "--signer-id", "local-dev"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("checkpoint create exit %d err=%s", exit, errOut.String())
	}
	policyPath := filepath.Join(t.TempDir(), "policy.json")
	writeJSONFileForTest(t, policyPath, checkpoint.TrustedCheckpointPolicy{
		SchemaVersion: checkpoint.PolicySchemaVersion,
		PolicyID:      "local-policy",
		AllowedSigners: []checkpoint.TrustedSigner{{
			SignerID:  "local-dev",
			Authority: checkpoint.AuthorityLocalDevelopment,
			PublicKey: key.PublicKey,
		}},
	})
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "2026-05-06T00:00:00Z",
	  "source": {"repository": "", "ref": "", "commit_sha": ""},
	  "run_artifacts": [],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	out.Reset()
	errOut.Reset()
	gatePath := filepath.Join(t.TempDir(), "protected-gate.json")
	exit = run([]string{
		"gate",
		"--profile", "protected",
		"--out", gatePath,
		"--checkpoint", checkpointPath,
		"--checkpoint-policy", policyPath,
		"--witness", witnessPath,
		runDir,
	}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("protected gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"protected_gate": "fail"`) ||
		!strings.Contains(out.String(), `"reason_code": "local_signed_not_protected"`) {
		t.Fatalf("protected gate missing local-signed failure: %s", out.String())
	}
	if strings.Contains(out.String(), "SECRET_TOKEN_BLOCK16") {
		t.Fatalf("protected gate leaked secret-like value: %s", out.String())
	}
}

func TestProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI(t *testing.T) {
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
		t.Fatalf("run exit %d err=%s", exit, errOut.String())
	}
	key, err := checkpoint.GenerateKeyPair("ci-signer")
	if err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(t.TempDir(), "key.json")
	writeJSONFileForTest(t, keyPath, key)
	checkpointPath := filepath.Join(t.TempDir(), "checkpoint.json")
	out.Reset()
	errOut.Reset()
	exit = run([]string{"checkpoint", "create", "--run", runDir, "--out", checkpointPath, "--private-key", keyPath, "--signer-id", "ci-signer"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("checkpoint create exit %d err=%s", exit, errOut.String())
	}
	var signed checkpoint.SignedCheckpoint
	readJSONFileForTest(t, checkpointPath, &signed)
	signed.Signer.Authority = checkpoint.AuthorityCIIsolatedJob
	writeJSONFileForTest(t, checkpointPath, signed)
	policyPath := filepath.Join(t.TempDir(), "policy.json")
	writeJSONFileForTest(t, policyPath, checkpoint.TrustedCheckpointPolicy{
		SchemaVersion: checkpoint.PolicySchemaVersion,
		PolicyID:      "ci-policy",
		AllowedSigners: []checkpoint.TrustedSigner{{
			SignerID:  "ci-signer",
			Authority: checkpoint.AuthorityCIIsolatedJob,
			PublicKey: key.PublicKey,
		}},
	})
	digest := sha256FileForTest(t, filepath.Join(runDir, "run.json"))
	runArtifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		t.Fatal(err)
	}
	freshGeneratedAt := time.Now().UTC().Format(time.RFC3339)
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "`+freshGeneratedAt+`",
	  "source": {"repository": "", "ref": "", "commit_sha": ""},
	  "ci_identity": {"run_id": "`+runArtifact.Manifest.RunID+`"},
	  "run_artifacts": [{"path": "run/run.json", "sha256": "`+digest+`"}],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	out.Reset()
	errOut.Reset()
	gatePath := filepath.Join(t.TempDir(), "protected-gate.json")
	exit = run([]string{
		"gate",
		"--profile", "protected",
		"--out", gatePath,
		"--checkpoint", checkpointPath,
		"--checkpoint-policy", policyPath,
		"--witness", witnessPath,
		runDir,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("protected gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"protected_gate": "pass"`) ||
		!strings.Contains(out.String(), `"reason_code": "protected_trust_scope_satisfied"`) ||
		!strings.Contains(out.String(), `"gate_mode": "protected"`) ||
		!strings.Contains(out.String(), `"trust_cap": "ci_signed"`) ||
		!strings.Contains(out.String(), `"ci_witness_gate": "pass"`) {
		t.Fatalf("protected gate did not pass with CI signed checkpoint: %s", out.String())
	}

	noFreshnessWitnessPath := filepath.Join(t.TempDir(), "ci-witness-no-freshness.json")
	if err := os.WriteFile(noFreshnessWitnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "source": {"repository": "", "ref": "", "commit_sha": ""},
	  "ci_identity": {"run_id": "`+runArtifact.Manifest.RunID+`"},
	  "run_artifacts": [{"path": "run/run.json", "sha256": "`+digest+`"}],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness without freshness: %v", err)
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{
		"gate",
		"--profile", "protected",
		"--out", filepath.Join(t.TempDir(), "protected-gate.json"),
		"--checkpoint", checkpointPath,
		"--checkpoint-policy", policyPath,
		"--witness", noFreshnessWitnessPath,
		runDir,
	}, &out, &errOut)
	if exit != exitCannotVerify || !strings.Contains(out.String(), `"protected_gate": "cannot_verify"`) {
		t.Fatalf("protected cannot_verify exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func TestBlock16CommittedFixturesHaveRequiredProtectedRows(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "examples", "block16-protected-gate")
	entries, err := os.ReadDir(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	seen := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".gate-result.json") {
			continue
		}
		seen++
		var result demo.GateResult
		readJSONFileForTest(t, filepath.Join(fixtureDir, entry.Name()), &result)
		if result.SchemaVersion != demo.GateSchemaVersionBlock16 || result.SelectedProfile != demo.GateProfileProtected {
			t.Fatalf("%s has schema/profile %s/%s", entry.Name(), result.SchemaVersion, result.SelectedProfile)
		}
		if result.GateMode != demo.GateProfileProtected {
			t.Fatalf("%s gate_mode = %s", entry.Name(), result.GateMode)
		}
		if len(result.ProtectedConditions) != 10 {
			t.Fatalf("%s protected condition count = %d", entry.Name(), len(result.ProtectedConditions))
		}
		for i, id := range []string{
			"protected_profile_explicitly_selected",
			"all_required_runs_present",
			"all_required_evidence_observed",
			"ci_witness_bound",
			"witness_freshness_valid",
			"checkpoint_signature_valid",
			"checkpoint_run_binding_valid",
			"checkpoint_signer_authorized",
			"protected_trust_scope_satisfied",
			"override_does_not_upgrade_profile",
		} {
			if result.ProtectedConditions[i].ID != id {
				t.Fatalf("%s condition %d = %s, want %s", entry.Name(), i, result.ProtectedConditions[i].ID, id)
			}
		}
		if got := protectedFixtureGate(result.ProtectedConditions); result.ProtectedGate != got {
			t.Fatalf("%s protected_gate = %s, want %s from condition rows", entry.Name(), result.ProtectedGate, got)
		}
		trustScope := result.ProtectedConditions[8]
		if trustScope.State == demo.GatePass {
			for _, condition := range result.ProtectedConditions[3:8] {
				if condition.State != demo.GatePass {
					t.Fatalf("%s protected trust scope passes while %s is %s", entry.Name(), condition.ID, condition.State)
				}
			}
			if result.CheckpointVerification == nil ||
				result.CheckpointVerification.Result != checkpoint.StatePass ||
				result.CheckpointVerification.TrustScope != checkpoint.TrustScopeCISigned {
				t.Fatalf("%s protected trust scope passes with checkpoint %+v", entry.Name(), result.CheckpointVerification)
			}
		}
		if result.CheckpointVerification != nil &&
			result.CheckpointVerification.TrustScope == checkpoint.TrustScopeLocalSigned &&
			result.ProtectedConditions[7].State == demo.GatePass {
			t.Fatalf("%s local-signed checkpoint overclaims protected signer authorization", entry.Name())
		}
	}
	if seen != 13 {
		t.Fatalf("fixture count = %d, want 13", seen)
	}
}

func protectedFixtureGate(conditions []demo.ProtectedCondition) string {
	gate := demo.GatePass
	for _, condition := range conditions {
		if condition.ID == "override_does_not_upgrade_profile" {
			continue
		}
		state := condition.State
		if state == demo.GateMissingTelemetry || state == "not_integrated" {
			state = demo.GateCannotVerify
		}
		if protectedFixtureSeverity(state) > protectedFixtureSeverity(gate) {
			gate = state
		}
	}
	return gate
}

func protectedFixtureSeverity(state string) int {
	switch state {
	case demo.GateFail:
		return 5
	case demo.GateCannotVerify:
		return 4
	case demo.GateNotAssessed:
		return 2
	case demo.GatePass:
		return 1
	default:
		return 0
	}
}

func TestGateExplainRendersProtectedFields(t *testing.T) {
	gatePath := filepath.Join(t.TempDir(), "protected-gate.json")
	if err := os.WriteFile(gatePath, []byte(`{
	  "schema_version": "block16-gate-result-v1",
	  "generated_at": "2026-05-06T00:00:00Z",
	  "selected_profile": "protected",
	  "local_gate": "pass",
	  "ci_witness_gate": "pass",
	  "audit_grade_gate": "cannot_verify",
	  "protected_gate": "fail",
	  "gate_mode": "protected",
	  "trust_cap": "local_signed",
	  "checkpoint_verification": {"schema_version":"block15-checkpoint-verification-v1","result":"pass","trust_scope":"local_signed","payload_digest_state":"pass","signature_state":"pass","run_binding_state":"pass","chain_binding_state":"pass","source_binding_state":"pass","nonce_binding_state":"pass","sequence_state":"pass","signer_authority_state":"pass","replay_freshness_state":"not_assessed","reasons":[]},
	  "protected_conditions": [{"id":"protected_trust_scope_satisfied","state":"fail","reason_code":"local_signed_not_protected","reason":"local signed is not protected"}],
	  "required_runs": [],
	  "required_evidence": [],
	  "observed_evidence": [],
	  "witness_bindings": [],
	  "override_requests": [],
	  "gate_conditions": [],
	  "reasons": ["local_signed_not_protected: local signed is not protected"],
	  "next_actions": ["Provide CI signed checkpoint evidence."],
	  "missing_audit_evidence": [],
	  "runs": [{"name": "run-a", "command": "deploy SECRET_TOKEN_BLOCK16"}]
	}`), 0o644); err != nil {
		t.Fatalf("write gate result: %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "explain", "--gate-result", gatePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("gate explain exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "Selected profile: protected") ||
		!strings.Contains(out.String(), "Protected gate: fail") ||
		!strings.Contains(out.String(), "Protected condition protected_trust_scope_satisfied: fail") ||
		!strings.Contains(out.String(), "Checkpoint result: pass") {
		t.Fatalf("protected explain missing protected fields: %s", out.String())
	}
	if strings.Contains(out.String(), "SECRET_TOKEN_BLOCK16") {
		t.Fatalf("protected explain leaked secret-like command: %s", out.String())
	}
}

func TestGateExplainUnsupportedArtifactCannotVerify(t *testing.T) {
	gatePath := filepath.Join(t.TempDir(), "unsupported-gate.json")
	if err := os.WriteFile(gatePath, []byte(`{"schema_version":"unknown-gate-result-v1"}`), 0o644); err != nil {
		t.Fatalf("write gate result: %v", err)
	}
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "explain", "--gate-result", gatePath}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("unsupported explain exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func TestGateExplainDoesNotPrintRawSecretLikeCommand(t *testing.T) {
	gatePath := filepath.Join(t.TempDir(), "gate-result.json")
	if err := os.WriteFile(gatePath, []byte(`{
	  "schema_version": "block14-gate-result-v1",
	  "local_gate": "fail",
	  "ci_witness_gate": "cannot_verify",
	  "audit_grade_gate": "cannot_verify",
	  "gate_mode": "observation",
	  "trust_cap": "local_observed",
	  "required_runs": [{"id": "verification_run", "wrapper_name": "verification-run", "profile": "observation", "state": "missing_telemetry", "reasons": ["missing"]}],
	  "reasons": ["missing required run"],
	  "next_actions": ["Run required wrapper verification-run through sdp-trace before evaluating advisory gate."],
	  "runs": [{"name": "run-a", "command": "deploy --token SECRET_TOKEN_BLOCK14"}]
	}`), 0o644); err != nil {
		t.Fatalf("write gate result: %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "explain", "--gate-result", gatePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("gate explain exit %d err=%s", exit, errOut.String())
	}
	if strings.Contains(out.String(), "SECRET_TOKEN_BLOCK14") {
		t.Fatalf("explain leaked secret-like command: %s", out.String())
	}
	if !strings.Contains(out.String(), "verification_run") || !strings.Contains(out.String(), "missing_telemetry") {
		t.Fatalf("explain omitted required run state: %s", out.String())
	}
}

func TestGatePreviewIsReadOnlyAndDoesNotPrintSecretLikeValues(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "SECRET_TOKEN_BLOCK14")
	contractPath := writeGateContract(t, t.TempDir())

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "preview", "--contract", contractPath, root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("gate preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if strings.Contains(out.String(), "SECRET_TOKEN_BLOCK14") {
		t.Fatalf("preview leaked secret-like value: %s", out.String())
	}
	if strings.Contains(out.String(), `"local_gate": "pass"`) {
		t.Fatalf("preview claimed pass: %s", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, "gate-result.json")); !os.IsNotExist(err) {
		t.Fatalf("preview wrote gate artifact")
	}
	if !strings.Contains(out.String(), `"command": "gate preview"`) {
		t.Fatalf("preview output missing command marker: %s", out.String())
	}
	if !strings.Contains(out.String(), `"gate_mode": "observation"`) || !strings.Contains(out.String(), `"trust_cap": "local_observed"`) {
		t.Fatalf("preview output missing mode/trust cap: %s", out.String())
	}
}

func TestGatePreviewReportsWitnessArtifactMismatch(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runAndWrapNamed(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "agent")
	contractPath := writeGateContract(t, t.TempDir())
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "source": {"repository": "org/repo", "ref": "refs/heads/main", "commit_sha": "abc123"},
	  "run_artifacts": [{"path": "001-agent-session/run.json", "sha256": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"}],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"gate", "preview", "--contract", contractPath, "--witness", witnessPath, root}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("gate preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "ci witness artifact digest mismatch") {
		t.Fatalf("preview did not report witness mismatch: %s", out.String())
	}
}

func TestGateExitCodeChecksRequiredRunStatesDirectly(t *testing.T) {
	result := demo.GateResult{
		LocalGate:      demo.GatePass,
		CIWitnessGate:  demo.GatePass,
		AuditGradeGate: demo.GatePass,
		RequiredRuns: []demo.RequiredRunResult{
			{ID: "verification_run", State: demo.GateMissingTelemetry},
		},
	}
	if got := gateExitCode(result); got != 1 {
		t.Fatalf("exit code = %d", got)
	}
}

func TestGateExitCodeUsesProtectedGateWhenSelected(t *testing.T) {
	cases := []struct {
		name          string
		protectedGate string
		want          int
	}{
		{name: "pass", protectedGate: demo.GatePass, want: 0},
		{name: "fail", protectedGate: demo.GateFail, want: 1},
		{name: "cannot_verify", protectedGate: demo.GateCannotVerify, want: exitCannotVerify},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := demo.GateResult{
				SelectedProfile: demo.GateProfileProtected,
				ProtectedGate:   tc.protectedGate,
				LocalGate:       demo.GateFail,
				CIWitnessGate:   demo.GateFail,
				AuditGradeGate:  demo.GateFail,
			}
			if got := gateExitCode(result); got != tc.want {
				t.Fatalf("exit code = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestOverrideRequestAppendsFlightRecorderEvent(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	runDir := filepath.Join(t.TempDir(), "run")
	runAndWrapNamed(t, runDir, "agent-session", echo, "agent")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"override", "request",
		"--out", runDir,
		"--id", "override-1",
		"--by", "release-captain",
		"--reason", "emergency fix",
		"--source-ref", "refs/heads/main",
		"--scope", "verification_run",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("override request exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		t.Fatalf("open run: %v", err)
	}
	last := artifact.Events[len(artifact.Events)-1]
	if last.EventType != trace.EventPolicyOverrideRequested {
		t.Fatalf("last event type = %s", last.EventType)
	}
	if got := last.EventPayload["origin"]; got != "native_cli" {
		t.Fatalf("origin = %v", got)
	}
}

func TestGateOutputIncludesOverrideWithoutPassingMissingEvidence(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runDir := filepath.Join(root, "001-agent-session")
	runAndWrapNamed(t, runDir, "agent-session", echo, "agent")
	exit := run([]string{
		"override", "request",
		"--out", runDir,
		"--id", "override-1",
		"--by", "release-captain",
		"--reason", "emergency fix",
		"--source-ref", "refs/heads/main",
		"--scope", "verification_run",
	}, &bytes.Buffer{}, &bytes.Buffer{})
	if exit != 0 {
		t.Fatalf("override request exit = %d", exit)
	}
	contractPath := writeGateContract(t, t.TempDir())

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit = run([]string{"gate", "--out", filepath.Join(t.TempDir(), "gate-result.json"), "--contract", contractPath, root}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("gate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"override_id": "override-1"`) {
		t.Fatalf("gate output missing override request: %s", out.String())
	}
	if strings.Contains(out.String(), `"local_gate": "pass"`) {
		t.Fatalf("override converted missing evidence to pass: %s", out.String())
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

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	_ = gitOutput(t, dir, args...)
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	git := mustFindCommand(t, "git")
	cmd := exec.Command(git, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func writeFile(t *testing.T, path string, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(value), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func surfaceByID(t *testing.T, status repoobserver.Status, id string) repoobserver.Surface {
	t.Helper()
	for _, surface := range status.Surfaces {
		if surface.SurfaceID == id {
			return surface
		}
	}
	t.Fatalf("surface %s not found in %+v", id, status.Surfaces)
	return repoobserver.Surface{}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore cwd %s: %v", previous, err)
		}
	})
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

func sha256FileForTest(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file for digest: %v", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
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

type managedFixturePaths struct {
	contract string
	run      string
	policy   string
	registry string
	witness  string
}

func writeManagedFixtureInputs(t *testing.T, dir string) managedFixturePaths {
	t.Helper()
	runDir := filepath.Join(dir, "run")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run: %v", err)
	}
	policy := managed.Policy{
		PolicyID: "managed-policy-v1",
		PolicyProvenance: managed.Provenance{
			Source: "vcs",
			Digest: "1111111111111111111111111111111111111111111111111111111111111111",
		},
		AuthorizedAdapters: []managed.AuthorizedAdapter{{
			AdapterID:       "opencode-adapter",
			HarnessID:       "opencode",
			AuthorityRef:    "sig:adapter",
			DeploymentRef:   "vcs:adapter",
			CapabilityIDs:   []string{"harness-events", "tool-events", "file-events", "test-events"},
			VersionRequired: "1.0.0",
		}},
		RequiredEventGroups: []managed.RequiredEventGroup{
			{ID: "harness", EventTypes: []string{"harness_lifecycle_observed"}, AcceptableProvenanceScopes: []string{"harness_observed"}},
			{ID: "tool", EventTypes: []string{"tool_call_observed"}, AcceptableProvenanceScopes: []string{"local_observed"}},
			{ID: "file", EventTypes: []string{"file_mutation_observed"}, AcceptableProvenanceScopes: []string{"local_observed"}},
			{ID: "test", EventTypes: []string{"test_observed"}, AcceptableProvenanceScopes: []string{"local_observed", "ci_witnessed"}},
		},
	}
	registry := managed.Registry{
		RegistryID: "adapter-registry-v1",
		Provenance: managed.Provenance{
			Source: "vcs",
			Digest: "2222222222222222222222222222222222222222222222222222222222222222",
		},
		Adapters: []managed.Adapter{{
			AdapterID:      "opencode-adapter",
			HarnessID:      "opencode",
			Version:        "1.0.0",
			DeploymentRef:  "vcs:adapter",
			IdentityState:  managed.IdentityVerified,
			AuthorityRef:   "sig:adapter",
			AllowedEvents:  []string{"harness_lifecycle_observed", "tool_call_observed", "file_mutation_observed", "test_observed"},
			CapabilityRefs: []string{"harness-events", "tool-events", "file-events", "test-events"},
			Capabilities: []managed.Capability{
				{ID: "harness-events", EventTypes: []string{"harness_lifecycle_observed"}, ProvenanceScope: "harness_observed"},
				{ID: "tool-events", EventTypes: []string{"tool_call_observed"}, ProvenanceScope: "local_observed"},
				{ID: "file-events", EventTypes: []string{"file_mutation_observed"}, ProvenanceScope: "local_observed"},
				{ID: "test-events", EventTypes: []string{"test_observed"}, ProvenanceScope: "local_observed"},
			},
		}},
	}
	runEvidence := managed.RunEvidence{
		RunID:           "managed-run-1",
		RunNonce:        "nonce-1",
		SourceCommit:    "abc123",
		ChainHead:       "chain-head",
		EventCount:      8,
		OutputArtifacts: []managed.ArtifactDigest{{Path: "run.json", SHA256: "run-digest"}},
		ManagedBoundaryEnrolled: &managed.ManagedBoundaryEnrolled{
			Sequence:              1,
			ManagedPolicyDigest:   policy.PolicyProvenance.Digest,
			AdapterRegistryDigest: registry.Provenance.Digest,
			AdapterID:             "opencode-adapter",
			EnrollmentSource:      "vcs",
			ManagedProfileID:      "managed-harness",
			ParentRunID:           "managed-run-1",
			RunNonce:              "nonce-1",
			EventDigest:           "enroll-digest",
		},
		ChildLaunch: managed.LaunchEvent{Sequence: 2, EventDigest: "launch-digest"},
		ObservedEvents: []managed.EvidenceEvent{
			{EventType: "harness_lifecycle_observed", ProvenanceScope: "harness_observed"},
			{EventType: "tool_call_observed", ProvenanceScope: "local_observed"},
			{EventType: "file_mutation_observed", ProvenanceScope: "local_observed"},
			{EventType: "test_observed", ProvenanceScope: "local_observed"},
		},
		TestEvidence: []managed.EvidenceEvent{{EventType: "test_observed", ProvenanceScope: "local_observed"}},
	}
	witness := managed.Witness{
		WitnessID:             "ci-witness-1",
		Status:                managed.StatePass,
		RunID:                 runEvidence.RunID,
		RunNonce:              runEvidence.RunNonce,
		SourceCommit:          runEvidence.SourceCommit,
		ManagedPolicyDigest:   policy.PolicyProvenance.Digest,
		AdapterRegistryDigest: registry.Provenance.Digest,
		AdapterIdentityDigest: "opencode-adapter:sig:adapter",
		EnrollmentEventDigest: "enroll-digest",
		LaunchEventDigest:     "launch-digest",
		ChainHead:             runEvidence.ChainHead,
		EventCount:            runEvidence.EventCount,
		FreshnessState:        managed.StatePass,
		ArtifactDigests:       []managed.ArtifactDigest{{Path: "run.json", SHA256: "run-digest"}},
	}
	contractPath := filepath.Join(dir, "contract.json")
	contract := trace.DefaultContract
	contract.RequiredEvents = []string{"harness_lifecycle_observed", "tool_call_observed", "file_mutation_observed", "test_observed"}
	writeTestJSON(t, contractPath, contract)
	paths := managedFixturePaths{
		contract: contractPath,
		run:      runDir,
		policy:   filepath.Join(dir, "managed-policy.json"),
		registry: filepath.Join(dir, "adapter-registry.json"),
		witness:  filepath.Join(dir, "managed-witness.json"),
	}
	writeTestJSON(t, filepath.Join(runDir, "run.json"), runEvidence)
	writeTestJSON(t, paths.policy, policy)
	writeTestJSON(t, paths.registry, registry)
	writeTestJSON(t, paths.witness, witness)
	return paths
}

func writeTestJSON(t *testing.T, path string, value any) {
	t.Helper()
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write json %s: %v", path, err)
	}
}

func readTestJSON(t *testing.T, path string, value any) {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read json %s: %v", path, err)
	}
	if err := json.Unmarshal(payload, value); err != nil {
		t.Fatalf("unmarshal json %s: %v", path, err)
	}
}
