package harnessobs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
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

func TestExtractCommandModel(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		want    string
	}{
		{name: "long flag", command: []string{"opencode", "run", "--model", "minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "long equals flag", command: []string{"opencode", "run", "--model=minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "short flag", command: []string{"opencode", "run", "-m", "minimax-coding-plan/MiniMax-M2.5"}, want: "minimax-coding-plan/MiniMax-M2.5"},
		{name: "missing value", command: []string{"opencode", "run", "--model"}, want: ""},
		{name: "unsafe url", command: []string{"opencode", "run", "--model", "https://example.com/model"}, want: ""},
		{name: "unsafe shell", command: []string{"opencode", "run", "--model", "model$(touch x)"}, want: ""},
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
