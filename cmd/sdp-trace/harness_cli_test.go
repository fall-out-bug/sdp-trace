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

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

func TestHarnessObserveValidateSummarizeCLI(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessCLIEvents(t, dir, []map[string]any{
		harnessCLIEvent("e1", "harness"),
		harnessCLIEvent("e2", "model"),
	})
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"harness", "observe", "--profile", "profile.json", "--source", "events.jsonl", "--out", "run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("observe exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "run", "--out", "validation.json"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"validation_state": "pass"`) {
		t.Fatalf("validate stdout missing pass: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "summarize", "--validation", "validation.json"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("summarize exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "no harness compliance") {
		t.Fatalf("summary missing boundary: %s", out.String())
	}
}

func TestHarnessObserveCLIRejectsUnsafePrompt(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	event := harnessCLIEvent("e1", "harness")
	event["raw_prompt"] = "secret prompt"
	writeHarnessCLIEvents(t, dir, []map[string]any{event})
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"harness", "observe", "--profile", "profile.json", "--source", "events.jsonl", "--out", "run"}, &out, &errOut)
	if exit == 0 {
		t.Fatalf("observe exit=0, want non-zero")
	}
	if !strings.Contains(errOut.String(), "unsafe_input") {
		t.Fatalf("stderr missing unsafe_input: %s", errOut.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "run")); !os.IsNotExist(err) {
		t.Fatalf("run dir exists after unsafe observe: %v", err)
	}
}

func TestHarnessCLIRequiresDocumentedFlags(t *testing.T) {
	var out, errOut bytes.Buffer
	exit := run([]string{"harness", "observe"}, &out, &errOut)
	if exit != exitUsage || !strings.Contains(errOut.String(), "requires --profile") {
		t.Fatalf("observe missing flags exit=%d stderr=%s", exit, errOut.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate"}, &out, &errOut)
	if exit != exitUsage || !strings.Contains(errOut.String(), "requires --profile") {
		t.Fatalf("validate missing flags exit=%d stderr=%s", exit, errOut.String())
	}
}

func writeHarnessCLIProfile(t *testing.T, dir string) {
	t.Helper()
	profile := harnessobs.Profile{
		SchemaVersion:         harnessobs.ProfileSchemaVersion,
		ProfileID:             "generic-harness-v1",
		HarnessFamily:         "generic-harness",
		EventSchemaVersion:    harnessobs.EventSchemaVersion,
		RequiredEventFamilies: []string{"harness", "model"},
		RawRetentionPolicy:    "digest_only",
		DegradationRules: map[string]harnessobs.Rule{
			"missing_required_family": {State: harnessobs.StateNotAssessed, ReasonCode: "required_event_family_absent"},
			"missing_optional_family": {State: harnessobs.StateNotAssessed, ReasonCode: "optional_event_family_absent"},
			"source_unavailable":      {State: harnessobs.StateCannotVerify, ReasonCode: "source_unavailable"},
			"unsafe_input":            {State: harnessobs.StateFail, ReasonCode: "unsafe_input"},
			"digest_mismatch":         {State: harnessobs.StateCannotVerify, ReasonCode: "source_digest_mismatch"},
			"schema_version_mismatch": {State: harnessobs.StateCannotVerify, ReasonCode: "schema_version_mismatch"},
			"cross_link_conflict":     {State: harnessobs.StateCannotVerify, ReasonCode: "adapter_harness_state_conflict"},
		},
	}
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "profile.json"), append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func harnessCLIEvent(id, family string) map[string]any {
	return map[string]any{
		"event_id":             id,
		"event_schema_version": harnessobs.EventSchemaVersion,
		"event_family":         family,
		"event_type":           family + "_observed",
		"observed_at":          "2026-05-09T12:00:00Z",
		"source_ref":           "src-" + id,
		"source_digest":        "",
		"content_state":        harnessobs.ContentDigestOnly,
	}
}

func writeHarnessCLIEvents(t *testing.T, dir string, events []map[string]any) {
	t.Helper()
	lines := make([]string, 0, len(events))
	for _, event := range events {
		event["source_digest"] = ""
		line := marshalHarnessCLI(t, event)
		event["source_digest"] = digestHarnessCLI(t, line)
		lines = append(lines, marshalHarnessCLI(t, event))
	}
	if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
}

func marshalHarnessCLI(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func digestHarnessCLI(t *testing.T, line string) string {
	t.Helper()
	var raw map[string]any
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		t.Fatal(err)
	}
	raw["source_digest"] = ""
	canonical, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:])
}

func chdirCLI(t *testing.T, dir string) func() {
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
