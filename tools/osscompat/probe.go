package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	for {
		if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "."
		}
		cwd = parent
	}
}

// runJSONSchemaFixtures validates a checked fixture against the schema.
func runJSONSchemaFixtures() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "check-jsonschema",
		"--schemafile", filepath.Join(repoRoot(), "schema/flight-recorder-run.schema.json"),
		filepath.Join(repoRoot(), "examples/flight-recorder/local-positive/run.json"),
	).CombinedOutput()
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
	tmpDir, err := os.MkdirTemp("", "osscompat-wrap-*")
	if err != nil {
		return stateCannotVerify, fmt.Sprintf("mkdir temp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	bin := filepath.Join(tmpDir, "sdp-trace")
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", bin, "./cmd/sdp-trace")
	buildCmd.Dir = root
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		return stateCannotVerify, fmt.Sprintf("go build failed: %v\n%s", err, strings.TrimSpace(string(buildOut)))
	}

	wrapOut := filepath.Join(tmpDir, "wrap.json")
	cmd := exec.CommandContext(ctx, bin, "wrap", "/bin/true")
	cmd.Dir = tmpDir
	stdout, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = exitErr.Stderr
		}
		return stateFail, fmt.Sprintf("wrap failed: %v\nstderr: %s", err, strings.TrimSpace(string(stderr)))
	}
	if err := os.WriteFile(wrapOut, stdout, 0644); err != nil {
		return stateCannotVerify, fmt.Sprintf("write wrap output: %v", err)
	}

	if !hasTool("check-jsonschema") {
		return stateCannotVerify, "check-jsonschema not in PATH; cannot validate live wrap output"
	}
	checkCtx, checkCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer checkCancel()
	schemaPath := filepath.Join(root, "schema/flight-recorder-run.schema.json")
	// Preflight: validate a known-good fixture to confirm the tool and schema
	// are functional. If this fails, the harness or schema is broken, not the
	// wrap output.
	positiveFixture := filepath.Join(root, "examples/flight-recorder/local-positive/run.json")
	if _, err := exec.CommandContext(checkCtx, "check-jsonschema",
		"--schemafile", schemaPath,
		positiveFixture,
	).CombinedOutput(); err != nil {
		return stateCannotVerify, fmt.Sprintf("check-jsonschema preflight failed on known-positive fixture (harness/tool error, not conformance): %v", err)
	}
	schemaOut, err := exec.CommandContext(checkCtx, "check-jsonschema",
		"--schemafile", schemaPath,
		wrapOut,
	).CombinedOutput()
	if err == nil {
		return statePass, "local wrap stdout passed schema validation — this is local checkout evidence only, not source-bound drift closure"
	}
	// Distinguish schema-validation failure (exit code 1) from harness/tool
	// errors (other exit codes, signals, crashes).
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return stateFail, fmt.Sprintf("live wrap output fails schema validation (conformance failure): %s", strings.TrimSpace(string(schemaOut)))
	}
	return stateCannotVerify, fmt.Sprintf("check-jsonschema exited abnormally on wrap output (harness/tool error, not conformance): %v\n%s", err, strings.TrimSpace(string(schemaOut)))
}

// runOPAPolicy evaluates the checked-in adapter.rego against the test fixture.
func runOPAPolicy() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	regoPath := filepath.Join(repoRoot(), "examples/oss-policy/adapter.rego")
	fixturePath := filepath.Join(repoRoot(), "examples/oss-policy/test-fixture.json")
	if _, err := os.Stat(regoPath); err != nil {
		return stateCannotVerify, fmt.Sprintf("adapter.rego not found: %v", err)
	}
	if _, err := os.Stat(fixturePath); err != nil {
		return stateCannotVerify, fmt.Sprintf("test-fixture.json not found: %v", err)
	}
	cmd := exec.CommandContext(ctx, "opa", "eval",
		"--data", regoPath,
		"--input", fixturePath,
		"--format", "json",
		"data.sdp_trace.adapter.pass",
	)
	stdout, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = exitErr.Stderr
		}
		return stateFail, fmt.Sprintf("opa eval failed: %v\nstderr: %s", err, strings.TrimSpace(string(stderr)))
	}
	// Parse OPA JSON output to assert the expression evaluates to true.
	// Expected top-level structure: {"result":[{"expressions":[{"value":true}]}]}
	var opaResult struct {
		Result []struct {
			Expressions []struct {
				Value interface{} `json:"value"`
			} `json:"expressions"`
		} `json:"result"`
	}
	if err := json.Unmarshal(stdout, &opaResult); err != nil {
		return stateFail, fmt.Sprintf("opa eval output is not valid JSON: %v", err)
	}
	if len(opaResult.Result) == 0 || len(opaResult.Result[0].Expressions) == 0 {
		return stateFail, "opa eval returned no expressions"
	}
	if v, ok := opaResult.Result[0].Expressions[0].Value.(bool); !ok || !v {
		return stateFail, "opa eval did not return boolean true for expected pass fixture"
	}
	return statePass, "adapter.rego evaluates test-fixture.json as expected"
}

// runOPANegativeFixture evaluates the checked-in adapter.rego against the
// negative test fixture and asserts pass is false.
func runOPANegativeFixture() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	regoPath := filepath.Join(repoRoot(), "examples/oss-policy/adapter.rego")
	fixturePath := filepath.Join(repoRoot(), "examples/oss-policy/test-fixture-fail.json")
	if _, err := os.Stat(regoPath); err != nil {
		return stateCannotVerify, fmt.Sprintf("adapter.rego not found: %v", err)
	}
	if _, err := os.Stat(fixturePath); err != nil {
		return stateCannotVerify, fmt.Sprintf("test-fixture-fail.json not found: %v", err)
	}
	cmd := exec.CommandContext(ctx, "opa", "eval",
		"--data", regoPath,
		"--input", fixturePath,
		"--format", "json",
		"data.sdp_trace.adapter.pass",
	)
	stdout, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = exitErr.Stderr
		}
		return stateFail, fmt.Sprintf("opa eval failed: %v\nstderr: %s", err, strings.TrimSpace(string(stderr)))
	}
	var opaResult struct {
		Result []struct {
			Expressions []struct {
				Value interface{} `json:"value"`
			} `json:"expressions"`
		} `json:"result"`
	}
	if err := json.Unmarshal(stdout, &opaResult); err != nil {
		return stateFail, fmt.Sprintf("opa eval output is not valid JSON: %v", err)
	}
	if len(opaResult.Result) == 0 || len(opaResult.Result[0].Expressions) == 0 {
		return stateFail, "opa eval returned no expressions"
	}
	if v, ok := opaResult.Result[0].Expressions[0].Value.(bool); !ok || v {
		return stateFail, "opa eval did not return boolean false for expected fail fixture"
	}
	return statePass, "adapter.rego correctly rejects test-fixture-fail.json"
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
	regoPath := filepath.Join(repoRoot(), "examples/oss-policy/adapter.rego")
	fixturePath := filepath.Join(repoRoot(), "examples/oss-policy", fixtureName)
	if _, err := os.Stat(regoPath); err != nil {
		return stateCannotVerify, fmt.Sprintf("adapter.rego not found: %v", err)
	}
	if _, err := os.Stat(fixturePath); err != nil {
		return stateCannotVerify, fmt.Sprintf("%s not found: %v", fixtureName, err)
	}
	cmd := exec.CommandContext(ctx, "opa", "eval",
		"--data", regoPath,
		"--input", fixturePath,
		"--format", "json",
		"data.sdp_trace.adapter.pass",
	)
	stdout, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = exitErr.Stderr
		}
		return stateFail, fmt.Sprintf("opa eval failed: %v\nstderr: %s", err, strings.TrimSpace(string(stderr)))
	}
	var opaResult struct {
		Result []struct {
			Expressions []struct {
				Value interface{} `json:"value"`
			} `json:"expressions"`
		} `json:"result"`
	}
	if err := json.Unmarshal(stdout, &opaResult); err != nil {
		return stateFail, fmt.Sprintf("opa eval output is not valid JSON: %v", err)
	}
	if len(opaResult.Result) == 0 || len(opaResult.Result[0].Expressions) == 0 {
		return stateFail, "opa eval returned no expressions"
	}
	if v, ok := opaResult.Result[0].Expressions[0].Value.(bool); !ok || v {
		return stateFail, fmt.Sprintf("opa eval did not return boolean false for %s negative fixture", label)
	}
	return statePass, fmt.Sprintf("adapter.rego correctly rejects %s negative fixture", label)
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
	if out, err := exec.CommandContext(ctx, "cue", args...).CombinedOutput(); err != nil {
		return stateFail, fmt.Sprintf("cue import failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return statePass, "cue can import flight-recorder JSON Schema to stdout"
}

// runInTotoWrap tests in-toto-run presence.
func runInTotoWrap() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := exec.CommandContext(ctx, "in-toto-run", "--version").CombinedOutput(); err != nil {
		return stateCannotVerify, fmt.Sprintf("in-toto-run version check failed: %v", err)
	}
	return stateCannotVerify, "in-toto-run present; run manual wrap per docs/oss-replacement-compatibility.md"
}

// runCosignLocalSign tests cosign presence.
func runCosignLocalSign() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := exec.CommandContext(ctx, "cosign", "version").CombinedOutput(); err != nil {
		return stateCannotVerify, fmt.Sprintf("cosign version check failed: %v", err)
	}
	return stateCannotVerify, "cosign present; run manual sign/verify per docs/oss-replacement-compatibility.md"
}

// runSLSANegative tests slsa-verifier presence.
func runSLSANegative() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := exec.CommandContext(ctx, "slsa-verifier", "version").CombinedOutput(); err != nil {
		return stateCannotVerify, fmt.Sprintf("slsa-verifier version check failed: %v", err)
	}
	return stateCannotVerify, "slsa-verifier present; run manual negative test per docs/oss-replacement-compatibility.md"
}
