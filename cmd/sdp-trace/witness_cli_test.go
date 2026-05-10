package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWitnessCommandRequiresTargetRunDir(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"witness",
		"--kind", "github-actions",
		"--out", filepath.Join(t.TempDir(), "ci-witness.json"),
	}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("expected usage exit, got %d stdout=%q stderr=%q", exit, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "witness requires <runs-root-or-run-dir>") {
		t.Fatalf("missing usage message: %s", errOut.String())
	}
}

func TestWitnessCommandRequiresOutPath(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"witness", "--kind", "github-actions", root}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("expected usage exit, got %d stdout=%q stderr=%q", exit, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "witness requires --out <file>") {
		t.Fatalf("missing usage message: %s", errOut.String())
	}
}

func TestWitnessCommandBuildkiteRunMismatchReturnsFailExit(t *testing.T) {
	root := t.TempDir()
	runPath := writeWitnessRunJSON(t, filepath.Join(root, "001-agent-session"), "pipeline-42")
	envelopePath := writeBuildkiteRunMismatchEnvelope(t, t.TempDir(), root, runPath, "wrong-run-id")
	outPath := filepath.Join(t.TempDir(), "buildkite-witness.json")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"witness",
		"--kind", "buildkite",
		"--witness-envelope", envelopePath,
		"--out", outPath,
		root,
	}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("expected fail exit, got %d stdout=%q stderr=%q", exit, out.String(), errOut.String())
	}

	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	if !strings.Contains(string(raw), `"status": "fail"`) {
		t.Fatalf("expected fail status: %s", string(raw))
	}
}

func writeWitnessRunJSON(t *testing.T, runDir, runID string) string {
	t.Helper()
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run: %v", err)
	}
	runPath := filepath.Join(runDir, "run.json")
	runJSON, err := json.Marshal(map[string]string{
		"run_id": runID,
	})
	if err != nil {
		t.Fatalf("marshal run payload: %v", err)
	}
	if err := os.WriteFile(runPath, runJSON, 0o644); err != nil {
		t.Fatalf("write run: %v", err)
	}
	return runPath
}

func writeBuildkiteRunMismatchEnvelope(t *testing.T, envelopeDir, runsRoot, runPath, runID string) string {
	t.Helper()
	runJSON, err := os.ReadFile(runPath)
	if err != nil {
		t.Fatalf("read run: %v", err)
	}
	sum := sha256.Sum256(runJSON)
	envelopePath := filepath.Join(envelopeDir, "buildkite-envelope.json")
	relRunPath, err := filepath.Rel(runsRoot, runPath)
	if err != nil {
		t.Fatalf("resolve run path: %v", err)
	}
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
			"run_id":   runID,
			"job":      "verify",
		},
		"run_artifacts": []map[string]string{
			{
				"path":   filepath.ToSlash(relRunPath),
				"sha256": hex.EncodeToString(sum[:]),
			},
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
	return envelopePath
}
