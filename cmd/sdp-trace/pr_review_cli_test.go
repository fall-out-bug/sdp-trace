package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
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
			"raw_output_ref":  map[string]any{"id": "raw-run-code", "kind": "reviewer_output", "ref": "runs/run-code.txt", "digest_sha256": strings.Repeat("c", 64), "content_type": "text/plain", "redaction_state": "none"},
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
		"--out", filepath.Join(root, "packet-with-rest"),
		"--repo-id", "demo_repo",
		"--change-ref", "pr-123",
		"--base", strings.Repeat("a", 40),
		"--head", strings.Repeat("b", 40),
		"--diff", filepath.Join(root, "missing.diff"),
		"unexpected",
	}, &out, &errOut)
	if exit != exitUsage || !strings.Contains(errOut.String(), "accepts only flags") {
		t.Fatalf("packet rest exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{
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
			"raw_output_ref": map[string]any{"id": "raw-run-code", "kind": "reviewer_output", "ref": "runs/run-code.txt", "digest_sha256": strings.Repeat("c", 64), "content_type": "text/plain", "redaction_state": "none"},
			"findings":       []map[string]any{{"id": "F1", "severity": "major", "citation": map[string]any{"context_ref_id": "diff", "diff_hunk_id": "hunk-1"}, "summary": "Missing behavior."}},
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

func TestPRReviewHandlersKeepSubcommands(t *testing.T) {
	want := map[string]subcommandHandler{
		"packet":     runPRReviewPacket,
		"run":        runPRReviewRun,
		"synthesize": runPRReviewSynthesize,
		"validate":   runPRReviewValidate,
		"summarize":  runPRReviewSummarize,
		"check":      runPRReviewCheck,
	}
	if len(prReviewHandlers) != len(want) {
		t.Fatalf("prReviewHandlers length = %d, want %d", len(prReviewHandlers), len(want))
	}
	for name, wantHandler := range want {
		gotHandler, ok := prReviewHandlers[name]
		if !ok {
			t.Fatalf("prReviewHandlers missing %s", name)
		}
		if functionName(gotHandler) != functionName(wantHandler) {
			t.Fatalf("prReviewHandlers[%s] = %s, want %s", name, functionName(gotHandler), functionName(wantHandler))
		}
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	if exit := run([]string{"pr-review"}, &out, &errOut); exit != exitUsage || !strings.Contains(errOut.String(), "pr-review requires") {
		t.Fatalf("pr-review without subcommand exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func TestPRReviewPacketFlagsKeepContract(t *testing.T) {
	wantRequired := []requiredCLIFlag{
		{"out", "pr-review packet requires --out"},
		{"repo-id", "pr-review packet requires --repo-id"},
		{"change-ref", "pr-review packet requires --change-ref"},
		{"base", "pr-review packet requires --base"},
		{"head", "pr-review packet requires --head"},
		{"diff", "pr-review packet requires --diff"},
	}
	if len(prReviewPacketRequiredFlags) != len(wantRequired) {
		t.Fatalf("required flags length = %d, want %d", len(prReviewPacketRequiredFlags), len(wantRequired))
	}
	for i := range wantRequired {
		if prReviewPacketRequiredFlags[i] != wantRequired[i] {
			t.Fatalf("required flag %d = %#v, want %#v", i, prReviewPacketRequiredFlags[i], wantRequired[i])
		}
	}

	defaults := map[string]string{}
	for _, flag := range prReviewPacketStringFlags {
		defaults[flag.name] = flag.defaultValue
	}
	for _, name := range []string{"out", "repo-id", "change-ref", "base", "head", "diff", "metadata", "context", "verification"} {
		if got := defaults[name]; got != "" {
			t.Fatalf("default %s = %q, want empty", name, got)
		}
	}
	if got := defaults["ci-state"]; got != "not_assessed" {
		t.Fatalf("ci-state default = %q", got)
	}
	if got := defaults["created-by"]; got != "sdp-trace-cli" {
		t.Fatalf("created-by default = %q", got)
	}
}

func TestParsePRReviewPacketArgsKeepsUsageBoundaries(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n")
	validArgs := []string{
		"--out", filepath.Join(root, "packet"),
		"--repo-id", "demo_repo",
		"--change-ref", "pr-123",
		"--base", strings.Repeat("a", 40),
		"--head", strings.Repeat("b", 40),
		"--diff", diffPath,
	}
	var errOut bytes.Buffer
	opts, code, ok := parsePRReviewPacketArgs(validArgs, &errOut)
	if !ok || code != 0 || opts == nil {
		t.Fatalf("parse valid ok=%v code=%d err=%s", ok, code, errOut.String())
	}

	errOut.Reset()
	_, code, ok = parsePRReviewPacketArgs(append(validArgs, "unexpected"), &errOut)
	if ok || code != exitUsage || !strings.Contains(errOut.String(), "accepts only flags") {
		t.Fatalf("parse rest ok=%v code=%d err=%s", ok, code, errOut.String())
	}

	errOut.Reset()
	_, code, ok = parsePRReviewPacketArgs(validArgs[:len(validArgs)-2], &errOut)
	if ok || code != exitUsage || !strings.Contains(errOut.String(), "--diff") {
		t.Fatalf("parse missing diff ok=%v code=%d err=%s", ok, code, errOut.String())
	}

	var out bytes.Buffer
	errOut.Reset()
	badDiffArgs := append([]string{}, validArgs...)
	badDiffArgs[len(badDiffArgs)-1] = filepath.Join(root, "missing.diff")
	if exit := runPRReviewPacket(badDiffArgs, &out, &errOut); exit != exitCannotVerify {
		t.Fatalf("runPRReviewPacket bad diff exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func TestPRReviewPacketOptionsKeepsEvidenceMapping(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n")
	metadataPath := writeFileStringForPRReviewTest(t, root, "metadata.json", "{}\n")
	args := []string{
		"--out", filepath.Join(root, "packet"),
		"--repo-id", "demo_repo",
		"--change-ref", "pr-123",
		"--base", strings.Repeat("a", 40),
		"--head", strings.Repeat("b", 40),
		"--diff", diffPath,
		"--metadata", metadataPath,
		"--context", "context-one.md",
		"--context", "context-three.md",
		"--verification", "verify-one.txt",
		"--verification", "verify-two.txt",
		"--verification=verify-three.txt",
		"--ci-state", "pass",
		"--created-by", "tester",
	}
	var errOut bytes.Buffer
	opts, code, ok := parsePRReviewPacketArgs(args, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse args ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	options := prReviewPacketOptions(opts, args, opts.stringValue("out"))
	if options.OutDir != filepath.Join(root, "packet") ||
		options.RepoID != "demo_repo" ||
		options.ChangeRef != "pr-123" ||
		options.BaseCommit != strings.Repeat("a", 40) ||
		options.HeadCommit != strings.Repeat("b", 40) ||
		options.DiffPath != diffPath ||
		options.MetadataPath != metadataPath {
		t.Fatalf("identity options = %+v", options)
	}
	wantContext := []string{"context-one.md", "context-three.md"}
	if strings.Join(options.ContextPaths, "|") != strings.Join(wantContext, "|") {
		t.Fatalf("context paths = %#v, want %#v", options.ContextPaths, wantContext)
	}
	wantVerification := []string{"verify-one.txt", "verify-two.txt", "verify-three.txt"}
	if strings.Join(options.VerificationPaths, "|") != strings.Join(wantVerification, "|") {
		t.Fatalf("verification paths = %#v, want %#v", options.VerificationPaths, wantVerification)
	}
	if options.CIState != "pass" || options.CreatedBy != "tester" {
		t.Fatalf("metadata options = %+v", options)
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

func TestPRReviewRunPreviewUsesPacketAndProfile(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n-old\n+new\n")
	profilePath := writePRReviewCheckProfile(t, root)
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
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("packet exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{
		"pr-review", "run",
		"--packet", packetDir,
		"--profile", profilePath,
		"--out", filepath.Join(root, "runs"),
		"--work-dir", root,
		"--preview",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run preview exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"schema_version": "block30-pr-review-runs-v1"`) {
		t.Fatalf("preview output missing run schema: %s", out.String())
	}
}

func TestParsePRReviewRunArgsKeepsUsageBoundaries(t *testing.T) {
	var errOut bytes.Buffer
	opts, code, ok := parsePRReviewRunArgs([]string{"--packet", "packet", "--profile", "profile", "--out", "runs"}, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse run args ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	if opts.stringValue("work-dir") != "." ||
		opts.stringValue("allow-external-runner") != "" ||
		opts.stringValue("not-assessed-reason") != "" ||
		opts.boolValue("preview") {
		t.Fatalf("run defaults changed")
	}

	errOut.Reset()
	_, code, ok = parsePRReviewRunArgs([]string{"--unknown"}, &errOut)
	if ok || code != exitUsage {
		t.Fatalf("unknown flag ok=%v code=%d err=%s", ok, code, errOut.String())
	}

	errOut.Reset()
	_, code, ok = parsePRReviewRunArgs([]string{"--packet", "packet", "unexpected"}, &errOut)
	if ok || code != exitUsage || !strings.Contains(errOut.String(), "accepts only flags") {
		t.Fatalf("rest arg ok=%v code=%d err=%s", ok, code, errOut.String())
	}
}

func TestExecutePRReviewRunKeepsRunnerOptions(t *testing.T) {
	root := t.TempDir()
	diffPath := writeFileStringForPRReviewTest(t, root, "change.diff", "diff --git a/a.go b/a.go\n")
	profilePath := writeJSONForPRReviewTest(t, root, "profile.json", map[string]any{
		"schema_version":  "block30-pr-review-profile-v1",
		"profile_id":      "runner-preview",
		"required_planes": []string{"code_correctness"},
		"roles": []map[string]any{{
			"role_id":         "code",
			"plane":           "code_correctness",
			"runner":          "opencode",
			"requested_model": "not_assessed",
		}},
	})
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
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("packet exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}

	args := []string{
		"--packet", packetDir,
		"--profile", profilePath,
		"--out", filepath.Join(root, "runs"),
		"--work-dir", root,
		"--allow-external-runner", "opencode",
		"--allow-external-runner=pi",
		"--not-assessed-reason", "ci_model_review_not_configured",
		"--preview",
	}
	opts, code, ok := parsePRReviewRunArgs(args, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse preview ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	runs, preview, err := executePRReviewRun(opts, args)
	if err != nil {
		t.Fatalf("executePRReviewRun() error = %v", err)
	}
	if len(runs.Results) != 0 {
		t.Fatalf("preview produced run results: %+v", runs)
	}
	if preview == nil || len(preview.Roles) != 1 || preview.Roles[0].Runner != "opencode" {
		t.Fatalf("preview = %+v", preview)
	}

	runArgs := []string{
		"--packet", packetDir,
		"--profile", profilePath,
		"--out", filepath.Join(root, "not-assessed-runs"),
		"--work-dir", root,
		"--not-assessed-reason", "ci_model_review_not_configured",
	}
	runOpts, code, ok := parsePRReviewRunArgs(runArgs, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse not assessed run ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	notAssessedRuns, preview, err := executePRReviewRun(runOpts, runArgs)
	if err != nil {
		t.Fatalf("executePRReviewRun(not assessed) error = %v", err)
	}
	if preview != nil || len(notAssessedRuns.Results) != 1 {
		t.Fatalf("not assessed run preview=%+v runs=%+v", preview, notAssessedRuns)
	}
	result := notAssessedRuns.Results[0]
	if result.Status != prreview.StatusNotAssessed || result.Reason != "ci_model_review_not_configured" {
		t.Fatalf("not assessed override result = %+v", result)
	}

	badOpts, code, ok := parsePRReviewRunArgs([]string{
		"--packet", packetDir,
		"--profile", profilePath,
		"--out", filepath.Join(root, "bad-runs"),
		"--work-dir", filepath.Join(root, "missing"),
	}, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse bad workdir ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	if _, _, err := executePRReviewRun(badOpts, nil); err == nil || !strings.Contains(err.Error(), "work-dir:") || !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("missing work-dir error = %v", err)
	}

	out.Reset()
	errOut.Reset()
	exit = runPRReviewRun([]string{
		"--packet", filepath.Join(root, "missing-packet"),
		"--profile", profilePath,
		"--out", filepath.Join(root, "runs"),
	}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("missing packet run exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func TestWritePRReviewRunOutputKeepsPreviewBoundary(t *testing.T) {
	var out bytes.Buffer
	writePRReviewRunOutput(&out, prreview.RunSet{SchemaVersion: "run-set"}, nil)
	if !strings.Contains(out.String(), `"schema_version": "run-set"`) {
		t.Fatalf("run output = %s", out.String())
	}
	out.Reset()
	writePRReviewRunOutput(&out, prreview.RunSet{SchemaVersion: "run-set"}, &prreview.RunPreview{SchemaVersion: "preview"})
	if !strings.Contains(out.String(), `"schema_version": "preview"`) || strings.Contains(out.String(), "run-set") {
		t.Fatalf("preview output = %s", out.String())
	}
}

func TestParsePRReviewSynthesizeArgsKeepsUsageBoundaries(t *testing.T) {
	var errOut bytes.Buffer
	opts, code, ok := parsePRReviewSynthesizeArgs([]string{"--packet", "packet", "--runs", "runs", "--out", "ledger.json"}, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse synthesize args ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	if opts.stringValue("existing-ledger") != "" {
		t.Fatalf("existing ledger default changed")
	}

	errOut.Reset()
	_, code, ok = parsePRReviewSynthesizeArgs([]string{"--packet", "packet", "--runs", "runs"}, &errOut)
	if ok || code != exitUsage || !strings.Contains(errOut.String(), "requires --out") {
		t.Fatalf("missing out ok=%v code=%d err=%s", ok, code, errOut.String())
	}

	errOut.Reset()
	_, code, ok = parsePRReviewSynthesizeArgs([]string{"--packet", "packet", "--runs", "runs", "--out", "ledger.json", "unexpected"}, &errOut)
	if ok || code != exitUsage || !strings.Contains(errOut.String(), "accepts only flags") {
		t.Fatalf("rest arg ok=%v code=%d err=%s", ok, code, errOut.String())
	}
}

func TestReadPRReviewSynthesisInputsKeepsOptionalLedgerBoundary(t *testing.T) {
	root := t.TempDir()
	packetDir, packetDigest := writePRReviewPacketForSynthesisTest(t, root)
	runsDir := filepath.Join(root, "runs")
	writePRReviewRunsForSynthesisTest(t, runsDir, packetDigest)

	opts := &flagSet{name: "test synthesize inputs"}
	opts.setString("packet", packetDir)
	opts.setString("runs", runsDir)
	opts.setString("existing-ledger", "")
	inputs, err := readPRReviewSynthesisInputs(opts)
	if err != nil {
		t.Fatalf("read synthesis inputs without existing ledger: %v", err)
	}
	if inputs.packet.PacketDigest != packetDigest || inputs.runs.PacketDigest != packetDigest || inputs.existing != nil {
		t.Fatalf("inputs = %+v", inputs)
	}

	opts.setString("existing-ledger", filepath.Join(root, "missing-ledger.json"))
	if _, err := readPRReviewSynthesisInputs(opts); err == nil {
		t.Fatalf("missing existing ledger should fail")
	}
}

func TestRunPRReviewSynthesizeKeepsLedgerDurability(t *testing.T) {
	root := t.TempDir()
	packetDir, packetDigest := writePRReviewPacketForSynthesisTest(t, root)
	runsDir := filepath.Join(root, "runs")
	writePRReviewRunsForSynthesisTest(t, runsDir, packetDigest)
	ledgerPath := filepath.Join(root, "ledger.json")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := runPRReviewSynthesize([]string{"--packet", packetDir, "--runs", runsDir, "--out", ledgerPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("synthesize exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	written, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("read ledger: %v", err)
	}
	if !bytes.Equal(bytes.TrimSpace(out.Bytes()), bytes.TrimSpace(written)) {
		t.Fatalf("stdout should mirror durable ledger\nstdout=%s\nfile=%s", out.String(), string(written))
	}

	out.Reset()
	errOut.Reset()
	exit = runPRReviewSynthesize([]string{"--packet", filepath.Join(root, "missing-packet"), "--runs", runsDir, "--out", filepath.Join(root, "packet-fail.json")}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("missing packet exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = runPRReviewSynthesize([]string{"--packet", packetDir, "--runs", filepath.Join(root, "missing-runs"), "--out", filepath.Join(root, "runs-fail.json")}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("missing runs exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

func writePRReviewPacketForSynthesisTest(t *testing.T, root string) (string, string) {
	t.Helper()
	diffPath := writeFileStringForPRReviewTest(t, root, "synthesis.diff", "diff --git a/a.go b/a.go\n")
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
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("packet exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	var packet struct {
		PacketDigest string `json:"packet_digest"`
	}
	if err := json.Unmarshal(out.Bytes(), &packet); err != nil {
		t.Fatalf("packet stdout: %v\n%s", err, out.String())
	}
	return packetDir, packet.PacketDigest
}

func writePRReviewRunsForSynthesisTest(t *testing.T, runsDir, packetDigest string) {
	t.Helper()
	writeJSONForPRReviewTest(t, runsDir, "results.json", map[string]any{
		"schema_version": "block30-pr-review-runs-v1",
		"packet_digest":  packetDigest,
		"results": []map[string]any{{
			"review_run_id":   "run-code",
			"packet_digest":   packetDigest,
			"plane":           "code_correctness",
			"role_id":         "code",
			"runner":          "manual_external",
			"requested_model": "not_assessed",
			"observed_model":  "not_assessed",
			"model_family":    "not_assessed",
			"model_version":   "not_assessed",
			"status":          "no_findings",
			"raw_output_ref":  map[string]any{"id": "raw-run-code", "kind": "reviewer_output", "ref": "runs/run-code.txt", "digest_sha256": strings.Repeat("c", 64), "content_type": "text/plain", "redaction_state": "none"},
			"findings":        []map[string]any{},
		}},
	})
}

func TestParsePRReviewValidateArgsKeepsUsageBoundaries(t *testing.T) {
	var errOut bytes.Buffer
	opts, code, ok := parsePRReviewValidateArgs([]string{"--packet", "packet", "--profile", "profile", "--runs", "runs", "--ledger", "ledger", "--out", "validation.json"}, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse validate args ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	if opts.stringValue("packet") != "packet" || opts.stringValue("ledger") != "ledger" {
		t.Fatalf("validate opts changed: %+v", opts)
	}

	errOut.Reset()
	_, code, ok = parsePRReviewValidateArgs([]string{"--packet", "packet", "--profile", "profile", "--runs", "runs", "--ledger", "ledger"}, &errOut)
	if ok || code != exitUsage || !strings.Contains(errOut.String(), "requires --out") {
		t.Fatalf("missing out ok=%v code=%d err=%s", ok, code, errOut.String())
	}

	errOut.Reset()
	_, code, ok = parsePRReviewValidateArgs([]string{"--packet", "packet", "--profile", "profile", "--runs", "runs", "--ledger", "ledger", "--out", "validation.json", "unexpected"}, &errOut)
	if ok || code != exitUsage || !strings.Contains(errOut.String(), "accepts only flags") {
		t.Fatalf("rest arg ok=%v code=%d err=%s", ok, code, errOut.String())
	}
}

func TestReadPRReviewValidationInputsKeepsArtifactBoundaries(t *testing.T) {
	root := t.TempDir()
	artifacts := writePRReviewValidationArtifactsForTest(t, root, true)

	opts := newPRReviewValidationTestOptions(artifacts.packetPath, artifacts.profilePath, artifacts.runsPath, artifacts.ledgerPath)
	inputs, err := readPRReviewValidationInputs(opts)
	if err != nil {
		t.Fatalf("read validation inputs: %v", err)
	}
	if inputs.packet.PacketDigest != artifacts.packetDigest ||
		inputs.profile.ProfileID != "default" ||
		inputs.runs.PacketDigest != artifacts.packetDigest ||
		inputs.ledger.PacketDigest != artifacts.packetDigest {
		t.Fatalf("inputs = %+v", inputs)
	}

	for _, test := range []struct {
		name string
		opts *flagSet
	}{
		{name: "packet", opts: newPRReviewValidationTestOptions(filepath.Join(root, "missing-packet"), artifacts.profilePath, artifacts.runsPath, artifacts.ledgerPath)},
		{name: "profile", opts: newPRReviewValidationTestOptions(artifacts.packetPath, filepath.Join(root, "missing-profile.json"), artifacts.runsPath, artifacts.ledgerPath)},
		{name: "runs", opts: newPRReviewValidationTestOptions(artifacts.packetPath, artifacts.profilePath, filepath.Join(root, "missing-runs"), artifacts.ledgerPath)},
		{name: "ledger", opts: newPRReviewValidationTestOptions(artifacts.packetPath, artifacts.profilePath, artifacts.runsPath, filepath.Join(root, "missing-ledger.json"))},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := readPRReviewValidationInputs(test.opts); err == nil {
				t.Fatalf("missing %s artifact should fail", test.name)
			}
		})
	}
}

func TestRunPRReviewValidateKeepsDurableVerdictAndExitMapping(t *testing.T) {
	root := t.TempDir()
	artifacts := writePRReviewValidationArtifactsForTest(t, root, true)
	validationPath := filepath.Join(root, "validation.json")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := runPRReviewValidate([]string{"--packet", artifacts.packetPath, "--profile", artifacts.profilePath, "--runs", artifacts.runsPath, "--ledger", artifacts.ledgerPath, "--out", validationPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
	written, err := os.ReadFile(validationPath)
	if err != nil {
		t.Fatalf("read validation: %v", err)
	}
	if !bytes.Equal(bytes.TrimSpace(out.Bytes()), bytes.TrimSpace(written)) {
		t.Fatalf("stdout should mirror durable validation\nstdout=%s\nfile=%s", out.String(), string(written))
	}

	unresolved := writePRReviewValidationArtifactsForTest(t, filepath.Join(root, "unresolved"), false)
	out.Reset()
	errOut.Reset()
	exit = runPRReviewValidate([]string{"--packet", unresolved.packetPath, "--profile", unresolved.profilePath, "--runs", unresolved.runsPath, "--ledger", unresolved.ledgerPath, "--out", filepath.Join(root, "unresolved-validation.json")}, &out, &errOut)
	if exit != exitCannotVerify || !strings.Contains(out.String(), `"review_coverage_state": "coverage_unresolved"`) {
		t.Fatalf("unresolved validation exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = runPRReviewValidate([]string{"--packet", filepath.Join(root, "missing-packet"), "--profile", artifacts.profilePath, "--runs", artifacts.runsPath, "--ledger", artifacts.ledgerPath, "--out", filepath.Join(root, "missing-validation.json")}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("missing packet validation exit=%d err=%s out=%s", exit, errOut.String(), out.String())
	}
}

type prReviewValidationTestArtifacts struct {
	packetDigest string
	packetPath   string
	profilePath  string
	runsPath     string
	ledgerPath   string
}

func writePRReviewValidationArtifactsForTest(t *testing.T, root string, resolved bool) prReviewValidationTestArtifacts {
	t.Helper()
	packetDigest := "sha256:" + strings.Repeat("7", 64)
	packetPath := writeJSONForPRReviewTest(t, root, "packet.json", map[string]any{
		"schema_version":  "block30-pr-review-packet-v1",
		"packet_id":       "packet-1",
		"packet_digest":   packetDigest,
		"repo_id":         "demo_repo",
		"change_ref":      "pr-123",
		"base_commit":     strings.Repeat("a", 40),
		"head_commit":     strings.Repeat("b", 40),
		"diff_ref":        map[string]any{"id": "diff", "kind": "diff", "ref": "inputs/diff.patch", "digest_sha256": strings.Repeat("2", 64), "content_type": "text/x-diff", "redaction_state": "none"},
		"ci_state":        "not_assessed",
		"created_at":      "2026-05-09T12:00:00Z",
		"created_by":      "test",
		"redaction_state": "none",
	})
	profilePath := writeJSONForPRReviewTest(t, root, "profile.json", map[string]any{
		"schema_version":  "block30-pr-review-profile-v1",
		"profile_id":      "default",
		"required_planes": []string{"code_correctness"},
		"roles":           []map[string]any{{"role_id": "code", "plane": "code_correctness", "runner": "manual_external", "requested_model": "not_assessed"}},
	})
	status := "no_findings"
	findings := []map[string]any{}
	disposition := "resolved"
	if !resolved {
		status = "findings_reported"
		findings = []map[string]any{{"id": "F1", "severity": "major", "citation": map[string]any{"context_ref_id": "diff", "diff_hunk_id": "hunk-1"}, "summary": "Missing behavior."}}
		disposition = "unresolved_review_blocker"
	}
	runsPath := writeJSONForPRReviewTest(t, root, "results.json", map[string]any{
		"schema_version": "block30-pr-review-runs-v1",
		"packet_digest":  packetDigest,
		"results": []map[string]any{{
			"review_run_id":   "run-code",
			"packet_digest":   packetDigest,
			"plane":           "code_correctness",
			"role_id":         "code",
			"runner":          "manual_external",
			"requested_model": "not_assessed",
			"observed_model":  "not_assessed",
			"model_family":    "not_assessed",
			"model_version":   "not_assessed",
			"status":          status,
			"raw_output_ref":  map[string]any{"id": "raw-run-code", "kind": "reviewer_output", "ref": "runs/run-code.txt", "digest_sha256": strings.Repeat("c", 64), "content_type": "text/plain", "redaction_state": "none"},
			"findings":        findings,
		}},
	})
	ledgerFindings := []map[string]any{}
	if !resolved {
		ledgerFindings = []map[string]any{{
			"id":            "F1",
			"review_run_id": "run-code",
			"plane":         "code_correctness",
			"role_id":       "code",
			"severity":      "major",
			"summary":       "Missing behavior.",
			"citation":      map[string]any{"context_ref_id": "diff", "diff_hunk_id": "hunk-1"},
			"disposition":   disposition,
		}}
	}
	ledgerPath := writeJSONForPRReviewTest(t, root, "ledger.json", map[string]any{
		"schema_version": "block30-pr-review-ledger-v1",
		"packet_digest":  packetDigest,
		"findings":       ledgerFindings,
	})
	return prReviewValidationTestArtifacts{packetDigest: packetDigest, packetPath: packetPath, profilePath: profilePath, runsPath: runsPath, ledgerPath: ledgerPath}
}

func newPRReviewValidationTestOptions(packetPath, profilePath, runsPath, ledgerPath string) *flagSet {
	opts := &flagSet{name: "test validation inputs"}
	opts.setString("packet", packetPath)
	opts.setString("profile", profilePath)
	opts.setString("runs", runsPath)
	opts.setString("ledger", ledgerPath)
	return opts
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
