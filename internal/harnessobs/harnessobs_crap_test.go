package harnessobs

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSessionCollectionRejectsMismatchedProfileAndMissingOptions(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
	})
	runDir := filepath.Join(dir, "session-run")
	if err := os.Mkdir(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJSONFixture(t, filepath.Join(runDir, "session.json"), SessionRun{
		SchemaVersion: SessionRunSchemaVersion,
		ProfileID:     "other-profile",
		CreatedAt:     "2026-05-10T12:00:00Z",
	})

	oldwd := chdir(t, dir)
	defer oldwd()
	for name, opts := range map[string]SessionCollectOptions{
		"missing profile": {RunDir: "session-run"},
		"missing run":     {ProfilePath: "session-profile.json"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, err := validateSessionCollectOptions(opts); err == nil {
				t.Fatalf("validateSessionCollectOptions() error = nil")
			}
		})
	}
	if _, err := prepareSessionCollection(SessionCollectOptions{ProfilePath: "session-profile.json", RunDir: "session-run"}); err == nil || !strings.Contains(err.Error(), "session profile mismatch") {
		t.Fatalf("prepareSessionCollection() error = %v, want session profile mismatch", err)
	}
}

func TestRunSessionCollectsCommandGeneratedEvents(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeEvents(t, dir, []map[string]any{eventMap("e1", "harness")})
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
	})
	oldwd := chdir(t, dir)
	defer oldwd()

	session, observed, err := RunSession(SessionOptions{
		ProfilePath: "session-profile.json",
		OutDir:      "session-run",
		Command:     []string{"sh", "-c", ":"},
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}
	if session.CollectionState != StatePass || observed.EventCount != 1 {
		t.Fatalf("session/observed = %+v / %+v", session, observed)
	}
	if session.ProcessIDState != StatePass || session.CommandDigestState != StatePass {
		t.Fatalf("process/command states = %s/%s", session.ProcessIDState, session.CommandDigestState)
	}
}

func TestCollectFinishedSessionReturnsCollectedRunWithWaitError(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeEvents(t, dir, []map[string]any{eventMap("e1", "harness")})
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
	})
	if err := os.Mkdir(filepath.Join(dir, "session-run"), 0o755); err != nil {
		t.Fatal(err)
	}
	session := SessionRun{
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
		CreatedAt:          "2026-05-10T12:00:00Z",
	}

	oldwd := chdir(t, dir)
	defer oldwd()
	waitErr := errors.New("command failed")
	collected, observed, err := collectFinishedSession(SessionOptions{ProfilePath: "session-profile.json", OutDir: "session-run"}, session, waitErr, time.Date(2026, 5, 10, 12, 30, 0, 0, time.UTC))
	if !errors.Is(err, waitErr) {
		t.Fatalf("collectFinishedSession() error = %v, want wait error", err)
	}
	if collected.CollectionState != StatePass || observed.EventCount != 1 {
		t.Fatalf("collected/observed = %+v / %+v", collected, observed)
	}
}

func TestRawNormalizationAndSignalUtilityBranches(t *testing.T) {
	dir := t.TempDir()
	rawPath := filepath.Join(dir, "raw.jsonl")
	outPath := filepath.Join(dir, "nested", "events.jsonl")
	if err := os.WriteFile(rawPath, []byte("\n{\"type\":\"tool_use\",\"model\":\"qwen\",\"tool\":\"edit\"}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	facts := []Event{normalizedEvent("fact", "model", "model_observed", "2026-05-10T12:00:00Z", "session", "qwen")}
	if err := normalizeRawEvents(OpenCodeJSONLRawFormat, rawPath, outPath, facts, time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("normalizeRawEvents() error = %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(strings.TrimSpace(string(data)), "\n") + 1; got != 4 {
		t.Fatalf("normalized line count = %d, want 4\n%s", got, string(data))
	}
	if err := normalizeRawEvents("unknown", rawPath, filepath.Join(dir, "bad.jsonl"), nil, time.Time{}); err == nil || !strings.Contains(err.Error(), "unsupported raw_event_format") {
		t.Fatalf("normalizeRawEvents() unsupported error = %v", err)
	}
	if err := normalizeRawEvents(OpenCodeJSONLRawFormat, rawPath, rawPath, nil, time.Time{}); err == nil || !strings.Contains(err.Error(), "must be different files") {
		t.Fatalf("normalizeRawEvents() same path error = %v", err)
	}

	signals := rawSignalsAt("items", []any{map[string]any{"model": "qwen"}, "literal"})
	if !containsString(signals, "model") || !containsString(signals, "qwen") {
		t.Fatalf("rawSignalsAt() = %+v, want nested map signals", signals)
	}
}

func TestSmallUtilityEdgeBranches(t *testing.T) {
	if got := unixMillisTimestamp(1700000000123); got != "2023-11-14T22:13:20Z" {
		t.Fatalf("unixMillisTimestamp() = %s", got)
	}
	if got := unixMillisTimestamp(1700000000); got != "2023-11-14T22:13:20Z" {
		t.Fatalf("unixMillisTimestamp(seconds) = %s", got)
	}
	if got := unixMillisTimestamp(-1); got != "" {
		t.Fatalf("unixMillisTimestamp(-1) = %q, want empty", got)
	}
	if got := safeToken("://"); got != "opencode" {
		t.Fatalf("safeToken() = %s, want opencode", got)
	}
	if got := safeCommandModel(strings.Repeat("a", 160)); got != "" {
		t.Fatalf("safeCommandModel() = %q, want empty for overlong model", got)
	}
	if safeEventRef("events/../e1.json") {
		t.Fatalf("safeEventRef accepted traversal")
	}
	if rank("unexpected") != 0 {
		t.Fatalf("rank(unexpected) = %d, want 0", rank("unexpected"))
	}
	if reason := unsafeStringReason("path", "/Users/name/.ssh/id_rsa", false); reason != "unsafe_path_or_private_path" {
		t.Fatalf("unsafeStringReason() = %s, want unsafe_path_or_private_path", reason)
	}
	if reason := unsafeStringReason("source_digest", strings.Repeat("a", 64), false); reason != "" {
		t.Fatalf("unsafeStringReason() digest reason = %s, want empty", reason)
	}
	if !hasKeyInSlice([]any{map[string]any{"model": "qwen"}}, map[string]bool{"model": true}) {
		t.Fatalf("hasKeyInSlice() = false, want true")
	}
	decoded, err := DecodeValidation(bytes.NewBufferString(`{"schema_version":"harness-observation-validation-v1","validation_state":"pass"}`))
	if err != nil {
		t.Fatalf("DecodeValidation() error = %v", err)
	}
	if decoded.SchemaVersion != ValidationSchemaVersion {
		t.Fatalf("DecodeValidation() schema = %s", decoded.SchemaVersion)
	}
	if _, err := DecodeValidation(bytes.NewBufferString("{")); err == nil {
		t.Fatalf("DecodeValidation(invalid) error = nil")
	}
}

func TestSetupActionAndDegradationUtilityBranches(t *testing.T) {
	for _, kind := range []string{"init", "profile", "wrapper", "hook", "context_isolation"} {
		if err := validateSessionSetupAction(SessionSetupAction{ID: "setup-1", Kind: kind}); err != nil {
			t.Fatalf("validateSessionSetupAction(%s) error = %v", kind, err)
		}
	}
	if err := validateSessionSetupActions([]SessionSetupAction{
		{ID: "setup-1", Kind: "init"},
		{ID: "setup-2", Kind: "profile"},
		{ID: "setup-3", Kind: "wrapper"},
		{ID: "setup-4", Kind: "hook"},
	}); err == nil || !strings.Contains(err.Error(), "too many setup actions") {
		t.Fatalf("validateSessionSetupActions(too many) error = %v", err)
	}
	if err := validateSessionSetupAction(SessionSetupAction{ID: "../bad", Kind: "hook"}); err == nil || !strings.Contains(err.Error(), "unsafe setup action id") {
		t.Fatalf("validateSessionSetupAction(bad id) error = %v", err)
	}
	if err := validateSessionSetupAction(SessionSetupAction{ID: "setup-1", Kind: "shell"}); err == nil || !strings.Contains(err.Error(), "unsupported setup action kind") {
		t.Fatalf("validateSessionSetupAction(bad kind) error = %v", err)
	}
	if !validDegradationRule(Rule{State: StateCannotVerify, ReasonCode: "missing-proof"}) {
		t.Fatalf("validDegradationRule(valid) = false")
	}
	if validDegradationRule(Rule{State: "unknown", ReasonCode: "missing-proof"}) {
		t.Fatalf("validDegradationRule(unknown state) = true")
	}
	if validDegradationRule(Rule{State: StateCannotVerify, ReasonCode: "../bad"}) {
		t.Fatalf("validDegradationRule(unsafe reason) = true")
	}
}

func TestCommandModelSafetyAndSourceDigest(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "source.txt")
	if err := os.WriteFile(sourcePath, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := safeCommandModel(" qwen3-coder "); got != "qwen3-coder" {
		t.Fatalf("safeCommandModel(trim) = %q, want qwen3-coder", got)
	}
	for _, model := range []string{
		"https://example.invalid/model",
		"qwen coder",
		`qwen"coder`,
		`qwen\coder`,
		"../qwen",
		"/tmp/qwen",
		strings.Repeat("a", 129),
	} {
		if got := safeCommandModel(model); got != "" {
			t.Fatalf("safeCommandModel(%q) = %q, want empty", model, got)
		}
	}

	if got := digestFile(filepath.Join(dir, "missing.txt")); got != "" {
		t.Fatalf("digestFile(missing) = %q, want empty", got)
	}
	if got := digestFile(sourcePath); got != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("digestFile(abc) = %q, want SHA-256 of abc", got)
	}
}

func TestNormalizedWriteAndShellAndSourceCommitBranches(t *testing.T) {
	dir := t.TempDir()
	blockingFile := filepath.Join(dir, "file")
	if err := os.WriteFile(blockingFile, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeNormalizedEvents(filepath.Join(blockingFile, "events.jsonl"), nil); err == nil {
		t.Fatalf("writeNormalizedEvents(blocked parent) error = nil")
	}
	if got := shellFields(`opencode run --model "qwen coder"`); len(got) != 4 || got[3] != "qwen coder" {
		t.Fatalf("shellFields(quoted model) = %#v", got)
	}

	oldwd := chdir(t, dir)
	if got := sourceCommit(); got != "" {
		t.Fatalf("sourceCommit(non-git) = %q, want empty", got)
	}
	oldwd()
	if got := sourceCommit(); got != "" && len(got) != 40 {
		t.Fatalf("sourceCommit(repo) = %q, want empty or full hash", got)
	}
}

func TestNormalizeSessionStreamCaptureErrors(t *testing.T) {
	for name, tc := range map[string]struct {
		mode    string
		wantErr string
	}{
		"digest only": {
			mode:    ContentDigestOnly,
			wantErr: "stream_capture mode not implemented",
		},
		"unsupported": {
			mode:    "full",
			wantErr: "unsupported stream_capture",
		},
	} {
		t.Run(name, func(t *testing.T) {
			profile := SessionProfile{StreamCapture: tc.mode}
			err := normalizeSessionStreamCapture(&profile)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("normalizeSessionStreamCapture() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestIsolationRuleValidationBranches(t *testing.T) {
	valid := SessionIsolationRule{
		ID:         "rule-1",
		Kind:       "json_read_deny",
		TargetPath: "settings.json",
		Pattern:    "/tmp/private",
	}
	if err := validateSessionIsolationRule(valid); err != nil {
		t.Fatalf("validateSessionIsolationRule(valid) error = %v", err)
	}
	if err := validateSessionIsolationRules([]SessionIsolationRule{valid, {
		ID:         "bad-rule",
		Kind:       "json_read_deny",
		TargetPath: "settings.json",
		Pattern:    " \t",
	}}); err == nil || !strings.Contains(err.Error(), "unsafe isolation rule pattern") {
		t.Fatalf("validateSessionIsolationRules() error = %v, want unsafe isolation rule pattern", err)
	}

	for name, tc := range map[string]struct {
		rule    SessionIsolationRule
		wantErr string
	}{
		"bad id": {
			rule:    SessionIsolationRule{ID: "../bad", Kind: "json_read_deny", TargetPath: "settings.json", Pattern: "/tmp/private"},
			wantErr: "unsafe isolation rule id",
		},
		"blank pattern": {
			rule:    SessionIsolationRule{ID: "rule-1", Kind: "json_read_deny", TargetPath: "settings.json", Pattern: " \t"},
			wantErr: "unsafe isolation rule pattern",
		},
		"newline pattern": {
			rule:    SessionIsolationRule{ID: "rule-1", Kind: "json_read_deny", TargetPath: "settings.json", Pattern: "a\nb"},
			wantErr: "unsafe isolation rule pattern",
		},
		"unsafe target": {
			rule:    SessionIsolationRule{ID: "rule-1", Kind: "json_read_deny", TargetPath: "../settings.json", Pattern: "/tmp/private"},
			wantErr: "unsafe isolation target path",
		},
		"unsupported kind": {
			rule:    SessionIsolationRule{ID: "rule-1", Kind: "write_allow", TargetPath: "settings.json", Pattern: "/tmp/private"},
			wantErr: "unsupported isolation rule kind",
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := validateSessionIsolationRule(tc.rule)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("validateSessionIsolationRule() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestSafeProfileRelativeIsolationFileBranches(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()

	path, err := safeProfileRelativeIsolationFile("demo/session-profile.json", "settings/permissions.json")
	if err != nil {
		t.Fatalf("safeProfileRelativeIsolationFile() error = %v", err)
	}
	if path != filepath.Join("demo", "settings", "permissions.json") {
		t.Fatalf("safeProfileRelativeIsolationFile() = %q", path)
	}

	for name, relPath := range map[string]string{
		"absolute":   filepath.Join(string(filepath.Separator), "tmp", "settings.json"),
		"traversal":  "../settings.json",
		"url":        "file://settings.json",
		"blank base": " ",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := safeProfileRelativeIsolationFile("session-profile.json", relPath); err == nil {
				t.Fatalf("safeProfileRelativeIsolationFile(%q) error = nil", relPath)
			}
		})
	}

	outside := t.TempDir()
	if err := os.Symlink(outside, "outside-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := safeProfileRelativeIsolationFile("session-profile.json", "outside-link/nested/settings.json"); err == nil || !strings.Contains(err.Error(), "escapes working directory") {
		t.Fatalf("safeProfileRelativeIsolationFile(symlink escape) error = %v, want escapes working directory", err)
	}
}

func TestEnsureJSONReadDenyRuleAndPresenceBranches(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()
	path := "settings.json"

	if err := ensureJSONReadDenyRule(path, "/tmp/private"); err != nil {
		t.Fatalf("ensureJSONReadDenyRule(missing) error = %v", err)
	}
	if present, err := jsonReadDenyRulePresent(path, "/tmp/private"); err != nil || !present {
		t.Fatalf("jsonReadDenyRulePresent() = %v/%v, want true/nil", present, err)
	}
	if present, err := jsonReadDenyRulePresent(path, "/tmp/other"); err != nil || present {
		t.Fatalf("jsonReadDenyRulePresent(absent) = %v/%v, want false/nil", present, err)
	}

	writeJSONFixture(t, path, map[string]any{
		"permission": "replace-me",
		"kept":       true,
	})
	if err := ensureJSONReadDenyRule(path, "/home/private"); err != nil {
		t.Fatalf("ensureJSONReadDenyRule(replace permission) error = %v", err)
	}
	var config map[string]any
	readJSONFixture(t, path, &config)
	if config["kept"] != true {
		t.Fatalf("kept field = %v, want true", config["kept"])
	}
	permission, ok := config["permission"].(map[string]any)
	if !ok {
		t.Fatalf("permission = %#v, want object", config["permission"])
	}
	read, ok := permission["read"].(map[string]any)
	if !ok || read["/home/private"] != "deny" {
		t.Fatalf("permission.read = %#v, want /home/private deny", permission["read"])
	}

	if err := os.WriteFile(path, []byte(" \n\t"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ensureJSONReadDenyRule(path, "/blank"); err != nil {
		t.Fatalf("ensureJSONReadDenyRule(blank JSON) error = %v", err)
	}

	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ensureJSONReadDenyRule(path, "/bad"); err == nil {
		t.Fatalf("ensureJSONReadDenyRule(invalid JSON) error = nil")
	}
}

func TestIsolationRulePresentBranches(t *testing.T) {
	dir := t.TempDir()
	linePath := filepath.Join(dir, "ignore")
	if err := os.WriteFile(linePath, []byte(".evidence/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	present, err := isolationRulePresent(SessionIsolationRule{Kind: "ignore_line", TargetPath: linePath, Pattern: ".evidence/"})
	if err != nil || !present {
		t.Fatalf("isolationRulePresent(ignore_line) = %v/%v, want true/nil", present, err)
	}
	present, err = isolationRulePresent(SessionIsolationRule{Kind: "ignore_line", TargetPath: linePath, Pattern: ".trace/"})
	if err != nil || present {
		t.Fatalf("isolationRulePresent(ignore_line absent) = %v/%v, want false/nil", present, err)
	}
	unreadableLinePath := filepath.Join(dir, "unreadable")
	if err := os.Mkdir(unreadableLinePath, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := lineIsolationRulePresent(unreadableLinePath, ".evidence/"); err == nil {
		t.Fatalf("lineIsolationRulePresent(unreadable directory) error = nil")
	}
	if _, err := isolationRulePresent(SessionIsolationRule{Kind: "unknown"}); err == nil || !strings.Contains(err.Error(), "unsupported isolation rule kind") {
		t.Fatalf("isolationRulePresent(unknown) error = %v, want unsupported kind", err)
	}
	if err := ensureIsolationRule(SessionIsolationRule{Kind: "unknown"}); err == nil || !strings.Contains(err.Error(), "unsupported isolation rule kind") {
		t.Fatalf("ensureIsolationRule(unknown) error = %v, want unsupported kind", err)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
