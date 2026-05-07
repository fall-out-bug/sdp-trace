package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/forensic"
)

func TestForensicAssessRequiresInputsWithoutWriting(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"assess", "--profile", "forensic-retention", "--out", outPath}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("forensic assess missing inputs exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Fatalf("forensic assess wrote artifact despite usage error")
	}
}

func TestForensicAssessPassesAndExplains(t *testing.T) {
	root := t.TempDir()
	paths := writeForensicFixtureInputs(t, root)
	outPath := filepath.Join(root, "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "forensic-retention",
		"--out", outPath,
		"--run", paths.run,
		"--redaction-policy", paths.policy,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("forensic assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if strings.Contains(out.String(), "secret-token") || strings.Contains(out.String(), "raw prompt") || strings.Contains(out.String(), "--password") {
		t.Fatalf("forensic assess leaked sensitive marker: %s", out.String())
	}
	var result forensic.AssessmentResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("assessment payload: %v", err)
	}
	if result.ForensicRetentionAssessment != forensic.StatePass || result.SchemaVersion != forensic.SchemaVersion {
		t.Fatalf("assessment result = %+v", result)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"assess", "explain", "--assessment-result", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("forensic explain exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "Forensic retention assessment: pass") ||
		!strings.Contains(out.String(), "Forensic condition raw_reference_bound: pass") {
		t.Fatalf("forensic explain missing fields: %s", out.String())
	}
}

func TestForensicAssessRejectsDigestOnlyAndWeakDigest(t *testing.T) {
	root := t.TempDir()
	paths := writeForensicFixtureInputs(t, root)
	var runEvidence forensic.RunEvidence
	readTestJSON(t, filepath.Join(paths.run, "run.json"), &runEvidence)
	runEvidence.Events[0].RetentionMode = forensic.RetentionModeDigestOnly
	runEvidence.Events[0].RawReference = nil
	runEvidence.Events[1].RawReference.Digest.Algorithm = "sha1"
	writeTestJSON(t, filepath.Join(paths.run, "run.json"), runEvidence)

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "forensic-retention",
		"--out", filepath.Join(root, "assessment.json"),
		"--run", paths.run,
		"--redaction-policy", paths.policy,
	}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("forensic assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"reason_code": "critical_evidence_digest_only"`) ||
		!strings.Contains(out.String(), `"reason_code": "weak_digest"`) {
		t.Fatalf("forensic assess missing fail reasons: %s", out.String())
	}
}

func TestForensicAssessPreviewDoesNotWriteOrLeak(t *testing.T) {
	root := t.TempDir()
	outPath := filepath.Join(root, "assessment.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"assess", "preview", "--profile", "forensic-retention", "--out", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("forensic preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"selected_profile": "forensic_retention"`) ||
		!strings.Contains(out.String(), `"claim": "preview is read-only and does not emit a forensic verdict"`) {
		t.Fatalf("forensic preview missing fields: %s", out.String())
	}
	if strings.Contains(out.String(), "secret-token") || strings.Contains(out.String(), "raw prompt") {
		t.Fatalf("forensic preview leaked sensitive marker: %s", out.String())
	}
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Fatalf("forensic preview wrote artifact")
	}
}

func TestForensicAssessPreviewCannotVerifyMalformedProvidedInputs(t *testing.T) {
	root := t.TempDir()
	badPolicy := filepath.Join(root, "bad-policy.json")
	if err := os.WriteFile(badPolicy, []byte(`{not-json`), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"assess", "preview", "--profile", "forensic-retention", "--redaction-policy", badPolicy}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("forensic preview malformed input exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"redaction_policy": "present_malformed"`) {
		t.Fatalf("forensic preview missing malformed status: %s", out.String())
	}
}

type forensicFixturePaths struct {
	run    string
	policy string
}

func writeForensicFixtureInputs(t *testing.T, dir string) forensicFixturePaths {
	t.Helper()
	runDir := filepath.Join(dir, "run")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run: %v", err)
	}
	input := forensic.ValidTestInput()
	paths := forensicFixturePaths{
		run:    runDir,
		policy: filepath.Join(dir, "redaction-policy.json"),
	}
	writeTestJSON(t, filepath.Join(runDir, "run.json"), input.Run)
	writeTestJSON(t, paths.policy, input.Policy)
	return paths
}
