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

func TestObserveSetupCollectSupportsSetUpAndForgetWorkflow(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessSessionProfile(t, dir, "events.jsonl")
	writeHarnessCLIEventsFile(t, filepath.Join(dir, "source-events.jsonl"), []map[string]any{
		harnessCLIEvent("e1", "harness"),
		harnessCLIEvent("e2", "model"),
	})
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run", "--command", "opencode run demo"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"command_digest_state": "pass"`) {
		t.Fatalf("setup output missing command digest state: %s", out.String())
	}

	if data, err := os.ReadFile("source-events.jsonl"); err != nil {
		t.Fatal(err)
	} else if err := os.WriteFile("events.jsonl", data, 0o644); err != nil {
		t.Fatal(err)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("collect exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"observed_run_dir": "observed"`) {
		t.Fatalf("collect output missing observed run dir: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed", "--out", "validation.json"}, &out, &errOut)
	if exit != 0 || !strings.Contains(out.String(), `"validation_state": "pass"`) {
		t.Fatalf("validate exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
}

func TestObserveSessionRunsControlledProxyWithoutRetainingStdout(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessSessionProfile(t, dir, "events.jsonl")
	writeHarnessCLIEventsFile(t, filepath.Join(dir, "source-events.jsonl"), []map[string]any{
		harnessCLIEvent("e1", "harness"),
		harnessCLIEvent("e2", "model"),
	})
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{
		"observe", "session", "--profile", "session-profile.json", "--out", "session-run", "--",
		"sh", "-c", "cp source-events.jsonl events.jsonl && printf 'sk-should-not-be-captured'",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("session exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if strings.Contains(out.String(), "sk-should-not-be-captured") {
		t.Fatalf("session leaked command stdout: %s", out.String())
	}
	if strings.Contains(errOut.String(), "sk-should-not-be-captured") {
		t.Fatalf("session leaked command stderr: %s", errOut.String())
	}
	if !strings.Contains(out.String(), `"command_digest"`) ||
		!strings.Contains(out.String(), `"command_digest_state": "pass"`) ||
		!strings.Contains(out.String(), `"observed_run_dir": "observed"`) {
		t.Fatalf("session output missing provenance: %s", out.String())
	}
}

func TestObserveSessionNormalizesOpenCodeRawJSONL(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "model", "interaction", "phase", "tool", "mutation", "test"})
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := strings.Join([]string{
		`{"type":"session.started","provider":"minimax","model":"minimax-coding-plan/MiniMax-M2.5","timestamp":"2026-05-09T12:00:00Z"}`,
		`{"type":"message","role":"assistant","content":"ack"}`,
		`{"type":"phase","name":"gsd.plan"}`,
		`{"type":"tool.call","tool":"edit"}`,
		`{"type":"file.write","path":"src/App.kt"}`,
		`{"type":"test.finished","status":"pass"}`,
	}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "raw-source.jsonl"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{
		"observe", "session", "--profile", "session-profile.json", "--out", "session-run", "--",
		"sh", "-c", "cp raw-source.jsonl opencode-raw.jsonl && printf 'raw output must not be retained'",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("session exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if strings.Contains(out.String(), "raw output must not be retained") || strings.Contains(errOut.String(), "raw output must not be retained") {
		t.Fatalf("session retained command output stdout=%s stderr=%s", out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), `"normalized_digest"`) || !strings.Contains(out.String(), `"collection_state": "pass"`) {
		t.Fatalf("session output missing raw normalization evidence: %s", out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "normalized-events.jsonl")); err != nil {
		t.Fatalf("normalized source not written: %v", err)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed", "--out", "validation.json"}, &out, &errOut)
	if exit != 0 || !strings.Contains(out.String(), `"validation_state": "pass"`) {
		t.Fatalf("validate exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
}

func TestObserveSessionNormalizesNativeOpenCodeJSONL(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "model", "interaction", "phase"})
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := strings.Join([]string{
		`{"type":"step_start","timestamp":1778326453680,"sessionID":"ses_fixture","part":{"id":"prt_start","messageID":"msg_fixture","sessionID":"ses_fixture","snapshot":"4b825dc642cb6eb9a060e54bf8d69288fbee4904","type":"step-start"}}`,
		`{"type":"text","timestamp":1778326453882,"sessionID":"ses_fixture","part":{"id":"prt_text","messageID":"msg_fixture","sessionID":"ses_fixture","type":"text","text":"OK","time":{"start":1778326453875,"end":1778326453882}}}`,
		`{"type":"step_finish","timestamp":1778326453906,"sessionID":"ses_fixture","part":{"id":"prt_finish","reason":"stop","snapshot":"4b825dc642cb6eb9a060e54bf8d69288fbee4904","messageID":"msg_fixture","sessionID":"ses_fixture","type":"step-finish","tokens":{"total":13137,"input":6,"output":27,"reasoning":0},"cost":0}}`,
	}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "raw-source.jsonl"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{
		"observe", "session", "--profile", "session-profile.json", "--out", "session-run", "--",
		"sh", "-c", "cp raw-source.jsonl opencode-raw.jsonl", "--model", "minimax-coding-plan/MiniMax-M2.5",
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("session exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"event_count": 4`) || !strings.Contains(out.String(), `"collection_state": "pass"`) {
		t.Fatalf("session output missing native OpenCode events: %s", out.String())
	}
	if !strings.Contains(out.String(), `"command_model": "minimax-coding-plan/MiniMax-M2.5"`) {
		t.Fatalf("session output missing command model fact: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed"}, &out, &errOut)
	if exit == 0 || !strings.Contains(out.String(), `"validation_state": "not_assessed"`) {
		t.Fatalf("validate exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	for _, family := range []string{"harness", "model", "interaction"} {
		if !strings.Contains(out.String(), `"family": "`+family+`"`) || !strings.Contains(out.String(), `"state": "pass"`) {
			t.Fatalf("validation missing pass family %s: %s", family, out.String())
		}
	}
	if !strings.Contains(out.String(), `"family": "phase"`) || !strings.Contains(out.String(), `"state": "not_assessed"`) {
		t.Fatalf("phase should remain not_assessed: %s", out.String())
	}
}

func TestObserveCollectNormalizesNativeOpenCodeToolUseWithPrivateOutput(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "tool", "mutation"})
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := strings.Join([]string{
		`{"type":"step_start","timestamp":1778326472992,"sessionID":"ses_fixture","part":{"type":"step-start"}}`,
		`{"type":"tool_use","timestamp":1778326473358,"sessionID":"ses_fixture","part":{"type":"tool","tool":"glob","callID":"call_fixture","state":{"status":"completed","input":{"pattern":"*"},"output":"/Users/fall_out_bug/projects/vibe_coding/sdp-trace/schema/harness-event.schema.json"}}}`,
	}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("collect exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if data, err := os.ReadFile(filepath.Join(dir, "normalized-events.jsonl")); err != nil {
		t.Fatal(err)
	} else if strings.Contains(string(data), "/Users/fall_out_bug") {
		t.Fatalf("normalized output retained private raw body: %s", string(data))
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed"}, &out, &errOut)
	if exit == 0 || !strings.Contains(out.String(), `"validation_state": "not_assessed"`) {
		t.Fatalf("validate exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"family": "tool"`) || !strings.Contains(out.String(), `"state": "pass"`) ||
		!strings.Contains(out.String(), `"family": "mutation"`) || !strings.Contains(out.String(), `"state": "not_assessed"`) {
		t.Fatalf("tool should pass while mutation remains not_assessed: %s", out.String())
	}
}

func TestObserveCollectTreatsNativeEditToolAsMutation(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "tool", "mutation"})
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := strings.Join([]string{
		`{"type":"step_start","timestamp":1778326472992,"sessionID":"ses_fixture","part":{"type":"step-start"}}`,
		`{"type":"tool_use","timestamp":1778326473358,"sessionID":"ses_fixture","part":{"type":"tool","tool":"edit","callID":"call_fixture","state":{"status":"completed","input":{"file":"src/App.kt"},"output":"updated"}}}`,
	}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("collect exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed"}, &out, &errOut)
	if exit != 0 || !strings.Contains(out.String(), `"validation_state": "pass"`) {
		t.Fatalf("validate exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	for _, family := range []string{"tool", "mutation"} {
		if !strings.Contains(out.String(), `"family": "`+family+`"`) || !strings.Contains(out.String(), `"state": "pass"`) {
			t.Fatalf("validation missing pass family %s: %s", family, out.String())
		}
	}
}

func TestObserveSessionDoesNotPromoteMessageTextToEvidence(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "model", "interaction", "phase", "tool", "mutation", "test"})
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := strings.Join([]string{
		`{"type":"session.started","provider":"minimax","model":"minimax-coding-plan/MiniMax-M2.5"}`,
		`{"type":"message","role":"assistant","content":"I used a tool, edited files, completed the phase, and tests pass"}`,
	}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("collect exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed", "--out", "validation.json"}, &out, &errOut)
	if exit == 0 || !strings.Contains(out.String(), `"validation_state": "not_assessed"`) {
		t.Fatalf("validate should not promote message text exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	for _, family := range []string{"phase", "tool", "mutation", "test"} {
		if !strings.Contains(out.String(), `"family": "`+family+`"`) || !strings.Contains(out.String(), `"state": "not_assessed"`) {
			t.Fatalf("validation missing not_assessed family %s: %s", family, out.String())
		}
	}
}

func TestObserveCollectDoesNotTreatFileReadAsMutation(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "model", "mutation"})
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := strings.Join([]string{
		`{"type":"session.started","provider":"minimax","model":"minimax-coding-plan/MiniMax-M2.5"}`,
		`{"type":"file.read","path":"src/App.kt"}`,
	}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("collect exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed"}, &out, &errOut)
	if exit == 0 || !strings.Contains(out.String(), `"family": "mutation"`) || !strings.Contains(out.String(), `"state": "not_assessed"`) {
		t.Fatalf("file.read should not count as mutation exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
}

func TestObserveCollectDoesNotFabricateEventsForUnrecognizedRawJSONL(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "model"})
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(`{"type":"custom.event","data":"x"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != 0 || !strings.Contains(out.String(), `"collection_state": "pass"`) {
		t.Fatalf("collect exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if data, err := os.ReadFile(filepath.Join(dir, "normalized-events.jsonl")); err != nil {
		t.Fatal(err)
	} else if strings.TrimSpace(string(data)) != "" {
		t.Fatalf("unrecognized raw source produced events: %s", string(data))
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"harness", "validate", "--profile", "profile.json", "--run", "session-run/observed"}, &out, &errOut)
	if exit == 0 || !strings.Contains(out.String(), `"validation_state": "not_assessed"`) {
		t.Fatalf("validate should remain not_assessed exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
}

func TestObserveCollectRejectsUnsafeRawOpenCodeJSONL(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := `{"type":"session.started","provider":"minimax","model":"minimax-coding-plan/MiniMax-M2.5","api_key":"redacted"}`
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(raw+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit == 0 || !strings.Contains(errOut.String(), "unsafe_input:api_key:sensitive_field") {
		t.Fatalf("collect should reject unsafe raw source exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "normalized-events.jsonl")); err == nil {
		t.Fatalf("normalized source written after unsafe raw input")
	}
}

func TestObserveCollectRejectsNestedUnsafeRawOpenCodeJSONL(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	raw := `{"type":"session.started","provider":"minimax","model":"minimax-coding-plan/MiniMax-M2.5","content":{"api_key":"redacted"}}`
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(raw+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit == 0 || !strings.Contains(errOut.String(), "unsafe_input:content.api_key:sensitive_field") {
		t.Fatalf("collect should reject nested unsafe raw source exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
}

func TestObserveCollectRejectsUnsafeRawEventSourcePath(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessSessionProfileWithRaw(t, dir, "normalized-events.jsonl", "../opencode-raw.jsonl", harnessobs.OpenCodeJSONLRawFormat)
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit == 0 || !strings.Contains(errOut.String(), "raw_event_source_path invalid") {
		t.Fatalf("collect should reject unsafe raw source path exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
}

func TestObserveCollectRecordsCannotVerifyForMissingSource(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessSessionProfile(t, dir, "missing-events.jsonl")
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("collect exit=%d, want cannot_verify; stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"collection_state": "cannot_verify"`) || !strings.Contains(out.String(), `"collection_reason": "source_unavailable"`) {
		t.Fatalf("collect output missing cannot_verify evidence: %s", out.String())
	}
}

func TestObserveCollectResolvesSourcesRelativeToSessionProfile(t *testing.T) {
	dir := t.TempDir()
	profileDir := filepath.Join(dir, "profiles")
	if err := os.Mkdir(profileDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeHarnessCLIProfile(t, profileDir)
	writeHarnessSessionProfile(t, profileDir, "events.jsonl")
	writeHarnessCLIEventsFile(t, filepath.Join(profileDir, "events.jsonl"), []map[string]any{
		harnessCLIEvent("e1", "harness"),
		harnessCLIEvent("e2", "model"),
	})
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "profiles/session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"observe", "collect", "--profile", "profiles/session-profile.json", "--run", "session-run"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("collect exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"event_count": 2`) {
		t.Fatalf("collect output missing observed event count: %s", out.String())
	}
}

func TestObserveSetupRejectsUnimplementedStreamCapture(t *testing.T) {
	dir := t.TempDir()
	writeHarnessCLIProfile(t, dir)
	writeHarnessSessionProfileWithStream(t, dir, "events.jsonl", harnessobs.ContentDigestOnly)
	oldwd := chdirCLI(t, dir)
	defer oldwd()

	var out, errOut bytes.Buffer
	exit := run([]string{"observe", "setup", "--profile", "session-profile.json", "--out", "session-run"}, &out, &errOut)
	if exit == 0 || !strings.Contains(errOut.String(), "stream_capture mode not implemented") {
		t.Fatalf("setup exit=%d stderr=%s stdout=%s", exit, errOut.String(), out.String())
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
	writeHarnessCLIProfileWithFamilies(t, dir, []string{"harness", "model"})
}

func writeHarnessCLIProfileWithFamilies(t *testing.T, dir string, families []string) {
	t.Helper()
	profile := harnessobs.Profile{
		SchemaVersion:         harnessobs.ProfileSchemaVersion,
		ProfileID:             "generic-harness-v1",
		HarnessFamily:         "generic-harness",
		EventSchemaVersion:    harnessobs.EventSchemaVersion,
		RequiredEventFamilies: families,
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

func writeHarnessSessionProfile(t *testing.T, dir, eventSource string) {
	t.Helper()
	writeHarnessSessionProfileWithStream(t, dir, eventSource, "disabled")
}

func writeHarnessSessionProfileWithStream(t *testing.T, dir, eventSource, streamCapture string) {
	t.Helper()
	profile := harnessobs.SessionProfile{
		SchemaVersion:      harnessobs.SessionProfileSchemaVersion,
		ProfileID:          "opencode-gsd-fixture-v1",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    eventSource,
		SetupActions: []harnessobs.SessionSetupAction{
			{ID: "init", Kind: "init", Required: true},
			{ID: "profile", Kind: "profile", Required: true},
		},
		StreamCapture: streamCapture,
	}
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "session-profile.json"), append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeHarnessSessionProfileWithRaw(t *testing.T, dir, eventSource, rawSource, rawFormat string) {
	t.Helper()
	profile := harnessobs.SessionProfile{
		SchemaVersion:      harnessobs.SessionProfileSchemaVersion,
		ProfileID:          "opencode-gsd-fixture-v1",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    eventSource,
		RawEventSourcePath: rawSource,
		RawEventFormat:     rawFormat,
		SetupActions: []harnessobs.SessionSetupAction{
			{ID: "init", Kind: "init", Required: true},
			{ID: "profile", Kind: "profile", Required: true},
		},
		StreamCapture: "disabled",
	}
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "session-profile.json"), append(data, '\n'), 0o644); err != nil {
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
	writeHarnessCLIEventsFile(t, filepath.Join(dir, "events.jsonl"), events)
}

func writeHarnessCLIEventsFile(t *testing.T, path string, events []map[string]any) {
	t.Helper()
	lines := make([]string, 0, len(events))
	for _, event := range events {
		event["source_digest"] = ""
		line := marshalHarnessCLI(t, event)
		event["source_digest"] = digestHarnessCLI(t, line)
		lines = append(lines, marshalHarnessCLI(t, event))
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
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
