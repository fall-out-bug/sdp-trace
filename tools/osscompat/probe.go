package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		Description: "Document live wrap output vs flight-recorder-run.schema.json drift",
		Run:         runJSONSchemaWrapDrift,
	},
	{
		Name:        "opa-policy",
		NeedsTool:   "opa",
		Description: "Verify opa is present and responds to version query",
		Run:         runOPAPolicy,
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
	out, err := exec.Command("check-jsonschema",
		"--schemafile", filepath.Join(repoRoot(), "schema/flight-recorder-run.schema.json"),
		filepath.Join(repoRoot(), "examples/flight-recorder/local-positive/run.json"),
	).CombinedOutput()
	if err != nil {
		return stateFail, fmt.Sprintf("fixture validation failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return statePass, "fixture validates against schema"
}

// runJSONSchemaWrapDrift checks live wrap output against schema (expected fail).
func runJSONSchemaWrapDrift() (verifierState, string) {
	// This probe is expected to fail until the wrap output/schema drift is resolved.
	// We cannot easily capture live wrap output here without mutating state,
	// so we report cannot_verify with the documented blocker reason.
	return stateCannotVerify, "live wrap output/schema drift is documented as a blocker; see docs/oss-replacement-compatibility.md"
}

// runOPAPolicy evaluates the checked-in adapter.rego against the test fixture.
func runOPAPolicy() (verifierState, string) {
	regoPath := filepath.Join(repoRoot(), "examples/oss-policy/adapter.rego")
	fixturePath := filepath.Join(repoRoot(), "examples/oss-policy/test-fixture.json")
	if _, err := os.Stat(regoPath); err != nil {
		return stateCannotVerify, fmt.Sprintf("adapter.rego not found: %v", err)
	}
	if _, err := os.Stat(fixturePath); err != nil {
		return stateCannotVerify, fmt.Sprintf("test-fixture.json not found: %v", err)
	}
	out, err := exec.Command("opa", "eval",
		"--data", regoPath,
		"--input", fixturePath,
		"data.sdp_trace.adapter.pass",
	).CombinedOutput()
	if err != nil {
		return stateFail, fmt.Sprintf("opa eval failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	if !strings.Contains(string(out), "true") {
		return stateFail, "opa eval did not return true for expected pass fixture"
	}
	return statePass, "adapter.rego evaluates test-fixture.json as expected"
}

// runCUEImport tests CUE JSON Schema import.
func runCUEImport() (verifierState, string) {
	args := []string{
		"import",
		"--package", "sdptrace",
		"-o", "-",
		filepath.Join(repoRoot(), "schema/flight-recorder-run.schema.json"),
	}
	if out, err := exec.Command("cue", args...).CombinedOutput(); err != nil {
		return stateFail, fmt.Sprintf("cue import failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return statePass, "cue can import flight-recorder JSON Schema to stdout"
}

// runInTotoWrap tests in-toto-run presence.
func runInTotoWrap() (verifierState, string) {
	if out, err := exec.Command("in-toto-run", "--version").CombinedOutput(); err != nil {
		return stateFail, fmt.Sprintf("in-toto-run version failed: %v", err)
	} else if !strings.Contains(string(out), "in-toto") {
		return stateFail, "unexpected in-toto-run version output"
	}
	return stateCannotVerify, "in-toto-run present; run manual wrap per docs/oss-replacement-compatibility.md"
}

// runCosignLocalSign tests cosign presence.
func runCosignLocalSign() (verifierState, string) {
	if _, err := exec.Command("cosign", "version").CombinedOutput(); err != nil {
		return stateFail, fmt.Sprintf("cosign version failed: %v", err)
	}
	return stateCannotVerify, "cosign present; run manual sign/verify per docs/oss-replacement-compatibility.md"
}

// runSLSANegative tests slsa-verifier presence.
func runSLSANegative() (verifierState, string) {
	if out, err := exec.Command("slsa-verifier", "version").CombinedOutput(); err != nil {
		return stateFail, fmt.Sprintf("slsa-verifier version failed: %v", err)
	} else if !strings.Contains(string(out), "slsa") && !strings.Contains(string(out), "SLSA") {
		return stateFail, "unexpected slsa-verifier version output"
	}
	return stateCannotVerify, "slsa-verifier present; run manual negative test per docs/oss-replacement-compatibility.md"
}
