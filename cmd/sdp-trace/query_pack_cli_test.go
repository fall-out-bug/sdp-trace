package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQueryPackRequiresExplicitPackAndOut(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"query-pack", "--run", t.TempDir()}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("exit = %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(errOut.String(), "error: ambiguous pack selection; --pack is required") {
		t.Fatalf("missing pack error: %s", errOut.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"query-pack", "--pack", "forensics-basic-v1", "--run", t.TempDir()}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("exit = %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(errOut.String(), "query-pack requires --out") {
		t.Fatalf("missing out error: %s", errOut.String())
	}
}

func TestQueryPackWritesSafeArtifactAndExplain(t *testing.T) {
	runDir := writeQueryPackCLIFixture(t)
	outPath := filepath.Join(t.TempDir(), "query-pack-result.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"query-pack",
		"--pack", "forensics-basic-v1",
		"--run", runDir,
		"--out", outPath,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("query-pack exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if strings.TrimSpace(out.String()) != "" {
		t.Fatalf("query-pack wrote to stdout: %s", out.String())
	}
	payload, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if strings.Contains(string(payload), "secret-token") ||
		strings.Contains(string(payload), "deploy-prod.sh") {
		t.Fatalf("query-pack leaked unsafe marker")
	}
	var envelope struct {
		SchemaVersion string `json:"schema_version"`
		QueryPackID   string `json:"query_pack_id"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if envelope.QueryPackID != "forensics-basic-v1" {
		t.Fatalf("envelope = %+v", envelope)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"query-pack", "explain", "--result", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("explain exit=%d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "forensics-redactions") ||
		strings.Contains(out.String(), "secret-token") ||
		strings.Contains(out.String(), "deploy-prod.sh") {
		t.Fatalf("unsafe or incomplete explain output")
	}
}

func writeQueryPackCLIFixture(t *testing.T) string {
	t.Helper()
	runDir := t.TempDir()
	writeCLITestJSON(t, filepath.Join(runDir, "run.json"), map[string]any{
		"run_id":          "cli-run",
		"run_nonce":       "cli-nonce",
		"source_baseline": "sha256:cli",
		"event_refs": []map[string]any{
			{"sequence": 1, "event_hash": "hash", "event_type": "command_started", "uri": "events/deploy-prod.sh"},
		},
		"verifier_states": map[string]any{
			"event_chain_structurally_valid": map[string]any{"state": "pass"},
		},
	})
	writeCLITestJSON(t, filepath.Join(runDir, "forensic-retention.assessment-result.json"), map[string]any{
		"schema_version":                "block18-forensic-retention-assessment-v1",
		"forensic_retention_assessment": "pass",
		"forensic_conditions": []map[string]any{
			{"id": "redaction_unresolved_visible", "state": "pass", "reason_code": "redaction_resolved", "reason": "resolved"},
		},
	})
	writeCLITestJSON(t, filepath.Join(runDir, "adapter-capture.assessment-result.json"), map[string]any{
		"schema_version":             "block19-adapter-capture-assessment-v1",
		"adapter_capture_assessment": "pass",
		"adapter_capture_conditions": []map[string]any{
			{"id": "provider_refs_portable", "state": "pass", "reason_code": "provider_refs_portable", "reason": "portable"},
		},
	})
	return runDir
}

func writeCLITestJSON(t *testing.T, path string, value any) {
	t.Helper()
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
}
