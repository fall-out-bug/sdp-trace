package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPRReviewPacketSynthesizeValidateSummarizeCLI(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	contextOne := writeFileStringForPRReviewTest(t, root, "context-one.md", "# Context one\n")
	contextTwo := writeFileStringForPRReviewTest(t, root, "context-two.md", "# Context two\n")
	verificationOne := writeFileStringForPRReviewTest(t, root, "verify-one.txt", "go test ./...\n")
	verificationTwo := writeFileStringForPRReviewTest(t, root, "verify-two.txt", "jq empty schema/*.json\n")
	profilePath := writeJSONForPRReviewTest(t, root, "profile.json", map[string]any{
		"schema_version":  "block30-pr-review-profile-v1",
		"profile_id":      "trust-sensitive-default",
		"required_planes": []string{"code_correctness"},
		"roles": []map[string]any{{
			"role_id":         "code",
			"plane":           "code_correctness",
			"runner":          "manual_external",
			"requested_model": "not_assessed",
		}},
	})
	runsPath := filepath.Join(root, "runs", "results.json")
	packetDir := filepath.Join(root, "packet")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"pr-review", "packet",
		"--out", packetDir,
		"--repo-id", "demo_repo",
		"--change-ref", "pr-123",
		"--base", strings.Repeat("a", 40),
		"--head", strings.Repeat("b", 40),
		"--diff", diffPath,
		"--context", contextOne,
		"--context", contextTwo,
		"--verification", verificationOne,
		"--verification", verificationTwo,
		"--ci-state", "not_assessed",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("packet exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	var packet struct {
		PacketDigest     string `json:"packet_digest"`
		ContextRefs      []any  `json:"context_refs"`
		VerificationRefs []any  `json:"verification_refs"`
	}
	if err := json.Unmarshal(out.Bytes(), &packet); err != nil {
		t.Fatalf("packet stdout: %v\n%s", err, out.String())
	}
	if packet.PacketDigest == "" {
		t.Fatalf("missing packet digest: %s", out.String())
	}
	if len(packet.ContextRefs) != 2 || len(packet.VerificationRefs) != 2 {
		t.Fatalf("packet did not preserve repeated refs: %s", out.String())
	}
	writeJSONForPRReviewTest(t, filepath.Dir(runsPath), filepath.Base(runsPath), map[string]any{
		"schema_version": "block30-pr-review-runs-v1",
		"packet_digest":  packet.PacketDigest,
		"results": []map[string]any{{
			"review_run_id":   "run-code",
			"packet_digest":   packet.PacketDigest,
			"plane":           "code_correctness",
			"role_id":         "code",
			"runner":          "manual_external",
			"requested_model": "not_assessed",
			"observed_model":  "not_assessed",
			"model_family":    "not_assessed",
			"model_version":   "not_assessed",
			"status":          "no_findings",
			"findings":        []map[string]any{},
		}},
	})
	out.Reset()
	errOut.Reset()
	ledgerPath := filepath.Join(root, "ledger.json")
	exit = run([]string{"pr-review", "synthesize", "--packet", packetDir, "--runs", filepath.Dir(runsPath), "--out", ledgerPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("synthesize exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	validationPath := filepath.Join(root, "validation.json")
	exit = run([]string{"pr-review", "validate", "--packet", packetDir, "--profile", profilePath, "--runs", filepath.Dir(runsPath), "--ledger", ledgerPath, "--out", validationPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"review_coverage_state": "coverage_satisfied"`) ||
		!strings.Contains(out.String(), `"merge_decision": "not_authorized_by_sdp_trace"`) {
		t.Fatalf("validation missing authority/coverage: %s", out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"pr-review", "summarize", "--validation", validationPath, "--ledger", ledgerPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("summarize exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "not_authorized_by_sdp_trace") {
		t.Fatalf("summary missing authority boundary: %s", out.String())
	}
	if strings.Contains(strings.ToLower(out.String()), "safe to merge") {
		t.Fatalf("summary implies merge approval: %s", out.String())
	}
}

func TestPRReviewCLIRequiresPacketInputsAndReturnsNonZeroForUnresolved(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"pr-review", "packet",
		"--out", filepath.Join(root, "packet"),
		"--repo-id", "demo_repo",
		"--change-ref", "pr-123",
		"--base", strings.Repeat("a", 40),
		"--head", strings.Repeat("b", 40),
	}, &out, &errOut)
	if exit != exitUsage || !strings.Contains(errOut.String(), "--diff") {
		t.Fatalf("missing diff exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}

	packetDigest := "sha256:" + strings.Repeat("1", 64)
	packetPath := writeJSONForPRReviewTest(t, root, "packet.json", map[string]any{
		"schema_version":    "block30-pr-review-packet-v1",
		"packet_id":         "packet-1",
		"packet_digest":     packetDigest,
		"repo_id":           "demo_repo",
		"change_ref":        "pr-123",
		"base_commit":       strings.Repeat("a", 40),
		"head_commit":       strings.Repeat("b", 40),
		"diff_ref":          map[string]any{"id": "diff", "kind": "diff", "ref": "inputs/diff.patch", "digest_sha256": strings.Repeat("2", 64), "content_type": "unified_diff", "redaction_state": "none"},
		"context_refs":      []map[string]any{},
		"verification_refs": []map[string]any{},
		"ci_state":          "not_assessed",
		"created_at":        "2026-05-09T12:00:00Z",
		"created_by":        "test",
		"redaction_state":   "none",
	})
	profilePath := writeJSONForPRReviewTest(t, root, "profile.json", map[string]any{
		"schema_version":  "block30-pr-review-profile-v1",
		"profile_id":      "default",
		"required_planes": []string{"code_correctness"},
		"roles":           []map[string]any{{"role_id": "code", "plane": "code_correctness", "runner": "manual_external", "requested_model": "not_assessed"}},
	})
	runsPath := writeJSONForPRReviewTest(t, root, "results.json", map[string]any{
		"schema_version": "block30-pr-review-runs-v1",
		"packet_digest":  packetDigest,
		"results": []map[string]any{{
			"review_run_id": "run-code", "packet_digest": packetDigest, "plane": "code_correctness", "role_id": "code", "runner": "manual_external",
			"requested_model": "not_assessed", "observed_model": "not_assessed", "model_family": "not_assessed", "model_version": "not_assessed", "status": "findings_reported",
			"findings": []map[string]any{{"id": "F1", "severity": "major", "citation": map[string]any{"context_ref_id": "diff", "diff_hunk_id": "hunk-1"}, "summary": "Missing behavior."}},
		}},
	})
	ledgerPath := writeJSONForPRReviewTest(t, root, "ledger.json", map[string]any{
		"schema_version": "block30-pr-review-ledger-v1",
		"packet_digest":  packetDigest,
		"findings": []map[string]any{{
			"id": "F1", "review_run_id": "run-code", "plane": "code_correctness", "role_id": "code", "severity": "major",
			"summary": "Missing behavior.", "citation": map[string]any{"context_ref_id": "diff", "diff_hunk_id": "hunk-1"}, "disposition": "unresolved_review_blocker",
		}},
	})
	out.Reset()
	errOut.Reset()
	validationPath := filepath.Join(root, "unresolved-validation.json")
	exit = run([]string{"pr-review", "validate", "--packet", packetPath, "--profile", profilePath, "--runs", runsPath, "--ledger", ledgerPath, "--out", validationPath}, &out, &errOut)
	if exit != exitCannotVerify || !strings.Contains(out.String(), `"review_coverage_state": "coverage_unresolved"`) {
		t.Fatalf("unresolved validation exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func TestPRReviewCheckWritesRunProvenance(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	profilePath := writeJSONForPRReviewTest(t, root, "profile.json", map[string]any{
		"schema_version":  "block30-pr-review-profile-v1",
		"profile_id":      "default",
		"required_planes": []string{"code_correctness"},
		"roles":           []map[string]any{{"role_id": "code", "plane": "code_correctness", "runner": "manual_external", "requested_model": "not_assessed"}},
	})
	outDir := filepath.Join(root, "review")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"pr-review", "check",
		"--out", outDir,
		"--repo-id", "demo_repo",
		"--change-ref", "pr-123",
		"--base", strings.Repeat("a", 40),
		"--head", strings.Repeat("b", 40),
		"--diff", diffPath,
		"--profile", profilePath,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("check exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if _, err := os.Stat(filepath.Join(outDir, "runs", "results.json")); err != nil {
		t.Fatalf("check did not persist run provenance: %v", err)
	}
}

func TestPRReviewCheckPreviewDoesNotWriteArtifacts(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	profilePath := writePRReviewCheckProfile(t, root)
	outDir := filepath.Join(root, "review")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run(prReviewCheckArgs(outDir, diffPath, profilePath, "--preview"), &out, &errOut)
	if exit != 0 {
		t.Fatalf("check preview exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"schema_version": "block30-pr-review-runs-v1"`) {
		t.Fatalf("preview output missing run schema: %s", out.String())
	}
	for _, path := range []string{
		filepath.Join(outDir, "runs", "results.json"),
		filepath.Join(outDir, "ledger.json"),
		filepath.Join(outDir, "validation.json"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("preview wrote artifact %s: %v", path, err)
		}
	}
}

func TestPRReviewCheckRequiresOutAndPacketInputs(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	profilePath := writePRReviewCheckProfile(t, root)
	for name, tc := range map[string]struct {
		args    []string
		wantErr string
	}{
		"missing-out": {
			args: []string{
				"pr-review", "check",
				"--repo-id", "demo_repo",
				"--change-ref", "pr-123",
				"--base", strings.Repeat("a", 40),
				"--head", strings.Repeat("b", 40),
				"--diff", diffPath,
				"--profile", profilePath,
			},
			wantErr: "pr-review check requires --out",
		},
		"missing-diff": {
			args: []string{
				"pr-review", "check",
				"--out", filepath.Join(root, "review-missing-diff"),
				"--repo-id", "demo_repo",
				"--change-ref", "pr-123",
				"--base", strings.Repeat("a", 40),
				"--head", strings.Repeat("b", 40),
				"--profile", profilePath,
			},
			wantErr: "pr-review packet requires --diff",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run(tc.args, &out, &errOut)
			if exit != exitUsage || !strings.Contains(errOut.String(), tc.wantErr) {
				t.Fatalf("check exit=%d err=%s out=%s want %q", exit, errOut.String(), out.String(), tc.wantErr)
			}
		})
	}
}

func TestPRReviewCheckRejectsMissingOrFileWorkDir(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	profilePath := writePRReviewCheckProfile(t, root)
	fileWorkDir := writeFileStringForPRReviewTest(t, root, "work-file", "not a dir\n")
	for name, workDir := range map[string]string{
		"missing-work-dir": filepath.Join(root, "missing"),
		"file-work-dir":    fileWorkDir,
	} {
		t.Run(name, func(t *testing.T) {
			var out bytes.Buffer
			var errOut bytes.Buffer
			args := prReviewCheckArgs(filepath.Join(root, "review-"+name), diffPath, profilePath, "--work-dir", workDir)
			exit := run(args, &out, &errOut)
			if exit != exitCannotVerify || !strings.Contains(errOut.String(), "work-dir") {
				t.Fatalf("check exit=%d err=%s out=%s", exit, errOut.String(), out.String())
			}
		})
	}
}

func writePRReviewCheckProfile(t *testing.T, root string) string {
	t.Helper()
	return writeJSONForPRReviewTest(t, root, "profile.json", map[string]any{
		"schema_version":  "block30-pr-review-profile-v1",
		"profile_id":      "default",
		"required_planes": []string{"code_correctness"},
		"roles":           []map[string]any{{"role_id": "code", "plane": "code_correctness", "runner": "manual_external", "requested_model": "not_assessed"}},
	})
}

func prReviewCheckArgs(outDir, diffPath, profilePath string, extra ...string) []string {
	args := []string{
		"pr-review", "check",
		"--out", outDir,
		"--repo-id", "demo_repo",
		"--change-ref", "pr-123",
		"--base", strings.Repeat("a", 40),
		"--head", strings.Repeat("b", 40),
		"--diff", diffPath,
		"--profile", profilePath,
	}
	return append(args, extra...)
}

func writeFileStringForPRReviewTest(t *testing.T, root, name, content string) string {
	t.Helper()
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeJSONForPRReviewTest(t *testing.T, root, name string, value any) string {
	t.Helper()
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
