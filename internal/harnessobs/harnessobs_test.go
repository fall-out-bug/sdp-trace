package harnessobs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestObserveValidateCompleteHarnessExport(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness", "model", "phase"}, []string{"review"})
	writeEvents(t, dir, []map[string]any{
		eventMap("e1", "harness"),
		eventMap("e2", "model"),
		eventMap("e3", "phase"),
	})

	oldwd := chdir(t, dir)
	defer oldwd()

	run, err := Observe(ObserveOptions{
		ProfilePath: "profile.json",
		SourcePath:  "events.jsonl",
		OutDir:      "run",
		Now:         time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if run.EventCount != 3 {
		t.Fatalf("EventCount = %d, want 3", run.EventCount)
	}
	validation, err := Validate(ValidateOptions{ProfilePath: "profile.json", RunDir: "run"})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if validation.ValidationState != StatePass {
		t.Fatalf("ValidationState = %s, want pass: %+v", validation.ValidationState, validation)
	}
}

func TestValidateZeroEventSourceIsNotAssessed(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness", "model"}, nil)
	if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdir(t, dir)
	defer oldwd()
	if _, err := Observe(ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl", OutDir: "run"}); err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	validation, err := Validate(ValidateOptions{ProfilePath: "profile.json", RunDir: "run"})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if validation.ValidationState != StateNotAssessed || validation.ReasonCode != "required_event_family_absent" {
		t.Fatalf("state = %s/%s, want not_assessed/required_event_family_absent", validation.ValidationState, validation.ReasonCode)
	}
}

func TestObserveRequiredOptionErrors(t *testing.T) {
	for name, tc := range map[string]struct {
		opts    ObserveOptions
		wantErr string
	}{
		"missing-profile": {
			opts:    ObserveOptions{SourcePath: "events.jsonl", OutDir: "run"},
			wantErr: "harness observe requires --profile",
		},
		"missing-source": {
			opts:    ObserveOptions{ProfilePath: "profile.json", OutDir: "run"},
			wantErr: "harness observe requires --source",
		},
		"missing-out": {
			opts:    ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl"},
			wantErr: "harness observe requires --out",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if err := requireObserveOptions(tc.opts); err == nil || err.Error() != tc.wantErr {
				t.Fatalf("requireObserveOptions() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestObserveRejectsUnsafeRawPromptAndDoesNotWriteRun(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	line := eventMap("e1", "harness")
	line["raw_prompt"] = "do the thing"
	writeEvents(t, dir, []map[string]any{line})
	oldwd := chdir(t, dir)
	defer oldwd()
	_, err := Observe(ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl", OutDir: "run"})
	if err == nil || !strings.Contains(err.Error(), "unsafe_input") {
		t.Fatalf("Observe() error = %v, want unsafe_input", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "run")); !os.IsNotExist(statErr) {
		t.Fatalf("run dir exists after unsafe observe: %v", statErr)
	}
}

func TestFindUnsafeAtReasonCodes(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		wantField  string
		wantReason string
	}{
		{
			name:       "raw field",
			value:      map[string]any{"nested": map[string]any{"raw_prompt": "secret"}},
			wantField:  "nested.raw_prompt",
			wantReason: "forbidden_raw_field",
		},
		{
			name:       "sensitive field",
			value:      map[string]any{"authorization": "Bearer token"},
			wantField:  "authorization",
			wantReason: "sensitive_field",
		},
		{
			name:       "unsafe path",
			value:      map[string]any{"items": []any{"../secret.txt"}},
			wantField:  "items[0]",
			wantReason: "unsafe_path_or_private_path",
		},
		{
			name:       "authenticated url",
			value:      map[string]any{"url": "https://example.test/callback?token=secret"},
			wantField:  "url",
			wantReason: "authenticated_url",
		},
		{
			name:       "token value",
			value:      map[string]any{"value": "sk-abcdefghijklmnop"},
			wantField:  "value",
			wantReason: "token_like_value",
		},
		{
			name:       "digest exemption",
			value:      map[string]any{"source_digest": strings.Repeat("a", 64)},
			wantField:  "",
			wantReason: "",
		},
		{
			name:       "empty string",
			value:      "",
			wantField:  "",
			wantReason: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, reason := findUnsafe(tt.value)
			if field != tt.wantField || reason != tt.wantReason {
				t.Fatalf("findUnsafe() = %q/%q, want %q/%q", field, reason, tt.wantField, tt.wantReason)
			}
		})
	}
}

func TestFindUnsafeRawEventAtReasonCodes(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		wantField  string
		wantReason string
	}{
		{
			name:       "raw prompt still forbidden outside retained tool input",
			value:      map[string]any{"raw_prompt": "secret"},
			wantField:  "raw_prompt",
			wantReason: "forbidden_raw_field",
		},
		{
			name:       "retained tool input prompt skipped",
			value:      map[string]any{"part": map[string]any{"state": map[string]any{"input": map[string]any{"prompt": "secret"}}}},
			wantField:  "",
			wantReason: "",
		},
		{
			name:       "unstructured body skipped",
			value:      map[string]any{"message": map[string]any{"content": "raw model text"}},
			wantField:  "",
			wantReason: "",
		},
		{
			name:       "structured body inspected",
			value:      map[string]any{"message": map[string]any{"content": map[string]any{"access_token": "secret"}}},
			wantField:  "message.content.access_token",
			wantReason: "sensitive_field",
		},
		{
			name:       "raw path token exemption",
			value:      map[string]any{"file_path": strings.Repeat("A", 48)},
			wantField:  "",
			wantReason: "",
		},
		{
			name:       "raw url still unsafe",
			value:      map[string]any{"url": "https://example.test/callback?api_key=secret"},
			wantField:  "url",
			wantReason: "authenticated_url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, reason := findUnsafeRawEvent(tt.value)
			if field != tt.wantField || reason != tt.wantReason {
				t.Fatalf("findUnsafeRawEvent() = %q/%q, want %q/%q", field, reason, tt.wantField, tt.wantReason)
			}
		})
	}
}

func TestFindStringByKeyIn(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		wanted map[string]bool
		want   string
	}{
		{
			name:   "direct hit",
			value:  map[string]any{"name": "agent"},
			wanted: map[string]bool{"name": true},
			want:   "agent",
		},
		{
			name:   "skip empty string and continue in list",
			value:  []any{map[string]any{"name": "   "}, map[string]any{"name": "agent"}},
			wanted: map[string]bool{"name": true},
			want:   "agent",
		},
		{
			name:   "nested map lookup",
			value:  []any{map[string]any{"meta": map[string]any{"name": "agent"}}},
			wanted: map[string]bool{"name": true},
			want:   "agent",
		},
		{
			name:   "case-insensitive key match",
			value:  map[string]any{"Display": "Agent"},
			wanted: map[string]bool{"display": true},
			want:   "Agent",
		},
		{
			name:   "not found",
			value:  map[string]any{"note": "noop"},
			wanted: map[string]bool{"name": true},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findStringByKeyIn(tt.value, tt.wanted)
			if got != tt.want {
				t.Fatalf("findStringByKeyIn() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFindNumberByKeyIn(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		wanted map[string]bool
		want   float64
		ok     bool
	}{
		{
			name:   "float64 direct",
			value:  map[string]any{"count": 42.5},
			wanted: map[string]bool{"count": true},
			want:   42.5,
			ok:     true,
		},
		{
			name:   "int direct",
			value:  map[string]any{"count": int(7)},
			wanted: map[string]bool{"count": true},
			want:   7,
			ok:     true,
		},
		{
			name:   "skip non-number and continue in list",
			value:  []any{map[string]any{"count": "ignore"}, map[string]any{"count": int(9)}},
			wanted: map[string]bool{"count": true},
			want:   9,
			ok:     true,
		},
		{
			name:   "nested map lookup",
			value:  map[string]any{"meta": map[string]any{"count": 3.14}},
			wanted: map[string]bool{"count": true},
			want:   3.14,
			ok:     true,
		},
		{
			name:   "case-insensitive key match",
			value:  map[string]any{"Count": 8},
			wanted: map[string]bool{"count": true},
			want:   8,
			ok:     true,
		},
		{
			name:   "not found",
			value:  map[string]any{"meta": "noop"},
			wanted: map[string]bool{"count": true},
			want:   0,
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := findNumberByKeyIn(tt.value, tt.wanted)
			if ok != tt.ok {
				t.Fatalf("findNumberByKeyIn() ok = %t, want %t", ok, tt.ok)
			}
			if ok && got != tt.want {
				t.Fatalf("findNumberByKeyIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEffectiveEventLimitsDefaultsAndOverrides(t *testing.T) {
	maxLine, maxEvents := effectiveEventLimits(Limits{})
	if maxLine != DefaultMaxLineBytes || maxEvents != DefaultMaxEvents {
		t.Fatalf("effectiveEventLimits(defaults) = %d/%d, want %d/%d", maxLine, maxEvents, DefaultMaxLineBytes, DefaultMaxEvents)
	}

	maxLine, maxEvents = effectiveEventLimits(Limits{MaxLineBytes: 32, MaxEvents: 7})
	if maxLine != 32 || maxEvents != 7 {
		t.Fatalf("effectiveEventLimits(overrides) = %d/%d, want 32/7", maxLine, maxEvents)
	}
}

func TestTimestampForKeyPrefersRFC3339ThenUnix(t *testing.T) {
	got := timestampForKey(map[string]any{"timestamp": "2026-05-10T15:00:00+03:00"}, "timestamp")
	if got != "2026-05-10T12:00:00Z" {
		t.Fatalf("timestampForKey(RFC3339) = %s, want UTC normalized value", got)
	}

	got = timestampForKey(map[string]any{"timestamp": float64(1_746_878_400_000)}, "timestamp")
	if got != "2025-05-10T12:00:00Z" {
		t.Fatalf("timestampForKey(unix ms) = %s, want 2025-05-10T12:00:00Z", got)
	}
}

func TestNormalizeOpenCodeRawLineClassifiesFamilies(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	raw := map[string]any{
		"type":      "tool.call",
		"model":     "minimax/MiniMax-M2.7",
		"role":      "assistant",
		"tool":      "edit",
		"phase":     "gsd.phase.1",
		"timestamp": "2026-05-10T11:59:00Z",
	}

	events := normalizeOpenCodeRawLine(raw, 7, now)
	got := map[string]Event{}
	for _, event := range events {
		got[event.EventFamily] = event
	}
	wantFamilies := []string{"interaction", "model", "mutation", "phase", "tool"}
	if len(got) != len(wantFamilies) {
		t.Fatalf("families = %v, want exactly %v", keys(got), wantFamilies)
	}
	for _, family := range wantFamilies {
		event, ok := got[family]
		if !ok {
			t.Fatalf("missing family %s in %+v", family, events)
		}
		if event.SourceRef != "raw-000007" {
			t.Fatalf("%s source_ref = %s", family, event.SourceRef)
		}
		if event.ActorRef != "minimax-MiniMax-M2.7" {
			t.Fatalf("%s actor_ref = %s", family, event.ActorRef)
		}
		if event.ObservedAt != "2026-05-10T11:59:00Z" {
			t.Fatalf("%s observed_at = %s", family, event.ObservedAt)
		}
	}
}

func keys(events map[string]Event) []string {
	out := make([]string, 0, len(events))
	for key := range events {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func TestNormalizeOpenCodeRawLineUsesDefaultActorAndTime(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	events := normalizeOpenCodeRawLine(map[string]any{"type": "session.started"}, 1, now)
	if len(events) != 1 {
		t.Fatalf("events = %+v, want one harness event", events)
	}
	if events[0].ActorRef != "opencode" {
		t.Fatalf("actor_ref = %s, want opencode", events[0].ActorRef)
	}
	if events[0].ObservedAt != now.Format(time.RFC3339) {
		t.Fatalf("observed_at = %s, want %s", events[0].ObservedAt, now.Format(time.RFC3339))
	}
}

func TestNormalizeOpenCodeRawLineBytesRejectsMalformedAndUnsafeRawInput(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	if events, err := normalizeOpenCodeRawLineBytes([]byte("  "), 1, now); err != nil || len(events) != 0 {
		t.Fatalf("blank line events/error = %+v/%v, want none/nil", events, err)
	}
	if _, err := normalizeOpenCodeRawLineBytes([]byte("{"), 2, now); err == nil || !strings.Contains(err.Error(), "raw source line 2: malformed_jsonl") {
		t.Fatalf("malformed error = %v", err)
	}
	if _, err := normalizeOpenCodeRawLineBytes([]byte(`{"raw_prompt":"secret"}`), 3, now); err == nil || !strings.Contains(err.Error(), "unsafe_input:raw_prompt:forbidden_raw_field") {
		t.Fatalf("unsafe error = %v", err)
	}
}

func TestNormalizeOpenCodeRawLineBytesComputesDigestForEachEvent(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	events, err := normalizeOpenCodeRawLineBytes([]byte(`{"type":"tool.call","model":"qwen3.6","timestamp":"2026-05-10T12:00:00Z"}`), 4, now)
	if err != nil {
		t.Fatalf("normalizeOpenCodeRawLineBytes() error = %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("events = %+v, want model and tool", events)
	}
	for _, event := range events {
		if event.SourceDigest == "" {
			t.Fatalf("%s missing source digest: %+v", event.EventFamily, event)
		}
	}
}

func TestWriteNormalizedEventsWritesJSONL(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "nested", "events.jsonl")
	events := []Event{
		normalizedEvent("e1", "harness", "harness_observed", "2026-05-10T12:00:00Z", "raw-000001", "opencode"),
		normalizedEvent("e2", "model", "model_observed", "2026-05-10T12:00:01Z", "raw-000002", "qwen"),
	}
	if err := writeNormalizedEvents(outPath, events); err != nil {
		t.Fatalf("writeNormalizedEvents() error = %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read normalized events: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("lines = %d, want 2\n%s", len(lines), string(data))
	}
	var got Event
	if err := json.Unmarshal([]byte(lines[1]), &got); err != nil {
		t.Fatalf("parse second line: %v", err)
	}
	if got.EventID != "e2" || got.EventFamily != "model" {
		t.Fatalf("second event = %+v", got)
	}
}

func TestSessionCommandFactsRequirePassedModelAndUseSafeTimestamp(t *testing.T) {
	if facts := sessionCommandFacts(SessionRun{CommandModelState: StateCannotVerify, CommandModel: "qwen"}); len(facts) != 0 {
		t.Fatalf("facts for cannot_verify model = %+v, want none", facts)
	}
	if facts := sessionCommandFacts(SessionRun{CommandModelState: StatePass, CommandModel: "   "}); len(facts) != 0 {
		t.Fatalf("facts for blank model = %+v, want none", facts)
	}

	facts := sessionCommandFacts(SessionRun{
		CommandModelState: StatePass,
		CommandModel:      "openrouter/qwen/qwen3.6-plus",
		StartTime:         "not-rfc3339",
		CreatedAt:         "2026-05-10T12:00:00Z",
	})
	if len(facts) != 1 {
		t.Fatalf("facts = %+v, want one", facts)
	}
	if facts[0].EventID != "session-command-model" || facts[0].EventFamily != "model" {
		t.Fatalf("fact identity = %+v", facts[0])
	}
	if facts[0].ActorRef != "openrouter-qwen-qwen3.6-plus" {
		t.Fatalf("actor_ref = %s", facts[0].ActorRef)
	}
	if _, err := time.Parse(time.RFC3339, facts[0].ObservedAt); err != nil {
		t.Fatalf("observed_at is not RFC3339: %s", facts[0].ObservedAt)
	}
	if facts[0].SourceDigest == "" {
		t.Fatalf("missing source digest: %+v", facts[0])
	}
}

func TestLoadSessionProfileDefaultsAndRejectsInvalidRawConfig(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()
	path := "session-profile.json"
	profile := SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
		SetupActions: []SessionSetupAction{{
			ID:       "init",
			Kind:     "init",
			Required: true,
		}},
	}
	writeJSONFixture(t, path, profile)

	loaded, err := LoadSessionProfile(path)
	if err != nil {
		t.Fatalf("LoadSessionProfile() error = %v", err)
	}
	if loaded.StreamCapture != "disabled" {
		t.Fatalf("StreamCapture = %s, want disabled", loaded.StreamCapture)
	}

	profile.RawEventFormat = OpenCodeJSONLRawFormat
	profile.RawEventSourcePath = ""
	writeJSONFixture(t, path, profile)
	if _, err := LoadSessionProfile(path); err == nil || !strings.Contains(err.Error(), "raw_event_source_path required") {
		t.Fatalf("LoadSessionProfile() raw config error = %v", err)
	}
}

func TestLoadSessionProfileRejectsUnsafeSetupAction(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()
	path := "session-profile.json"
	writeJSONFixture(t, path, SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
		SetupActions: []SessionSetupAction{{
			ID:   "../init",
			Kind: "init",
		}},
	})

	if _, err := LoadSessionProfile(path); err == nil || !strings.Contains(err.Error(), "unsafe setup action id") {
		t.Fatalf("LoadSessionProfile() error = %v, want unsafe setup action id", err)
	}
}

func TestSetupSessionRequireOptions(t *testing.T) {
	dir := t.TempDir()
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
	})
	oldwd := chdir(t, dir)
	defer oldwd()

	for name, tc := range map[string]struct {
		opts    SessionSetupOptions
		wantErr string
	}{
		"missing-profile": {
			opts:    SessionSetupOptions{OutDir: "run"},
			wantErr: "observe setup requires --profile",
		},
		"missing-out": {
			opts:    SessionSetupOptions{ProfilePath: "session-profile.json"},
			wantErr: "observe setup requires --out",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, err := validateSessionSetupOptions(tc.opts); err == nil || err.Error() != tc.wantErr {
				t.Fatalf("validateSessionSetupOptions() error = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestSetupSessionRejectsInvalidOptions(t *testing.T) {
	dir := t.TempDir()
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
	})
	oldwd := chdir(t, dir)
	defer oldwd()

	if _, err := SetupSession(SessionSetupOptions{
		ProfilePath: "../session-profile.json",
		OutDir:      "run",
	}); err == nil || err.Error() != "unsafe profile path: path must be relative local file without traversal" {
		t.Fatalf("SetupSession() error = %v, want unsafe profile path: path must be relative local file without traversal", err)
	}
	if _, err := SetupSession(SessionSetupOptions{
		ProfilePath: "session-profile.json",
		OutDir:      "../run",
	}); err == nil || err.Error() != "out must be a relative local directory without traversal" {
		t.Fatalf("SetupSession() error = %v, want out path traversal error", err)
	}
}

func TestSetupSessionRejectsInvalidProfilePayload(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "session-profile.json"), []byte(`{"schema_version":"bad"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	oldwd := chdir(t, dir)
	defer oldwd()

	if _, err := SetupSession(SessionSetupOptions{
		ProfilePath: "session-profile.json",
		OutDir:      "run",
	}); err == nil || !strings.Contains(err.Error(), "unsupported session profile schema_version") {
		t.Fatalf("SetupSession() error = %v, want unsupported session profile schema_version", err)
	}
}

func TestSetupSessionWritesSessionRunWithCommand(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	oldwd := chdir(t, dir)
	defer oldwd()
	sessionProfile := SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
		SetupActions: []SessionSetupAction{{
			ID:       "init",
			Kind:     "init",
			Required: true,
		}},
	}
	writeJSONFixture(t, "session-profile.json", sessionProfile)

	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	command := "--model=minimax-coding-plan/MiniMax-M2.7"
	run, err := SetupSession(SessionSetupOptions{
		ProfilePath: "session-profile.json",
		OutDir:      "run",
		Command:     command,
		Now:         now,
	})
	if err != nil {
		t.Fatalf("SetupSession() error = %v", err)
	}
	if run.CreatedAt != now.Format(time.RFC3339) {
		t.Fatalf("CreatedAt = %s, want %s", run.CreatedAt, now.Format(time.RFC3339))
	}
	if run.CommandDigest == "" || run.CommandDigestState != StatePass {
		t.Fatalf("command digest state = %s/%s, want pass", run.CommandDigestState, run.CommandDigest)
	}
	if run.CommandModel != "minimax-coding-plan/MiniMax-M2.7" || run.CommandModelState != StatePass {
		t.Fatalf("command model = %s/%s, want minimax-coding-plan/MiniMax-M2.7/pass", run.CommandModel, run.CommandModelState)
	}
	if run.SourceCommitState != StateCannotVerify || run.SourceCommit != "" {
		t.Fatalf("source commit state/commit = %s/%s, want cannot_verify with empty source_commit in test workspace", run.SourceCommitState, run.SourceCommit)
	}
	var saved SessionRun
	readJSONFixture(t, filepath.Join("run", "session.json"), &saved)
	if saved.CommandDigest != run.CommandDigest || saved.CommandDigestState != run.CommandDigestState {
		t.Fatalf("saved command digest = %s/%s, run = %s/%s", saved.CommandDigest, saved.CommandDigestState, run.CommandDigest, run.CommandDigestState)
	}
}

func TestSetupSessionWritesBlankCommandDefaults(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	oldwd := chdir(t, dir)
	defer oldwd()
	sessionProfile := SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
	}
	writeJSONFixture(t, "session-profile.json", sessionProfile)

	run, err := SetupSession(SessionSetupOptions{
		ProfilePath: "session-profile.json",
		OutDir:      "run",
		Command:     "   ",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("SetupSession() error = %v", err)
	}
	if run.CommandDigest != "" || run.CommandDigestState != StateCannotVerify {
		t.Fatalf("command digest state = %s/%s, want cannot_verify/", run.CommandDigestState, run.CommandDigest)
	}
	if run.CommandModel != "" || run.CommandModelState != "" {
		t.Fatalf("command model state = %s/%s, want empty/empty", run.CommandModelState, run.CommandModel)
	}
}

func TestSetupSessionCommandRejectsModelAndWritesDigest(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	oldwd := chdir(t, dir)
	defer oldwd()
	sessionProfile := SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
	}
	writeJSONFixture(t, "session-profile.json", sessionProfile)

	run, err := SetupSession(SessionSetupOptions{
		ProfilePath: "session-profile.json",
		OutDir:      "run",
		Command:     "opencode run --model model name",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("SetupSession() error = %v", err)
	}
	if run.CommandDigest == "" || run.CommandDigestState != StatePass {
		t.Fatalf("command digest state = %s/%s, want pass", run.CommandDigestState, run.CommandDigest)
	}
	if run.CommandModel != "" || run.CommandModelState != "" {
		t.Fatalf("command model state = %s/%s, want empty/empty", run.CommandModelState, run.CommandModel)
	}
}

func TestCollectSessionWritesObservedRun(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeEvents(t, dir, []map[string]any{eventMap("e1", "harness")})
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
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "events.jsonl",
		CreatedAt:          "2026-05-10T12:00:00Z",
	})

	oldwd := chdir(t, dir)
	defer oldwd()
	session, observed, err := CollectSession(SessionCollectOptions{
		ProfilePath: "session-profile.json",
		RunDir:      "session-run",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("CollectSession() error = %v", err)
	}
	if session.CollectionState != StatePass || session.ObservedRunDir != "observed" {
		t.Fatalf("session = %+v", session)
	}
	if observed.EventCount != 1 || len(observed.EventRefs) != 1 {
		t.Fatalf("observed = %+v", observed)
	}
	var written Run
	readJSONFixture(t, filepath.Join(dir, "session-run", "observed", "run.json"), &written)
	if written.SchemaVersion != RunSchemaVersion ||
		written.ProfileID != "generic-harness-v1" ||
		written.HarnessFamily != "generic-harness" ||
		written.EventSchemaVersion != EventSchemaVersion ||
		written.SourcePath != "events.jsonl" ||
		written.SourceDigest == "" ||
		written.CreatedAt != "2026-05-10T12:00:00Z" {
		t.Fatalf("written observed run = %+v", written)
	}
}

func TestCollectSessionMarksMissingSourceUnavailable(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "missing-events.jsonl",
	})
	runDir := filepath.Join(dir, "session-run")
	if err := os.Mkdir(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJSONFixture(t, filepath.Join(runDir, "session.json"), SessionRun{
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "missing-events.jsonl",
		CreatedAt:          "2026-05-10T12:00:00Z",
	})

	oldwd := chdir(t, dir)
	defer oldwd()
	session, observed, err := CollectSession(SessionCollectOptions{
		ProfilePath: "session-profile.json",
		RunDir:      "session-run",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("CollectSession() error = %v", err)
	}
	if session.CollectionState != StateCannotVerify || session.CollectionReason != "source_unavailable" {
		t.Fatalf("session = %+v", session)
	}
	if observed.EventCount != 0 || observed.SourcePath != "missing-events.jsonl" {
		t.Fatalf("observed = %+v", observed)
	}
}

func TestCollectSessionNormalizesRawEventsWhenSourceIsMissing(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"tool"}, nil)
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(`{"type":"tool_use","model":"qwen","tool":"edit"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "normalized-events.jsonl",
		RawEventSourcePath: "opencode-raw.jsonl",
		RawEventFormat:     OpenCodeJSONLRawFormat,
	})
	runDir := filepath.Join(dir, "session-run")
	if err := os.Mkdir(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJSONFixture(t, filepath.Join(runDir, "session.json"), SessionRun{
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "normalized-events.jsonl",
		RawEventSourcePath: "opencode-raw.jsonl",
		RawEventFormat:     OpenCodeJSONLRawFormat,
		CommandModel:       "qwen",
		CommandModelState:  StatePass,
		CreatedAt:          "2026-05-10T12:00:00Z",
	})

	oldwd := chdir(t, dir)
	defer oldwd()
	session, observed, err := CollectSession(SessionCollectOptions{
		ProfilePath: "session-profile.json",
		RunDir:      "session-run",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("CollectSession() error = %v", err)
	}
	if session.CollectionState != StatePass || session.NormalizedDigest == "" {
		t.Fatalf("session = %+v, want normalized pass", session)
	}
	if observed.EventCount != 4 || observed.SourcePath != "normalized-events.jsonl" {
		t.Fatalf("observed = %+v, want four normalized events", observed)
	}
	if _, err := os.Stat(filepath.Join(dir, "normalized-events.jsonl")); err != nil {
		t.Fatalf("normalized source was not written: %v", err)
	}
}

func TestCollectSessionPropagatesUnsafeRawSourcePath(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "normalized-events.jsonl",
		RawEventSourcePath: "../raw.jsonl",
		RawEventFormat:     OpenCodeJSONLRawFormat,
	})
	runDir := filepath.Join(dir, "session-run")
	if err := os.Mkdir(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJSONFixture(t, filepath.Join(runDir, "session.json"), SessionRun{
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "normalized-events.jsonl",
		RawEventSourcePath: "../raw.jsonl",
		RawEventFormat:     OpenCodeJSONLRawFormat,
		CreatedAt:          "2026-05-10T12:00:00Z",
	})

	oldwd := chdir(t, dir)
	defer oldwd()
	_, _, err := CollectSession(SessionCollectOptions{
		ProfilePath: "session-profile.json",
		RunDir:      "session-run",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err == nil || !strings.Contains(err.Error(), "raw_event_source_path invalid") {
		t.Fatalf("CollectSession() error = %v, want raw_event_source_path invalid", err)
	}
}

func TestCollectSessionPropagatesUnsafeRawEventContent(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	if err := os.WriteFile(filepath.Join(dir, "opencode-raw.jsonl"), []byte(`{"type":"tool_use","raw_prompt":"do secret work"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeJSONFixture(t, filepath.Join(dir, "session-profile.json"), SessionProfile{
		SchemaVersion:      SessionProfileSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "normalized-events.jsonl",
		RawEventSourcePath: "opencode-raw.jsonl",
		RawEventFormat:     OpenCodeJSONLRawFormat,
	})
	runDir := filepath.Join(dir, "session-run")
	if err := os.Mkdir(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeJSONFixture(t, filepath.Join(runDir, "session.json"), SessionRun{
		SchemaVersion:      SessionRunSchemaVersion,
		ProfileID:          "session-profile",
		HarnessProfilePath: "profile.json",
		EventSourcePath:    "normalized-events.jsonl",
		RawEventSourcePath: "opencode-raw.jsonl",
		RawEventFormat:     OpenCodeJSONLRawFormat,
		CreatedAt:          "2026-05-10T12:00:00Z",
	})

	oldwd := chdir(t, dir)
	defer oldwd()
	_, _, err := CollectSession(SessionCollectOptions{
		ProfilePath: "session-profile.json",
		RunDir:      "session-run",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err == nil || !strings.Contains(err.Error(), "unsafe_input:raw_prompt:forbidden_raw_field") {
		t.Fatalf("CollectSession() error = %v, want unsafe raw_prompt hard error", err)
	}
}

func TestExtractCommandModel(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		want    string
	}{
		{name: "long flag", command: []string{"opencode", "run", "--model", "minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "long equals flag", command: []string{"opencode", "run", "--model=minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "short flag", command: []string{"opencode", "run", "-m", "minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "shell command flag", command: []string{"sh", "-c", `opencode run --format json --model minimax-coding-plan/MiniMax-M2.5 --dir demo "Respond with OK only." > opencode-gsd-run.jsonl`}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "shell command equals flag", command: []string{"sh", "-c", "opencode run --format json --model=minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "bash shell command flag", command: []string{"bash", "-c", "opencode run --model minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "shell command preferred over positional argv", command: []string{"sh", "-c", "opencode run --model minimax-coding-plan/MiniMax-M2.5", "--model", "wrong-model"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "shell quoted prompt flag ignored", command: []string{"sh", "-c", `opencode run 'please ignore --model fake' --model minimax-coding-plan/MiniMax-M2.5`}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "missing value", command: []string{"opencode", "run", "--model"}, want: ""},
		{name: "unsafe url", command: []string{"opencode", "run", "--model", "https://example.com/model"}, want: ""},
		{name: "unsafe shell", command: []string{"opencode", "run", "--model", "model$(touch x)"}, want: ""},
		{name: "unsafe shell command model", command: []string{"sh", "-c", "opencode run --model https://example.com/model"}, want: ""},
		{name: "unsafe shell escape", command: []string{"sh", "-c", `opencode run --model "model\name"`}, want: ""},
		{name: "unsafe whitespace", command: []string{"opencode", "run", "--model", "model name"}, want: ""},
		{name: "unsafe newline", command: []string{"sh", "-c", "opencode run --model 'model\nname'"}, want: ""},
		{name: "unsafe path", command: []string{"opencode", "run", "--model", "../model"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractCommandModel(tt.command); got != tt.want {
				t.Fatalf("extractCommandModel() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShellFieldsControlledSyntax(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    []string
	}{
		{
			name:    "plain whitespace",
			command: "opencode run --model minimax/MiniMax-M2.7",
			want:    []string{"opencode", "run", "--model", "minimax/MiniMax-M2.7"},
		},
		{
			name:    "quoted prompt",
			command: `opencode run 'please ignore --model fake' --model minimax/MiniMax-M2.7`,
			want:    []string{"opencode", "run", "please ignore --model fake", "--model", "minimax/MiniMax-M2.7"},
		},
		{
			name:    "double quoted prompt",
			command: `opencode run "hello world" --model=minimax/MiniMax-M2.7`,
			want:    []string{"opencode", "run", "hello world", "--model=minimax/MiniMax-M2.7"},
		},
		{
			name:    "backslash preserved",
			command: `opencode run --model "model\name"`,
			want:    []string{"opencode", "run", "--model", `model\name`},
		},
		{
			name:    "line continuation removed",
			command: "opencode run \\\n--model minimax/MiniMax-M2.7",
			want:    []string{"opencode", "run", "--model", "minimax/MiniMax-M2.7"},
		},
		{
			name:    "trailing escape preserved",
			command: `opencode run --model model\`,
			want:    []string{"opencode", "run", "--model", `model\`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellFields(tt.command); !equalStrings(got, tt.want) {
				t.Fatalf("shellFields() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestObserveRejectsDigestMismatch(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	event := eventMap("e1", "harness")
	event["source_digest"] = strings.Repeat("0", 64)
	writeRawEvents(t, dir, []map[string]any{event})
	oldwd := chdir(t, dir)
	defer oldwd()
	_, err := Observe(ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl", OutDir: "run"})
	if err == nil || !strings.Contains(err.Error(), "source_digest_mismatch") {
		t.Fatalf("Observe() error = %v, want source_digest_mismatch", err)
	}
}

func TestObserveRejectsUnsafeEventIDBeforeWrite(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeEvents(t, dir, []map[string]any{eventMap("..", "harness")})
	oldwd := chdir(t, dir)
	defer oldwd()
	_, err := Observe(ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl", OutDir: "run"})
	if err == nil || !strings.Contains(err.Error(), "unsafe event_id") {
		t.Fatalf("Observe() error = %v, want unsafe event_id", err)
	}
}

func TestValidateEventRejectsInvalidFields(t *testing.T) {
	profile := Profile{EventSchemaVersion: EventSchemaVersion}
	base := Event{
		EventID:            "event-1",
		EventSchemaVersion: EventSchemaVersion,
		EventFamily:        "harness",
		EventType:          "observe",
		ObservedAt:         "2026-05-10T12:00:00Z",
		SourceRef:          "source-1",
		SourceDigest:       strings.Repeat("a", 64),
		TaskRef:            "task-1",
		OperationRef:       "adapter-run:fixture",
		ActorRef:           "actor-1",
		ContentState:       ContentDigestOnly,
		UnavailableFields: []UnavailableField{{
			Field:      "raw_prompt",
			State:      StateNotAssessed,
			ReasonCode: "redacted",
		}},
	}
	tests := []struct {
		name string
		edit func(*Event)
		want string
	}{
		{name: "schema", edit: func(event *Event) { event.EventSchemaVersion = "old" }, want: "schema_version_mismatch"},
		{name: "family", edit: func(event *Event) { event.EventFamily = "unknown" }, want: "unsupported event_family"},
		{name: "event type", edit: func(event *Event) { event.EventType = "../bad" }, want: "unsafe event_type"},
		{name: "observed at", edit: func(event *Event) { event.ObservedAt = "not-time" }, want: "invalid observed_at"},
		{name: "source ref", edit: func(event *Event) { event.SourceRef = "../source" }, want: "unsafe source_ref"},
		{name: "source digest", edit: func(event *Event) { event.SourceDigest = "bad" }, want: "invalid source_digest"},
		{name: "task ref", edit: func(event *Event) { event.TaskRef = "../task" }, want: "unsafe task_ref"},
		{name: "operation ref", edit: func(event *Event) { event.OperationRef = "adapter-run://bad" }, want: "unsafe operation_ref"},
		{name: "actor ref", edit: func(event *Event) { event.ActorRef = "../actor" }, want: "unsafe actor_ref"},
		{name: "content state", edit: func(event *Event) { event.ContentState = "raw" }, want: "invalid content_state"},
		{name: "unavailable field name", edit: func(event *Event) { event.UnavailableFields[0].Field = "../raw_prompt" }, want: "invalid unavailable_fields"},
		{name: "unavailable field state", edit: func(event *Event) { event.UnavailableFields[0].State = StatePass }, want: "invalid unavailable_fields"},
		{name: "unavailable field reason", edit: func(event *Event) { event.UnavailableFields[0].ReasonCode = "../redacted" }, want: "invalid unavailable_fields"},
	}

	if err := validateEvent(profile, base); err != nil {
		t.Fatalf("base validateEvent() error = %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := base
			event.UnavailableFields = append([]UnavailableField(nil), base.UnavailableFields...)
			tt.edit(&event)
			if err := validateEvent(profile, event); err == nil || err.Error() != tt.want {
				t.Fatalf("validateEvent() error = %v, want %s", err, tt.want)
			}
		})
	}
}

func TestValidateWritesCannotVerifyOutForMissingRun(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	oldwd := chdir(t, dir)
	defer oldwd()
	if err := os.Mkdir("empty-run", 0o755); err != nil {
		t.Fatal(err)
	}
	validation, err := Validate(ValidateOptions{ProfilePath: "profile.json", RunDir: "empty-run", OutPath: "validation.json"})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if validation.ValidationState != StateCannotVerify {
		t.Fatalf("ValidationState = %s, want cannot_verify", validation.ValidationState)
	}
	if _, err := os.Stat("validation.json"); err != nil {
		t.Fatalf("validation out not written: %v", err)
	}
}

func TestValidateWritesOutPathWhenPasses(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeEvents(t, dir, []map[string]any{eventMap("e1", "harness")})
	oldwd := chdir(t, dir)
	defer oldwd()

	run, err := Observe(ObserveOptions{
		ProfilePath: "profile.json",
		SourcePath:  "events.jsonl",
		OutDir:      "run",
		Now:         time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if run.EventCount != 1 {
		t.Fatalf("run.EventCount = %d, want 1", run.EventCount)
	}

	validation, err := Validate(ValidateOptions{
		ProfilePath: "profile.json",
		RunDir:      "run",
		OutPath:     "validation.json",
	})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if validation.ValidationState != StatePass {
		t.Fatalf("ValidationState = %s, want pass", validation.ValidationState)
	}
	if _, err := os.Stat("validation.json"); err != nil {
		t.Fatalf("validation out not written: %v", err)
	}

	var onDisk Validation
	raw, err := os.ReadFile("validation.json")
	if err != nil {
		t.Fatalf("os.ReadFile(validation.json) error = %v", err)
	}
	if err := json.Unmarshal(raw, &onDisk); err != nil {
		t.Fatalf("json.Unmarshal(validation.json) error = %v", err)
	}
	if onDisk.ValidationState != StatePass || onDisk.ReasonCode != "all_required_dimensions_observed" {
		t.Fatalf("on-disk validation = %+v", onDisk)
	}
}

func TestValidateCannotVerifyWhenRunFileInvalid(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	runDir := filepath.Join(dir, "run")
	if err := os.Mkdir(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "run.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}

	oldwd := chdir(t, dir)
	defer oldwd()

	validation, err := Validate(ValidateOptions{
		ProfilePath: "profile.json",
		RunDir:      "run",
		OutPath:     "validation.json",
	})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if validation.ValidationState != StateCannotVerify {
		t.Fatalf("ValidationState = %s, want cannot_verify", validation.ValidationState)
	}
	if validation.ReasonCode != "source_unavailable" {
		t.Fatalf("ReasonCode = %s, want source_unavailable", validation.ReasonCode)
	}
	var onDisk Validation
	raw, err := os.ReadFile("validation.json")
	if err != nil {
		t.Fatalf("os.ReadFile(validation.json) error = %v", err)
	}
	if err := json.Unmarshal(raw, &onDisk); err != nil {
		t.Fatalf("json.Unmarshal(validation.json) error = %v", err)
	}
	if onDisk.ValidationState != StateCannotVerify || onDisk.ReasonCode != "source_unavailable" {
		t.Fatalf("on-disk validation = %+v", onDisk)
	}
}

func TestLoadRunRejectsUnsafeEventRefs(t *testing.T) {
	dir := t.TempDir()
	run := Run{
		SchemaVersion:      RunSchemaVersion,
		ProfileID:          "generic-harness-v1",
		HarnessFamily:      "generic-harness",
		EventSchemaVersion: EventSchemaVersion,
		SourcePath:         "events.jsonl",
		SourceDigest:       strings.Repeat("a", 64),
		EventRefs:          []string{"../escape.json"},
	}
	data, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "run.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := LoadRun(dir); err == nil || !strings.Contains(err.Error(), "unsafe event ref") {
		t.Fatalf("LoadRun() error = %v, want unsafe event ref", err)
	}
}

func TestValidateRejectsUnsafePaths(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	oldwd := chdir(t, dir)
	defer oldwd()
	if _, err := Validate(ValidateOptions{ProfilePath: "../profile.json", RunDir: "run"}); err == nil || !strings.Contains(err.Error(), "unsafe profile path") {
		t.Fatalf("Validate() profile error = %v, want unsafe profile path", err)
	}
	if _, err := Validate(ValidateOptions{ProfilePath: "profile.json", RunDir: "../run"}); err == nil || !strings.Contains(err.Error(), "unsafe run path") {
		t.Fatalf("Validate() run error = %v, want unsafe run path", err)
	}
	if err := os.Mkdir("run", 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Validate(ValidateOptions{ProfilePath: "profile.json", RunDir: "run", OutPath: "../validation.json"}); err == nil || !strings.Contains(err.Error(), "unsafe out path") {
		t.Fatalf("Validate() out error = %v, want unsafe out path", err)
	}
}

func TestLoadValidationRejectsUnsafePathAndSchemaVersion(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()
	if _, err := LoadValidation("../validation.json"); err == nil {
		t.Fatalf("LoadValidation accepted unsafe path")
	}
	data := []byte(`{"schema_version":"old","profile_id":"p","harness_family":"h","event_schema_version":"harness-event-v1","validation_state":"pass","reason_code":"ok","dimensions":[],"event_count":0,"non_authority":"boundary"}`)
	if err := os.WriteFile("validation.json", data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadValidation("validation.json"); err == nil || !strings.Contains(err.Error(), "unsupported validation schema_version") {
		t.Fatalf("LoadValidation error = %v, want unsupported schema", err)
	}
}

func TestObserveRejectsOutParentSymlinkEscape(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeEvents(t, dir, []map[string]any{eventMap("e1", "harness")})
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(dir, "escape")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	oldwd := chdir(t, dir)
	defer oldwd()
	_, err := Observe(ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl", OutDir: "escape/run"})
	if err == nil || !strings.Contains(err.Error(), "out parent path escapes") {
		t.Fatalf("Observe() error = %v, want out parent escape", err)
	}
}

func TestObserveRejectsExistingOutSymlinkEscapeAndNonEmptyOut(t *testing.T) {
	dir := t.TempDir()
	writeProfile(t, dir, []string{"harness"}, nil)
	writeEvents(t, dir, []map[string]any{eventMap("e1", "harness")})
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(dir, "out-link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	oldwd := chdir(t, dir)
	defer oldwd()
	_, err := Observe(ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl", OutDir: "out-link"})
	if err == nil || !strings.Contains(err.Error(), "out path escapes") {
		t.Fatalf("Observe() error = %v, want out escape", err)
	}

	if err := os.Mkdir("occupied", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("occupied", "existing.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = Observe(ObserveOptions{ProfilePath: "profile.json", SourcePath: "events.jsonl", OutDir: "occupied"})
	if err == nil || !strings.Contains(err.Error(), "refuses existing non-empty") {
		t.Fatalf("Observe() error = %v, want non-empty out refusal", err)
	}
}

func TestPathEscapesWorkingDirectory(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "..", want: true},
		{path: "../escape", want: true},
		{path: "safe", want: false},
		{path: string(filepath.Separator) + "abs", want: true},
		{path: ".", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := pathEscapesWorkingDirectory(tt.path); got != tt.want {
				t.Fatalf("pathEscapesWorkingDirectory(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestRelativeSymlinkTarget(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(dir, "link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	oldwd := chdir(t, dir)
	defer oldwd()
	rel, err := relativeSymlinkTarget("link")
	if err != nil {
		t.Fatal(err)
	}
	if rel != "target" {
		t.Fatalf("relativeSymlinkTarget() = %q, want target", rel)
	}
}

func TestSafeExistingOutDir(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()
	if err := os.Mkdir("empty-target", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("empty-target", "empty-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	got, err := safeExistingOutDir("empty-link")
	if err != nil {
		t.Fatal(err)
	}
	if got != "empty-target" {
		t.Fatalf("safeExistingOutDir() = %q, want empty-target", got)
	}

	outside := t.TempDir()
	if err := os.Symlink(outside, "escape-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := safeExistingOutDir("escape-link"); err == nil || !strings.Contains(err.Error(), "out path escapes") {
		t.Fatalf("safeExistingOutDir() error = %v, want escape", err)
	}
}

func TestSafeExistingFilePathValidation(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()

	if err := os.WriteFile("events.jsonl", []byte("record\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("run", 0o755); err != nil {
		t.Fatal(err)
	}

	if _, err := safeExistingFile("events.jsonl"); err != nil {
		t.Fatalf("safeExistingFile(\"events.jsonl\") error = %v", err)
	}
	if _, err := safeExistingFile("../events.jsonl"); err == nil || !strings.Contains(err.Error(), "path must be relative local file without traversal") {
		t.Fatalf("safeExistingFile(\"../events.jsonl\") error = %v", err)
	}
	if _, err := safeExistingFile("https://example.com/file"); err == nil || !strings.Contains(err.Error(), "path must be relative local file without traversal") {
		t.Fatalf("safeExistingFile(url) error = %v", err)
	}
	if _, err := safeExistingFile("run"); err == nil || !strings.Contains(err.Error(), "path must be a file") {
		t.Fatalf("safeExistingFile(directory) error = %v", err)
	}
	if err := os.WriteFile("safe-target", []byte("record\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("safe-target", "safe-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	got, err := safeExistingFile("safe-link")
	if err != nil {
		t.Fatalf("safeExistingFile(safe-link) error = %v", err)
	}
	if got != "safe-target" {
		t.Fatalf("safeExistingFile(safe-link) = %q, want safe-target", got)
	}

	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "outside.jsonl"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "outside.jsonl"), "escape-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := safeExistingFile("escape-link"); err == nil || !strings.Contains(err.Error(), "path escapes working directory") {
		t.Fatalf("safeExistingFile(escape-link) error = %v", err)
	}
}

func TestSafeExistingDirPathValidation(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()

	if err := os.Mkdir("events", 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := safeExistingDir("events"); err != nil {
		t.Fatalf("safeExistingDir(\"events\") error = %v", err)
	}
	if _, err := safeExistingDir("../events"); err == nil || !strings.Contains(err.Error(), "path must be relative local directory without traversal") {
		t.Fatalf("safeExistingDir(\"../events\") error = %v", err)
	}
	if _, err := safeExistingDir("https://example.com/dir"); err == nil || !strings.Contains(err.Error(), "path must be relative local directory without traversal") {
		t.Fatalf("safeExistingDir(url) error = %v", err)
	}
	if err := os.WriteFile("events.jsonl", []byte("record\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := safeExistingDir("events.jsonl"); err == nil || !strings.Contains(err.Error(), "path must be a directory") {
		t.Fatalf("safeExistingDir(file) error = %v", err)
	}

	if err := os.WriteFile(filepath.Join("events", "note.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("events", "safe-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	got, err := safeExistingDir("safe-link")
	if err != nil {
		t.Fatalf("safeExistingDir(safe-link) error = %v", err)
	}
	if got != "events" {
		t.Fatalf("safeExistingDir(safe-link) = %q, want events", got)
	}

	outside := t.TempDir()
	if err := os.Mkdir(filepath.Join(outside, "outside-dir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "outside-dir"), "escape-link-dir"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := safeExistingDir("escape-link-dir"); err == nil || !strings.Contains(err.Error(), "path escapes working directory") {
		t.Fatalf("safeExistingDir(escape-link-dir) error = %v", err)
	}
}

func TestOutParentEscapes(t *testing.T) {
	dir := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(dir, "escape-parent")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	oldwd := chdir(t, dir)
	defer oldwd()
	escapes, err := outParentEscapes(filepath.Join("escape-parent", "child"))
	if err != nil {
		t.Fatal(err)
	}
	if !escapes {
		t.Fatalf("outParentEscapes() = false, want true")
	}
	escapes, err = outParentEscapes(filepath.Join("safe-parent", "child"))
	if err != nil {
		t.Fatal(err)
	}
	if escapes {
		t.Fatalf("missing safe parent should not escape")
	}
}

func TestSafeParentDirRejectsUnsafeInput(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "absolute-path", path: string(filepath.Separator) + "tmp"},
		{name: "url", path: "https://example.com/out"},
		{name: "traversal", path: "../run"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := safeParentDir(tt.path); err == nil || !strings.Contains(err.Error(), "relative local directory without traversal") {
				t.Fatalf("safeParentDir(%q) error = %v, want traversal rejection", tt.path, err)
			}
		})
	}
}

func TestSafeParentDirCharacterizesMissingAndExistingPaths(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()

	got, err := safeParentDir("")
	if err != nil || got != "." {
		t.Fatalf("safeParentDir(\"\") = %q/%v, want .", got, err)
	}
	if err := os.MkdirAll(filepath.Join("existing", "leaf"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err = safeParentDir(filepath.Join("missing", "child"))
	if err != nil || got != "." {
		t.Fatalf("safeParentDir(\"missing/child\") = %q/%v, want .", got, err)
	}
	got, err = safeParentDir(filepath.Join("existing", "leaf"))
	if err != nil || got != filepath.Join("existing", "leaf") {
		t.Fatalf("safeParentDir(\"existing/leaf\") = %q/%v, want existing/leaf", got, err)
	}
	got, err = safeParentDir(filepath.Join("existing", "missing-child"))
	if err != nil || got != "existing" {
		t.Fatalf("safeParentDir(\"existing/missing-child\") = %q/%v, want existing", got, err)
	}
}

func TestSafeParentDirResolvesSymlinks(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()

	safeTarget := "safe-link-target"
	if err := os.Mkdir(safeTarget, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(safeTarget, "safe-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	got, err := safeParentDir(filepath.Join("safe-link", "child"))
	if err != nil {
		t.Fatalf("safeParentDir(\"safe-link/child\") error = %v", err)
	}
	if got != safeTarget {
		t.Fatalf("safeParentDir(\"safe-link/child\") = %q, want %q", got, safeTarget)
	}

	outside := t.TempDir()
	if err := os.Symlink(outside, "escape-link"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := safeParentDir(filepath.Join("escape-link", "child")); err == nil || !strings.Contains(err.Error(), "parent path escapes working directory") {
		t.Fatalf("safeParentDir(\"escape-link/child\") error = %v, want parent path escapes", err)
	}
}

func TestEnsureOutDirEmptyOrMissing(t *testing.T) {
	dir := t.TempDir()
	oldwd := chdir(t, dir)
	defer oldwd()
	if got, err := ensureOutDirEmptyOrMissing("missing"); err != nil || got != "missing" {
		t.Fatalf("missing path = %q/%v, want pass", got, err)
	}
	if err := os.Mkdir("empty", 0o755); err != nil {
		t.Fatal(err)
	}
	if got, err := ensureOutDirEmptyOrMissing("empty"); err != nil || got != "empty" {
		t.Fatalf("empty path = %q/%v, want pass", got, err)
	}
	if err := os.Mkdir("occupied-unit", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("occupied-unit", "file.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ensureOutDirEmptyOrMissing("occupied-unit"); err == nil || !strings.Contains(err.Error(), "refuses existing non-empty") {
		t.Fatalf("ensureOutDirEmptyOrMissing() error = %v, want non-empty refusal", err)
	}
}

func TestSummarizeDoesNotPrintUnsafeValues(t *testing.T) {
	validation := Validation{
		ValidationState:    StateCannotVerify,
		ReasonCode:         "source_unavailable",
		ProfileID:          "generic-harness-v1",
		EventSchemaVersion: EventSchemaVersion,
		Dimensions: []Dimension{
			{Family: "model", Required: true, State: StateNotAssessed, ReasonCode: "required_event_family_absent"},
		},
	}
	summary := Summarize(validation)
	for _, forbidden := range []string{"sk-", "/Users/", "access_token=", "raw prompt"} {
		if strings.Contains(summary, forbidden) {
			t.Fatalf("summary leaked forbidden marker %q:\n%s", forbidden, summary)
		}
	}
	if !strings.Contains(summary, "no harness compliance") {
		t.Fatalf("summary missing non-authority boundary:\n%s", summary)
	}
}

func writeProfile(t *testing.T, dir string, required, optional []string) {
	t.Helper()
	profile := Profile{
		SchemaVersion:         ProfileSchemaVersion,
		ProfileID:             "generic-harness-v1",
		HarnessFamily:         "generic-harness",
		EventSchemaVersion:    EventSchemaVersion,
		RequiredEventFamilies: required,
		OptionalEventFamilies: optional,
		RawRetentionPolicy:    "digest_only",
		DegradationRules: map[string]Rule{
			"missing_required_family": {State: StateNotAssessed, ReasonCode: "required_event_family_absent"},
			"missing_optional_family": {State: StateNotAssessed, ReasonCode: "optional_event_family_absent"},
			"source_unavailable":      {State: StateCannotVerify, ReasonCode: "source_unavailable"},
			"unsafe_input":            {State: StateFail, ReasonCode: "unsafe_input"},
			"digest_mismatch":         {State: StateCannotVerify, ReasonCode: "source_digest_mismatch"},
			"schema_version_mismatch": {State: StateCannotVerify, ReasonCode: "schema_version_mismatch"},
			"cross_link_conflict":     {State: StateCannotVerify, ReasonCode: "adapter_harness_state_conflict"},
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

func eventMap(id, family string) map[string]any {
	return map[string]any{
		"event_id":             id,
		"event_schema_version": EventSchemaVersion,
		"event_family":         family,
		"event_type":           family + "_observed",
		"observed_at":          "2026-05-09T12:00:00Z",
		"source_ref":           "src-" + id,
		"source_digest":        "",
		"task_ref":             "task-1",
		"operation_ref":        "op-1",
		"actor_ref":            "agent-1",
		"content_state":        ContentDigestOnly,
	}
}

func writeEvents(t *testing.T, dir string, events []map[string]any) {
	t.Helper()
	for _, event := range events {
		event["source_digest"] = ""
		line := marshalCompact(t, event)
		event["source_digest"] = digestForTest(line)
	}
	writeRawEvents(t, dir, events)
}

func writeRawEvents(t *testing.T, dir string, events []map[string]any) {
	t.Helper()
	var lines []string
	for _, event := range events {
		lines = append(lines, marshalCompact(t, event))
	}
	if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeJSONFixture(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readJSONFixture(t *testing.T, path string, target any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatal(err)
	}
}

func marshalCompact(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func digestForTest(line string) string {
	var raw map[string]any
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		panic(err)
	}
	raw["source_digest"] = ""
	canonical, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:])
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func chdir(t *testing.T, dir string) func() {
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
