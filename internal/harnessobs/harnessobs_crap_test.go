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
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
