package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
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

func TestRepoRoot_FindsGit(t *testing.T) {
	// repoRoot should find the actual repository root because this test runs
	// inside the project tree.
	root := repoRoot()
	if root == "." {
		t.Skip("repoRoot fell back to '.'; cannot verify in this environment")
	}
	gitPath := filepath.Join(root, ".git")
	if _, err := os.Stat(gitPath); err != nil {
		t.Fatalf("repoRoot returned %q but .git not found: %v", root, err)
	}
}

func TestRepoRoot_Fallback(t *testing.T) {
	// Create a temporary directory with no .git and verify walking up falls back.
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		if err := os.Chdir(origWd); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if got := repoRoot(); got != "." {
		t.Fatalf("expected '.', got %q", got)
	}
}

func TestFindGitRoot_FindsGit(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if got := findGitRoot(tmpDir); got != tmpDir {
		t.Fatalf("expected %q, got %q", tmpDir, got)
	}
}

func TestFindGitRoot_Fallback(t *testing.T) {
	// tmpDir has no .git, and its parent also has no .git (in a temp tree),
	// so it should eventually reach root and fallback to "."
	tmpDir := t.TempDir()
	if got := findGitRoot(tmpDir); got != "." {
		t.Fatalf("expected '.', got %q", got)
	}
}

func TestHasGitDir(t *testing.T) {
	tmpDir := t.TempDir()
	if hasGitDir(tmpDir) {
		t.Error("expected false for dir without .git")
	}
	if err := os.Mkdir(filepath.Join(tmpDir, ".git"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if !hasGitDir(tmpDir) {
		t.Error("expected true for dir with .git")
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

func skipUnlessIntegration(t *testing.T) {
	if os.Getenv("SDPTRACE_INTEGRATION") != "1" {
		t.Skip("set SDPTRACE_INTEGRATION=1 to run tests that invoke optional external CLIs")
	}
}

func TestRunExternalTool_Success(t *testing.T) {
	out, err := runExternalTool(context.TODO(), "go", "version")
	if err != nil {
		t.Fatalf("expected go version to succeed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "go version") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRunExternalTool_Missing(t *testing.T) {
	_, err := runExternalTool(context.TODO(), "this-tool-definitely-does-not-exist-017")
	if err == nil {
		t.Fatal("expected error for missing tool")
	}
}

func TestParseWrapRunDir_Valid(t *testing.T) {
	runDir, err := parseWrapRunDir("run_dir: ./some/path\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runDir != "./some/path" {
		t.Fatalf("expected './some/path', got %q", runDir)
	}
}

func TestParseWrapRunDir_Invalid(t *testing.T) {
	cases := []string{
		"",
		"unexpected format",
		"run_dir:",
		"run_dir: a b c",
	}
	for _, c := range cases {
		_, err := parseWrapRunDir(c)
		if err == nil {
			t.Fatalf("expected error for %q", c)
		}
	}
}

func TestValidateRunDirUnderTmp_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	runDir := "run-123"
	if err := os.MkdirAll(filepath.Join(tmpDir, runDir), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, runDir, "run.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("write run.json: %v", err)
	}
	if err := validateRunDirUnderTmp(runDir, tmpDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRunDirUnderTmp_Absolute(t *testing.T) {
	tmpDir := t.TempDir()
	if err := validateRunDirUnderTmp("/etc/passwd", tmpDir); err == nil {
		t.Fatal("expected error for absolute path")
	}
}

func TestValidateRunDirUnderTmp_Traversal(t *testing.T) {
	tmpDir := t.TempDir()
	if err := validateRunDirUnderTmp("../../etc/passwd", tmpDir); err == nil {
		t.Fatal("expected error for traversal path")
	}
}

func TestCheckRunDirSafe_Absolute(t *testing.T) {
	if err := checkRunDirSafe("/absolute/path"); err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckRunDirSafe_Traversal(t *testing.T) {
	if err := checkRunDirSafe("../escape"); err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckRunDirSafe_OK(t *testing.T) {
	if err := checkRunDirSafe("subdir"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckRunJSONUnderTmp_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	runDir := "run-123"
	if err := os.MkdirAll(filepath.Join(tmpDir, runDir), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, runDir, "run.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("write run.json: %v", err)
	}
	if err := checkRunJSONUnderTmp(filepath.Join(tmpDir, runDir, "run.json"), tmpDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckRunJSONUnderTmp_Missing(t *testing.T) {
	tmpDir := t.TempDir()
	if err := checkRunJSONUnderTmp(filepath.Join(tmpDir, "run.json"), tmpDir); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunCheckJSONSchema_NotAssessedWithoutTool(t *testing.T) {
	if hasTool("check-jsonschema") {
		t.Skip("check-jsonschema is present; skipping negative test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := runCheckJSONSchema(ctx, "schema.json", "data.json")
	if err == nil {
		t.Fatal("expected error when check-jsonschema is missing")
	}
}

func TestParseOPAEvalResult_ParsesTrue(t *testing.T) {
	stdout := `{"result":[{"expressions":[{"value":true}]}]}`
	pass, err := parseOPAEvalResult([]byte(stdout))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pass {
		t.Fatal("expected true")
	}
}

func TestParseOPAEvalResult_ParsesFalse(t *testing.T) {
	stdout := `{"result":[{"expressions":[{"value":false}]}]}`
	pass, err := parseOPAEvalResult([]byte(stdout))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass {
		t.Fatal("expected false")
	}
}

func TestParseOPAEvalResult_InvalidJSON(t *testing.T) {
	_, err := parseOPAEvalResult([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseOPAEvalResult_EmptyResult(t *testing.T) {
	_, err := parseOPAEvalResult([]byte(`{"result":[]}`))
	if err == nil {
		t.Fatal("expected error for empty result")
	}
}

func TestParseOPAEvalResult_NonBool(t *testing.T) {
	_, err := parseOPAEvalResult([]byte(`{"result":[{"expressions":[{"value":"string"}]}]}`))
	if err == nil {
		t.Fatal("expected error for non-boolean value")
	}
}

func TestAssertOPAResult_Pass(t *testing.T) {
	s, r := assertOPAResult(true, true, "fixture")
	if s != statePass {
		t.Fatalf("expected pass, got %s", s)
	}
	if !strings.Contains(r, "evaluates") {
		t.Fatalf("unexpected reason: %s", r)
	}
}

func TestAssertOPAResult_Fail(t *testing.T) {
	s, r := assertOPAResult(false, true, "fixture")
	if s != stateFail {
		t.Fatalf("expected fail, got %s", s)
	}
	if !strings.Contains(r, "boolean true") {
		t.Fatalf("unexpected reason: %s", r)
	}
}

func TestAssertOPAResult_NegativePass(t *testing.T) {
	s, r := assertOPAResult(false, false, "fixture")
	if s != statePass {
		t.Fatalf("expected pass, got %s", s)
	}
	if !strings.Contains(r, "rejects") {
		t.Fatalf("unexpected reason: %s", r)
	}
}

func TestAssertOPAResult_NegativeFail(t *testing.T) {
	s, r := assertOPAResult(true, false, "fixture")
	if s != stateFail {
		t.Fatalf("expected fail, got %s", s)
	}
	if !strings.Contains(r, "boolean false") {
		t.Fatalf("unexpected reason: %s", r)
	}
}

func TestOPAFixturePaths_MissingRego(t *testing.T) {
	_, _, err := opaFixturePaths("nonexistent.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWrapArgs(t *testing.T) {
	args := wrapArgs()
	if len(args) == 0 {
		t.Fatal("expected non-empty args")
	}
	if args[0] != "wrap" {
		t.Fatalf("expected 'wrap', got %q", args[0])
	}
}

func TestRunWrapCommand_MissingBinary(t *testing.T) {
	ctx := context.Background()
	_, err := runWrapCommand(ctx, "/nonexistent/binary-017", []string{"wrap"}, ".")
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
}

func TestRunSchemaValidation_MissingTool(t *testing.T) {
	if hasTool("check-jsonschema") {
		t.Skip("check-jsonschema is present")
	}
	ctx := context.Background()
	_, reason := runSchemaValidation(ctx, ".", ".", "run-1")
	if !strings.Contains(reason, "preflight failed") {
		t.Fatalf("expected preflight failure, got: %s", reason)
	}
}

func TestOPAPreflight_MissingTool(t *testing.T) {
	if hasTool("opa") {
		t.Skip("opa is present")
	}
	ctx := context.Background()
	s, r := opaPreflight(ctx)
	if s != stateCannotVerify {
		t.Fatalf("expected cannot_verify, got %s", s)
	}
	if r == "" {
		t.Fatal("expected reason")
	}
}

// Direct probe function tests (skip when external tools are present).

func TestRunJSONSchemaFixtures_Direct(t *testing.T) {
	if hasTool("check-jsonschema") {
		t.Skip("check-jsonschema is present; skipping environment-sensitive test")
	}
	state, reason := runJSONSchemaFixtures()
	if state != stateFail {
		t.Fatalf("expected fail when tool missing, got %s: %s", state, reason)
	}
	if reason == "" {
		t.Fatal("expected reason")
	}
}

func TestRunCUEImport_Direct(t *testing.T) {
	if hasTool("cue") {
		t.Skip("cue is present; skipping environment-sensitive test")
	}
	state, reason := runCUEImport()
	if state != stateFail {
		t.Fatalf("expected fail when tool missing, got %s: %s", state, reason)
	}
	if reason == "" {
		t.Fatal("expected reason")
	}
}

func TestRunInTotoWrap_Direct(t *testing.T) {
	if hasTool("in-toto-run") {
		t.Skip("in-toto-run is present; skipping environment-sensitive test")
	}
	state, reason := runInTotoWrap()
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify, got %s: %s", state, reason)
	}
	if !strings.Contains(reason, "in-toto-run") || !strings.Contains(reason, "failed") {
		t.Fatalf("expected actionable reason when tool missing, got: %s", reason)
	}
}


func TestRunCosignLocalSign_Direct(t *testing.T) {
	if hasTool("cosign") {
		t.Skip("cosign is present; skipping environment-sensitive test")
	}
	state, reason := runCosignLocalSign()
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify, got %s: %s", state, reason)
	}
	if !strings.Contains(reason, "cosign") || !strings.Contains(reason, "failed") {
		t.Fatalf("expected actionable reason when tool missing, got: %s", reason)
	}
}


func TestRunSLSANegative_Direct(t *testing.T) {
	if hasTool("slsa-verifier") {
		t.Skip("slsa-verifier is present; skipping environment-sensitive test")
	}
	state, reason := runSLSANegative()
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify, got %s: %s", state, reason)
	}
	if !strings.Contains(reason, "slsa-verifier") || !strings.Contains(reason, "failed") {
		t.Fatalf("expected actionable reason when tool missing, got: %s", reason)
	}
}

func TestRunOPAPolicy_Direct(t *testing.T) {
	if hasTool("opa") {
		t.Skip("opa is present; skipping environment-sensitive test")
	}
	state, reason := runOPAPolicy()
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify when opa missing, got %s: %s", state, reason)
	}
	if reason == "" {
		t.Fatal("expected reason")
	}
}

func TestRunOPANegativeFixture_Direct(t *testing.T) {
	if hasTool("opa") {
		t.Skip("opa is present; skipping environment-sensitive test")
	}
	state, reason := runOPANegativeFixture()
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify when opa missing, got %s: %s", state, reason)
	}
	if reason == "" {
		t.Fatal("expected reason")
	}
}

func TestRunOPANegativeTraceID_Direct(t *testing.T) {
	if hasTool("opa") {
		t.Skip("opa is present; skipping environment-sensitive test")
	}
	state, reason := runOPANegativeTraceID()
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify when opa missing, got %s: %s", state, reason)
	}
	if reason == "" {
		t.Fatal("expected reason")
	}
}

func TestRunOPANegativeProvenance_Direct(t *testing.T) {
	if hasTool("opa") {
		t.Skip("opa is present; skipping environment-sensitive test")
	}
	state, reason := runOPANegativeProvenance()
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify when opa missing, got %s: %s", state, reason)
	}
	if reason == "" {
		t.Fatal("expected reason")
	}
}

func TestRunJSONSchemaWrapDrift_Direct(t *testing.T) {

	state, reason := runJSONSchemaWrapDrift()
	// This probe builds sdp-trace; if go is present it should at least get
	// past the build step.
	if state != statePass && state != stateFail && state != stateCannotVerify {
		t.Fatalf("unexpected state %s: %s", state, reason)
	}
}

// Integration tests (require SDPTRACE_INTEGRATION=1 and tools on PATH).

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
func TestInterpretSchemaCheckResult_Pass(t *testing.T) {
	s, r := interpretSchemaCheckResult(nil, nil)
	if s != statePass {
		t.Fatalf("expected pass, got %s", s)
	}
	if r == "" {
		t.Fatal("expected reason")
	}
}

func TestInterpretSchemaCheckResult_Fail(t *testing.T) {
	// Use a real command that exits 1 to get a genuine *exec.ExitError.
	cmd := exec.Command("go", "tool", "nonexistent-subcommand-017")
	_, _ = cmd.CombinedOutput()
	if cmd.ProcessState == nil || cmd.ProcessState.ExitCode() == 0 {
		cmd = exec.Command("false")
		_ = cmd.Run()
	}
	exitErr := &exec.ExitError{ProcessState: cmd.ProcessState}
	if exitErr.ExitCode() != 1 {
		t.Skip("cannot produce exit code 1")
	}
	s, r := interpretSchemaCheckResult([]byte("validation error details"), exitErr)
	if s != stateFail {
		t.Fatalf("expected fail, got %s", s)
	}
	if r == "" {
		t.Fatal("expected reason")
	}
}

func TestInterpretSchemaCheckResult_CannotVerify(t *testing.T) {
	s, r := interpretSchemaCheckResult(nil, fmt.Errorf("some error"))
	if s != stateCannotVerify {
		t.Fatalf("expected cannot_verify, got %s", s)
	}
	if r == "" {
		t.Fatal("expected reason")
	}
}

func TestRunOPAPolicyEval_MissingFixture(t *testing.T) {
	ctx := context.Background()
	state, reason := runOPAPolicyEval(ctx, "nonexistent-fixture-017.json", true, "label")
	if state != stateCannotVerify {
		t.Fatalf("expected cannot_verify, got %s: %s", state, reason)
	}
}

func TestRunOPAPolicyEval_MissingOPA(t *testing.T) {
	if hasTool("opa") {
		t.Skip("opa is present")
	}
	ctx := context.Background()
	// Use an existing fixture path; runOPAEval will fail because opa is missing.
	state, reason := runOPAPolicyEval(ctx, "test-fixture.json", true, "label")
	if state != stateFail {
		t.Fatalf("expected fail, got %s: %s", state, reason)
	}
}

func TestRunOPAEval_MissingTool(t *testing.T) {
	if hasTool("opa") {
		t.Skip("opa is present")
	}
	ctx := context.Background()
	_, err := runOPAEval(ctx, "/dev/null", "/dev/null", "data.test")
	if err == nil {
		t.Fatal("expected error when opa is missing")
	}
}

func TestRunWrapAndParse_MissingBinary(t *testing.T) {
	ctx := context.Background()
	_, s, r := runWrapAndParse(ctx, "/nonexistent-binary-017", t.TempDir())
	if s != stateFail {
		t.Fatalf("expected fail, got %s: %s", s, r)
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
	if sv, ok := obj["schema_version"].(string); ok && sv == "1.0.0" {
		t.Fatal("expected schema_version mismatch in frozen run.json")
	}
}

func TestBuildSDPTraceInTemp_MissingRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, _, err := buildSDPTraceInTemp(ctx, "/nonexistent-root-017")
	if err == nil {
		t.Fatal("expected build failure for nonexistent root")
	}
}

func TestSummarize(t *testing.T) {
	results := []probeResult{
		{State: statePass},
		{State: stateFail},
		{State: stateCannotVerify},
		{State: stateNotAssessed},
		{State: statePass},
	}
	p, f, c, n := summarize(results)
	if p != 2 || f != 1 || c != 1 || n != 1 {
		t.Fatalf("expected 2,1,1,1 got %d,%d,%d,%d", p, f, c, n)
	}
}

func TestBoolToInt(t *testing.T) {
	if boolToInt(true) != 1 {
		t.Error("expected 1")
	}
	if boolToInt(false) != 0 {
		t.Error("expected 0")
	}
}

func TestFormatResultLine(t *testing.T) {
	line := formatResultLine(probeResult{Name: "abc", State: statePass, Reason: "ok"}, 10)
	if !strings.Contains(line, "abc") {
		t.Error("missing name")
	}
	if !strings.Contains(line, "pass") {
		t.Error("missing state")
	}
	if !strings.Contains(line, "ok") {
		t.Error("missing reason")
	}
}

func TestMaxNameWidth(t *testing.T) {
	results := []probeResult{{Name: "short"}}
	if maxNameWidth(results) != 24 {
		t.Fatalf("expected 24, got %d", maxNameWidth(results))
	}
	results = []probeResult{{Name: "this-is-a-very-long-name-indeed"}}
	if maxNameWidth(results) != len("this-is-a-very-long-name-indeed") {
		t.Fatalf("expected %d, got %d", len("this-is-a-very-long-name-indeed"), maxNameWidth(results))
	}
}
