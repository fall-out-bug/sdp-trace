package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// verifierState is the canonical probe result state.
type verifierState string

const (
	statePass         verifierState = "pass"
	stateFail         verifierState = "fail"
	stateCannotVerify verifierState = "cannot_verify"
	stateNotAssessed  verifierState = "not_assessed"
)

// probeResult is the output of a single compatibility probe.
type probeResult struct {
	Name   string        `json:"name"`
	State  verifierState `json:"state"`
	Reason string        `json:"reason,omitempty"`
}

// probe is a single compatibility probe that can be run.
type probe struct {
	Name        string
	NeedsTool   string
	Run         func() (verifierState, string)
	Description string
}

// registry holds all defined probes.
var registry = []probe{
	{
		Name:        "jsonschema-fixtures",
		NeedsTool:   "check-jsonschema",
		Description: "Verify check-jsonschema is present and can validate fixtures (run manually per docs)",
		Run:         runJSONSchemaFixtures,
	},
	{
		Name:        "jsonschema-wrap-drift",
		NeedsTool:   "check-jsonschema",
		Description: "Document live wrap output vs flight-recorder-run.schema.json drift",
		Run:         runJSONSchemaWrapDrift,
	},
	{
		Name:        "opa-policy",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the positive test fixture",
		Run:         runOPAPolicy,
	},
	{
		Name:        "opa-negative",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the combined negative test fixture",
		Run:         runOPANegativeFixture,
	},
	{
		Name:        "opa-negative-traceid",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the negative trace_id fixture",
		Run:         runOPANegativeTraceID,
	},
	{
		Name:        "opa-negative-provenance",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the negative provenance fixture",
		Run:         runOPANegativeProvenance,
	},
	{
		Name:        "cue-import",
		NeedsTool:   "cue",
		Description: "Verify cue can import JSON Schema to stdout without mutating working tree",
		Run:         runCUEImport,
	},
	{
		Name:        "intoto-wrap",
		NeedsTool:   "in-toto-run",
		Description: "Verify in-toto-run is present and responds to version query",
		Run:         runInTotoWrap,
	},
	{
		Name:        "cosign-local-sign",
		NeedsTool:   "cosign",
		Description: "Verify cosign is present and responds to version query",
		Run:         runCosignLocalSign,
	},
	{
		Name:        "slsa-negative",
		NeedsTool:   "slsa-verifier",
		Description: "Verify slsa-verifier is present and responds to version query",
		Run:         runSLSANegative,
	},
}

// hasTool reports whether tool is in $PATH.
func hasTool(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// repoRoot returns the repository root by walking up from the current
// working directory until it finds a .git directory or reaches the filesystem
// root. It falls back to "." if the root cannot be determined.
func repoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return findGitRoot(cwd)
}

func findGitRoot(cwd string) string {
	for {
		if hasGitDir(cwd) {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "."
		}
		cwd = parent
	}
}

func hasGitDir(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// runExternalTool executes an external command and returns combined output.
func runExternalTool(ctx context.Context, name string, args ...string) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// buildSDPTraceInTemp builds the sdp-trace binary in a temporary directory.
func buildSDPTraceInTemp(ctx context.Context, root string) (bin string, tmpDir string, err error) {
	tmpDir, err = os.MkdirTemp("", "osscompat-wrap-*")
	if err != nil {
		return "", "", fmt.Errorf("mkdir temp: %w", err)
	}
	bin = filepath.Join(tmpDir, "sdp-trace")
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", bin, "./cmd/sdp-trace")
	buildCmd.Dir = root
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		return "", tmpDir, fmt.Errorf("go build failed: %w\n%s", err, strings.TrimSpace(string(buildOut)))
	}
	return bin, tmpDir, nil
}

// parseWrapRunDir extracts the run directory from wrap stdout.
func parseWrapRunDir(stdout string) (string, error) {
	fields := strings.Fields(strings.TrimSpace(stdout))
	if len(fields) != 2 || fields[0] != "run_dir:" {
		return "", fmt.Errorf("unexpected wrap stdout format: %q", stdout)
	}
	return fields[1], nil
}

// validateRunDirUnderTmp ensures runDir is safely contained within tmpDir.
func validateRunDirUnderTmp(runDir, tmpDir string) error {
	if err := checkRunDirSafe(runDir); err != nil {
		return err
	}
	runJSONPath := filepath.Join(tmpDir, filepath.Clean(runDir), "run.json")
	if err := checkRunJSONUnderTmp(runJSONPath, tmpDir); err != nil {
		return err
	}
	if _, err := os.Stat(runJSONPath); err != nil {
		return fmt.Errorf("run.json not found at expected path %s: %w", runJSONPath, err)
	}
	return nil
}

func checkRunDirSafe(runDir string) error {
	if filepath.IsAbs(runDir) {
		return fmt.Errorf("run_dir is absolute (possible traversal): %q", runDir)
	}
	if strings.HasPrefix(filepath.Clean(runDir), "..") {
		return fmt.Errorf("run_dir escapes tmpDir (possible traversal): %q", runDir)
	}
	return nil
}

func checkRunJSONUnderTmp(runJSONPath, tmpDir string) error {
	resolvedPath, err := filepath.EvalSymlinks(runJSONPath)
	if err != nil {
		return fmt.Errorf("run.json path resolution failed: %w", err)
	}
	resolvedTmp, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		return fmt.Errorf("tmpDir path resolution failed: %w", err)
	}
	if !strings.HasPrefix(resolvedPath, resolvedTmp+string(filepath.Separator)) {
		return fmt.Errorf("run.json resolved outside tmpDir: %s", resolvedPath)
	}
	return nil
}

// runCheckJSONSchema runs check-jsonschema against a JSON file.
func runCheckJSONSchema(ctx context.Context, schemaPath, jsonPath string) ([]byte, error) {
	out, err := runExternalTool(ctx, "check-jsonschema", "--schemafile", schemaPath, jsonPath)
	return out, err
}

// runOPAEval executes OPA evaluation and returns the boolean result of the query.
func runOPAEval(ctx context.Context, regoPath, fixturePath, query string) (bool, error) {
	cmd := exec.CommandContext(ctx, "opa", "eval",
		"--data", regoPath,
		"--input", fixturePath,
		"--format", "json",
		query,
	)
	stdout, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = exitErr.Stderr
		}
		return false, fmt.Errorf("opa eval failed: %w\nstderr: %s", err, strings.TrimSpace(string(stderr)))
	}
	return parseOPAEvalResult(stdout)
}

type opaEvalOutput struct {
	Result []opaExpressionSet `json:"result"`
}

type opaExpressionSet struct {
	Expressions []struct {
		Value interface{} `json:"value"`
	} `json:"expressions"`
}

func parseOPAEvalResult(stdout []byte) (bool, error) {
	var out opaEvalOutput
	if err := json.Unmarshal(stdout, &out); err != nil {
		return false, fmt.Errorf("opa eval output is not valid JSON: %w", err)
	}
	if err := checkOPAExpressions(out.Result); err != nil {
		return false, err
	}
	v, ok := out.Result[0].Expressions[0].Value.(bool)
	if !ok {
		return false, fmt.Errorf("opa eval result is not a boolean")
	}
	return v, nil
}

func checkOPAExpressions(result []opaExpressionSet) error {
	if len(result) == 0 || len(result[0].Expressions) == 0 {
		return fmt.Errorf("opa eval returned no expressions")
	}
	return nil
}

// runJSONSchemaFixtures validates a checked fixture against the schema.
func runJSONSchemaFixtures() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	root := repoRoot()
	out, err := runCheckJSONSchema(ctx,
		filepath.Join(root, "schema/flight-recorder-run.schema.json"),
		filepath.Join(root, "examples/flight-recorder/local-positive/run.json"),
	)
	if err != nil {
		return stateFail, fmt.Sprintf("fixture validation failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return statePass, "fixture validates against schema"
}

// runJSONSchemaWrapDrift builds sdp-trace from source, runs wrap in an isolated
// temp directory, and checks the live stdout against the schema. The probe
// reports fail when schema validation fails (conformance failure) and pass
// only if the drift is unexpectedly fixed.
func runJSONSchemaWrapDrift() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	root := repoRoot()
	bin, tmpDir, err := buildSDPTraceInTemp(ctx, root)
	if tmpDir != "" {
		defer os.RemoveAll(tmpDir)
	}
	if err != nil {
		return stateCannotVerify, err.Error()
	}
	return checkWrapDrift(ctx, root, bin, tmpDir)
}

func checkWrapDrift(ctx context.Context, root, bin, tmpDir string) (verifierState, string) {
	runDir, s, r := runWrapAndParse(ctx, bin, tmpDir)
	if s != "" {
		return s, r
	}
	if !hasTool("check-jsonschema") {
		return stateCannotVerify, "check-jsonschema not in PATH; cannot validate live wrap output"
	}
	return runSchemaValidation(ctx, root, tmpDir, runDir)
}

func runWrapAndParse(ctx context.Context, bin, tmpDir string) (string, verifierState, string) {
	stdout, err := runWrapCommand(ctx, bin, wrapArgs(), tmpDir)
	if err != nil {
		return "", stateFail, err.Error()
	}
	runDir, err := parseWrapRunDir(string(stdout))
	if err != nil {
		return "", stateCannotVerify, err.Error()
	}
	if err := validateRunDirUnderTmp(runDir, tmpDir); err != nil {
		return "", stateCannotVerify, err.Error()
	}
	return runDir, "", ""
}

func wrapArgs() []string {
	if runtime.GOOS == "windows" {
		return []string{"wrap", "cmd", "/c", "exit", "0"}
	}
	return []string{"wrap", "true"}
}

func runWrapCommand(ctx context.Context, bin string, args []string, dir string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	stdout, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = exitErr.Stderr
		}
		return nil, fmt.Errorf("wrap failed: %v\nstderr: %s", err, strings.TrimSpace(string(stderr)))
	}
	return stdout, nil
}

func runSchemaValidation(ctx context.Context, root, tmpDir, runDir string) (verifierState, string) {
	checkCtx, checkCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer checkCancel()
	schemaPath := filepath.Join(root, "schema/flight-recorder-run.schema.json")
	positiveFixture := filepath.Join(root, "examples/flight-recorder/local-positive/run.json")
	if _, err := runCheckJSONSchema(checkCtx, schemaPath, positiveFixture); err != nil {
		return stateCannotVerify, fmt.Sprintf("check-jsonschema preflight failed on known-positive fixture (harness/tool error, not conformance): %v", err)
	}
	runJSONPath := filepath.Join(tmpDir, filepath.Clean(runDir), "run.json")
	out, err := runCheckJSONSchema(checkCtx, schemaPath, runJSONPath)
	return interpretSchemaCheckResult(out, err)
}

func interpretSchemaCheckResult(out []byte, err error) (verifierState, string) {
	if err == nil {
		return statePass, "local wrap run.json passed schema validation — this is local checkout evidence only, not source-bound drift closure"
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return stateFail, fmt.Sprintf("live wrap run.json fails schema validation (conformance failure): %s", strings.TrimSpace(string(out)))
	}
	return stateCannotVerify, fmt.Sprintf("check-jsonschema exited abnormally on wrap run.json (harness/tool error, not conformance): %v\n%s", err, strings.TrimSpace(string(out)))
}

// opaPreflight checks whether the installed OPA supports `import rego.v1`.
// If not, it returns `cannot_verify` so that version-mismatch failures are
// not misreported as conformance `fail`.
func opaPreflight(ctx context.Context) (verifierState, string) {
	tmpDir, err := os.MkdirTemp("", "osscompat-opa-preflight-*")
	if err != nil {
		return stateCannotVerify, fmt.Sprintf("mkdir temp for opa preflight: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	preflightPath := filepath.Join(tmpDir, "preflight.rego")
	if err := os.WriteFile(preflightPath, []byte(opaPreflightRego), 0644); err != nil {
		return stateCannotVerify, fmt.Sprintf("write opa preflight: %v", err)
	}
	out, err := runExternalTool(ctx, "opa", "eval", "--data", preflightPath, "--format", "raw", "data.sdp_trace.preflight.allow")
	if err != nil {
		return stateCannotVerify, fmt.Sprintf("opa does not support rego.v1 syntax (version too old?): %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return "", ""
}

const opaPreflightRego = `package sdp_trace.preflight
import rego.v1
allow := true
`

func opaFixturePaths(name string) (regoPath, fixturePath string, err error) {
	root := repoRoot()
	regoPath = filepath.Join(root, "examples/oss-policy/adapter.rego")
	fixturePath = filepath.Join(root, "examples/oss-policy", name)
	if _, err := os.Stat(regoPath); err != nil {
		return "", "", fmt.Errorf("adapter.rego not found: %w", err)
	}
	if _, err := os.Stat(fixturePath); err != nil {
		return "", "", fmt.Errorf("%s not found: %w", name, err)
	}
	return regoPath, fixturePath, nil
}

func assertOPAResult(pass bool, expectPass bool, label string) (verifierState, string) {
	if pass == expectPass {
		if expectPass {
			return statePass, fmt.Sprintf("adapter.rego evaluates %s as expected", label)
		}
		return statePass, fmt.Sprintf("adapter.rego correctly rejects %s", label)
	}
	if expectPass {
		return stateFail, "opa eval did not return boolean true for expected pass fixture"
	}
	return stateFail, fmt.Sprintf("opa eval did not return boolean false for %s negative fixture", label)
}

// runOPAPolicy evaluates the checked-in adapter.rego against the test fixture.
func runOPAPolicy() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if s, r := opaPreflight(ctx); s != "" {
		return s, r
	}
	return runOPAPolicyEval(ctx, "test-fixture.json", true, "test-fixture.json")
}

// runOPANegativeFixture evaluates the checked-in adapter.rego against the
// negative test fixture and asserts pass is false.
func runOPANegativeFixture() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if s, r := opaPreflight(ctx); s != "" {
		return s, r
	}
	return runOPAPolicyEval(ctx, "test-fixture-fail.json", false, "test-fixture-fail.json")
}

// runOPANegativeTraceID evaluates adapter.rego against the trace_id-only
// negative fixture and asserts pass is false.
func runOPANegativeTraceID() (verifierState, string) {
	return runOPANegativeFixturePath("test-fixture-fail-traceid.json", "trace_id-only")
}

// runOPANegativeProvenance evaluates adapter.rego against the provenance-only
// negative fixture and asserts pass is false.
func runOPANegativeProvenance() (verifierState, string) {
	return runOPANegativeFixturePath("test-fixture-fail-provenance.json", "provenance-only")
}

// runOPANegativeFixturePath is a shared helper for per-rule negative fixtures.
func runOPANegativeFixturePath(fixtureName, label string) (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if s, r := opaPreflight(ctx); s != "" {
		return s, r
	}
	return runOPAPolicyEval(ctx, fixtureName, false, label)
}

func runOPAPolicyEval(ctx context.Context, fixtureName string, expectPass bool, label string) (verifierState, string) {
	regoPath, fixturePath, err := opaFixturePaths(fixtureName)
	if err != nil {
		return stateCannotVerify, err.Error()
	}
	pass, err := runOPAEval(ctx, regoPath, fixturePath, "data.sdp_trace.adapter.pass")
	if err != nil {
		return stateFail, err.Error()
	}
	return assertOPAResult(pass, expectPass, label)
}

// runCUEImport tests CUE JSON Schema import.
func runCUEImport() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	args := []string{
		"import",
		"--package", "sdptrace",
		"-o", "-",
		filepath.Join(repoRoot(), "schema/flight-recorder-run.schema.json"),
	}
	if out, err := runExternalTool(ctx, "cue", args...); err != nil {
		return stateFail, fmt.Sprintf("cue import failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return statePass, "cue can import flight-recorder JSON Schema to stdout"
}

// runInTotoWrap tests in-toto-run presence.
func runInTotoWrap() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := runExternalTool(ctx, "in-toto-run", "--version"); err != nil {
		return stateCannotVerify, fmt.Sprintf("in-toto-run version check failed: %v", err)
	}
	return stateCannotVerify, "in-toto-run present; run manual wrap per docs/oss-replacement-compatibility.md"
}

// runCosignLocalSign tests cosign presence.
func runCosignLocalSign() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := runExternalTool(ctx, "cosign", "version"); err != nil {
		return stateCannotVerify, fmt.Sprintf("cosign version check failed: %v", err)
	}
	return stateCannotVerify, "cosign present; run manual sign/verify per docs/oss-replacement-compatibility.md"
}

// runSLSANegative tests slsa-verifier presence.
func runSLSANegative() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := runExternalTool(ctx, "slsa-verifier", "version"); err != nil {
		return stateCannotVerify, fmt.Sprintf("slsa-verifier version check failed: %v", err)
	}
	return stateCannotVerify, "slsa-verifier present; run manual negative test per docs/oss-replacement-compatibility.md"
}
