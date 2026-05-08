package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCrossRepoPostureExportWritesSafeArtifactAndExplain(t *testing.T) {
	root := t.TempDir()
	withCLIChdir(t, root)
	current := writePostureCLIQueryPack(t, ".", "current", "missing_telemetry")
	previous := writePostureCLIQueryPack(t, ".", "previous", "present")
	selectionPath := writePostureCLISelection(t, ".", current, previous)
	outPath := "cross-repo-export.json"

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"export", "cross-repo-posture",
		"--profile", "cross-repo-evidence-posture-v1",
		"--selection", selectionPath,
		"--out", outPath,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("export exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if strings.TrimSpace(out.String()) != "" {
		t.Fatalf("export wrote to stdout: %s", out.String())
	}
	payload, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	if strings.Contains(string(payload), "secret-token") ||
		strings.Contains(string(payload), "/private") ||
		strings.Contains(string(payload), "https://") {
		t.Fatalf("export leaked unsafe marker")
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"export", "cross-repo-posture", "explain", "--result", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("explain exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "movement_summary") ||
		strings.Contains(out.String(), "secret-token") ||
		strings.Contains(out.String(), "/private") {
		t.Fatalf("unsafe or incomplete explain output: %s", out.String())
	}
}

func TestCrossRepoPostureValidateOnlyAndRequiredFlags(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"export", "cross-repo-posture", "--profile", "cross-repo-evidence-posture-v1"}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("exit=%d err=%s", exit, errOut.String())
	}
	if !strings.Contains(errOut.String(), "requires --selection") {
		t.Fatalf("missing selection error: %s", errOut.String())
	}

	root := t.TempDir()
	withCLIChdir(t, root)
	current := writePostureCLIQueryPack(t, ".", "current", "present")
	previous := writePostureCLIQueryPack(t, ".", "previous", "present")
	selectionPath := writePostureCLISelection(t, ".", current, previous)
	out.Reset()
	errOut.Reset()
	exit = run([]string{
		"export", "cross-repo-posture",
		"--profile", "cross-repo-evidence-posture-v1",
		"--selection", selectionPath,
		"--validate-only",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate-only exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func writePostureCLISelection(t *testing.T, root, current, previous string) string {
	t.Helper()
	currentDigest := writePostureCLIDigest(t, current)
	previousDigest := writePostureCLIDigest(t, previous)
	currentSignals := filepath.Join(root, "current-signals.json")
	writeCLITestJSON(t, currentSignals, map[string]any{
		"schema_version": "block21-posture-signal-manifest-v1",
		"signals": []map[string]any{{
			"row_ref":         "timeline.0001",
			"witness_scope":   "ci_witnessed",
			"override_marker": "override_present",
		}},
	})
	selectionPath := filepath.Join(root, "selection.json")
	writeCLITestJSON(t, selectionPath, map[string]any{
		"schema_version":            "block21-cross-repo-selection-v1",
		"profile_id":                "cross-repo-evidence-posture-v1",
		"profile_version":           "v1",
		"grouping_set_id":           "repo_window_v1",
		"freshness_boundary":        "2026-01-01T00:00:00Z",
		"dimension_exposure_policy": []string{"repo", "team", "service", "harness", "change_type"},
		"current_window":            "2026-w02",
		"previous_window":           "2026-w01",
		"handoff":                   map[string]string{"consumer": "sdp-report"},
		"repositories": []map[string]any{
			{
				"input_id":                 "current",
				"repo":                     "repo-a",
				"team":                     "platform",
				"service":                  "api",
				"harness":                  "generic",
				"change_type":              "feature",
				"time_window":              "2026-w02",
				"input_observed_at":        "2026-01-05T00:00:00Z",
				"query_pack_result":        current,
				"artifact_digest_manifest": currentDigest,
				"posture_signal_manifest":  currentSignals,
			},
			{
				"input_id":                 "previous",
				"repo":                     "repo-a",
				"team":                     "platform",
				"service":                  "api",
				"harness":                  "generic",
				"change_type":              "feature",
				"time_window":              "2026-w01",
				"input_observed_at":        "2026-01-05T00:00:00Z",
				"query_pack_result":        previous,
				"artifact_digest_manifest": previousDigest,
			},
		},
	})
	return selectionPath
}

func writePostureCLIQueryPack(t *testing.T, root, name, state string) string {
	t.Helper()
	path := filepath.Join(root, name+"-query-pack.json")
	writeCLITestJSON(t, path, map[string]any{
		"schema_version":     "block20-forensics-query-pack-result-v1",
		"query_pack_id":      "forensics-basic-v1",
		"query_pack_version": "v1",
		"input_artifacts": []map[string]any{
			{"role": "run", "path_redacted_id": "run", "artifact_required": true},
		},
		"query_rows": map[string]any{
			"forensics-summary": []map[string]any{},
			"forensics-timeline": []map[string]any{
				{"id": "timeline.0001", "query": "forensics-timeline", "evidence_state": state, "evidence_family": "command", "source_ref": "block_09.event.command.e0001"},
				{"id": "timeline.0002", "query": "forensics-timeline", "evidence_state": "not_assessed", "evidence_family": "test", "source_ref": "block_09.event.test.e0002"},
			},
			"forensics-gaps":              []map[string]any{},
			"forensics-redactions":        []map[string]any{},
			"forensics-capture-depth":     []map[string]any{},
			"forensics-unverified-claims": []map[string]any{},
		},
		"output_safety": map[string]any{"verified_absent_sensitive_classes": []string{"tokens"}},
	})
	return path
}

func writePostureCLIDigest(t *testing.T, artifact string) string {
	t.Helper()
	payload, err := os.ReadFile(artifact)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(payload)
	path := artifact + ".digest.json"
	writeCLITestJSON(t, path, map[string]any{
		"schema_version": "block21-artifact-digest-manifest-v1",
		"artifacts": []map[string]any{{
			"role":   "query_pack_result",
			"path":   filepath.Base(artifact),
			"sha256": hex.EncodeToString(sum[:]),
		}},
	})
	return path
}

func withCLIChdir(t *testing.T, dir string) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
}
