package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSpec019MonitoringGateProofPack(t *testing.T) {
	dir := t.TempDir()
	runDir := filepath.Join(dir, "run")

	wrapArgs := append([]string{"wrap", "--name", "spec019-proof", "--output-dir", runDir, "--"}, proofCommandForSpec019()...)
	exit, stdout, stderr := runCLIForSpec019(wrapArgs...)
	if exit != 0 {
		t.Fatalf("wrap exit=%d stdout=%s stderr=%s", exit, stdout, stderr)
	}

	exit, stdout, stderr = runCLIForSpec019("verify", runDir)
	if exit != 0 || !strings.Contains(stdout, `"result": "observed"`) {
		t.Fatalf("verify exit=%d stdout=%s stderr=%s", exit, stdout, stderr)
	}

	exit, stdout, stderr = runCLIForSpec019("query", "--query", "missing-evidence", runDir)
	if exit != 0 || !strings.Contains(stdout, `"rows": []`) {
		t.Fatalf("query exit=%d stdout=%s stderr=%s", exit, stdout, stderr)
	}

	reportDir := filepath.Join(dir, "report")
	exit, stdout, stderr = runCLIForSpec019("report", "--out", reportDir, runDir)
	if exit != 0 || !strings.Contains(stdout, `"run_count": 1`) {
		t.Fatalf("report exit=%d stdout=%s stderr=%s", exit, stdout, stderr)
	}
	for _, name := range []string{"summary.json", "evidence-table.json", "missing-telemetry.json", "timeline.md"} {
		if _, err := os.Stat(filepath.Join(reportDir, name)); err != nil {
			t.Fatalf("report artifact %s missing: %v", name, err)
		}
	}

	gatePath := filepath.Join(dir, "gate-result.json")
	exit, stdout, stderr = runCLIForSpec019("gate", "--out", gatePath, runDir)
	if exit != exitCannotVerify {
		t.Fatalf("gate exit=%d, want cannot_verify; stdout=%s stderr=%s", exit, stdout, stderr)
	}
	gate := readSpec019JSON(t, gatePath)
	assertSpec019Field(t, gate, "local_gate", "pass")
	assertSpec019Field(t, gate, "ci_witness_gate", "cannot_verify")
	assertSpec019Field(t, gate, "audit_grade_gate", "cannot_verify")
	assertSpec019Field(t, gate, "gate_mode", "observation")
}

func TestSpec019HarnessProofPack(t *testing.T) {
	dir := t.TempDir()
	fixtureDir := spec019ProofFixtureDir(t)
	copySpec019Fixture(t, fixtureDir, dir, "harness-profile.json")
	copySpec019Fixture(t, fixtureDir, dir, "harness-events-missing-model.jsonl")
	copySpec019Fixture(t, fixtureDir, dir, "harness-events-unsafe-raw-prompt.jsonl")
	restore := chdirSpec019(t, dir)
	defer restore()

	exit, stdout, stderr := runCLIForSpec019(
		"harness", "observe",
		"--profile", "harness-profile.json",
		"--source", "harness-events-missing-model.jsonl",
		"--out", "harness-run",
	)
	if exit != 0 {
		t.Fatalf("harness observe exit=%d stdout=%s stderr=%s", exit, stdout, stderr)
	}

	exit, stdout, stderr = runCLIForSpec019(
		"harness", "validate",
		"--profile", "harness-profile.json",
		"--run", "harness-run",
		"--out", "harness-validation.json",
	)
	if exit != exitCannotVerify {
		t.Fatalf("harness validate exit=%d, want cannot_verify; stdout=%s stderr=%s", exit, stdout, stderr)
	}
	validation := readSpec019JSON(t, "harness-validation.json")
	assertSpec019Field(t, validation, "validation_state", "not_assessed")
	assertSpec019Field(t, validation, "reason_code", "required_event_family_absent")

	exit, _, stderr = runCLIForSpec019(
		"harness", "observe",
		"--profile", "harness-profile.json",
		"--source", "harness-events-unsafe-raw-prompt.jsonl",
		"--out", "unsafe-run",
	)
	if exit == 0 || !strings.Contains(stderr, "unsafe_input:raw_prompt:forbidden_raw_field") {
		t.Fatalf("unsafe observe exit=%d stderr=%s", exit, stderr)
	}
	if _, err := os.Stat("unsafe-run"); !os.IsNotExist(err) {
		t.Fatalf("unsafe run dir exists or stat failed unexpectedly: %v", err)
	}
}

func TestSpec019TelemetryProofPack(t *testing.T) {
	fixture := filepath.Join("..", "..", "examples", "block21-cross-repo-posture", "valid-movement", "cross-repo-posture-export.json")
	exit, stdout, stderr := runCLIForSpec019(
		"export", "telemetry",
		"--profile", "prometheus-text-v1",
		"--cross-repo-posture", fixture,
		"--out", "-",
	)
	if exit != 0 {
		t.Fatalf("telemetry export exit=%d stdout=%s stderr=%s", exit, stdout, stderr)
	}
	for _, want := range []string{
		"# TYPE sdp_trace_posture_metric_numerator gauge",
		"sdp_trace_posture_metric_numerator",
		"sdp_trace_posture_movement_delta",
		"# EOF\n",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("telemetry output missing %q:\n%s", want, stdout)
		}
	}
	for _, forbidden := range []string{"https://", "access_token=", "secret", "/private"} {
		if strings.Contains(stdout, forbidden) {
			t.Fatalf("telemetry output leaked forbidden marker %q:\n%s", forbidden, stdout)
		}
	}
}

func TestSpec019HarnessFixtureDigest(t *testing.T) {
	fixturePath := filepath.Join(spec019ProofFixtureDir(t), "harness-events-missing-model.jsonl")
	line, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	var event map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(line), &event); err != nil {
		t.Fatal(err)
	}
	got, _ := event["source_digest"].(string)
	event["source_digest"] = ""
	canonical, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(canonical)
	if want := hex.EncodeToString(sum[:]); got != want {
		t.Fatalf("source_digest = %s, want %s", got, want)
	}
}

func runCLIForSpec019(args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run(args, &stdout, &stderr)
	return exit, stdout.String(), stderr.String()
}

func proofCommandForSpec019() []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c", "echo", "ok"}
	}
	return []string{"/bin/echo", "ok"}
}

func spec019ProofFixtureDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join("..", "..", "examples", "spec019-monitoring-gate-proof")
	if _, err := os.Stat(dir); err != nil {
		t.Fatal(err)
	}
	return dir
}

func readSpec019JSON(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func assertSpec019Field(t *testing.T, value map[string]any, key, want string) {
	t.Helper()
	if got, _ := value[key].(string); got != want {
		t.Fatalf("%s = %q, want %q in %#v", key, got, want, value)
	}
}

func copySpec019Fixture(t *testing.T, fromDir, toDir, name string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fromDir, name))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(toDir, name), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func chdirSpec019(t *testing.T, dir string) func() {
	t.Helper()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal(err)
		}
	}
}
