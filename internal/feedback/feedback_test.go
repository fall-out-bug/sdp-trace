package feedback

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecordRetainsCorrectiveFeedbackMessage(t *testing.T) {
	dir := t.TempDir()
	messagePath := filepath.Join(dir, "message.md")
	message := "User correction: do not edit the GSD plan; observe and report.\n"
	if err := os.WriteFile(messagePath, []byte(message), 0o644); err != nil {
		t.Fatal(err)
	}
	event, err := Record(Options{
		Kind:        "corrective_feedback",
		From:        "human",
		To:          "gsd",
		SourceRef:   "codex-thread",
		Summary:     "User corrected agent boundary for GSD plan observation.",
		MessageFile: messagePath,
		Now:         time.Date(2026, 5, 8, 20, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.SchemaVersion != SchemaVersion {
		t.Fatalf("schema = %s", event.SchemaVersion)
	}
	if event.Kind != "corrective_feedback" || event.TrustScope != "local_structural" || event.ProofState != "not_assessed" {
		t.Fatalf("unexpected event boundary: %+v", event)
	}
	if !event.Message.Retained || event.Message.Body != message || event.Message.Bytes != len(message) {
		t.Fatalf("message not retained correctly: %+v", event.Message)
	}
	if event.Message.SHA256 == "" || !strings.Contains(event.EventID, "corrective_feedback") {
		t.Fatalf("missing digest/id: %+v", event)
	}
	if err := ValidateEvent(event); err != nil {
		t.Fatalf("event should validate: %v", err)
	}
}

func TestWriteJSONWritesEvent(t *testing.T) {
	dir := t.TempDir()
	messagePath := filepath.Join(dir, "message.md")
	if err := os.WriteFile(messagePath, []byte("body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	event, err := Record(Options{
		Kind:        "corrective_feedback",
		From:        "human",
		To:          "opencode",
		Summary:     "Boundary correction.",
		MessageFile: messagePath,
		Now:         time.Date(2026, 5, 8, 20, 31, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join(dir, ".sdp-trace", "feedback", "event.json")
	if err := WriteJSON(outPath, event); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.EventID != event.EventID {
		t.Fatalf("event id = %s", decoded.EventID)
	}
}

func TestRecordRejectsUnsafeTokensAndOversizedMessages(t *testing.T) {
	dir := t.TempDir()
	messagePath := filepath.Join(dir, "message.md")
	if err := os.WriteFile(messagePath, []byte("body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Record(Options{
		Kind:        "corrective_feedback",
		From:        "/tmp/human",
		To:          "gsd",
		Summary:     "Bad token.",
		MessageFile: messagePath,
	})
	if err == nil || !strings.Contains(err.Error(), "--from") {
		t.Fatalf("expected unsafe token error, got %v", err)
	}
	largePath := filepath.Join(dir, "large.md")
	if err := os.WriteFile(largePath, []byte(strings.Repeat("x", MaxBodyBytes+1)), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = Record(Options{
		Kind:        "corrective_feedback",
		From:        "human",
		To:          "gsd",
		Summary:     "Too large.",
		MessageFile: largePath,
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("expected size error, got %v", err)
	}
}

func TestRecordRequiresActorAndTarget(t *testing.T) {
	dir := t.TempDir()
	messagePath := filepath.Join(dir, "message.md")
	if err := os.WriteFile(messagePath, []byte("body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Record(Options{
		Kind:        "corrective_feedback",
		To:          "gsd",
		Summary:     "Missing actor.",
		MessageFile: messagePath,
	})
	if err == nil || !strings.Contains(err.Error(), "--from") {
		t.Fatalf("expected missing from error, got %v", err)
	}
	_, err = Record(Options{
		Kind:        "corrective_feedback",
		From:        "human",
		Summary:     "Missing target.",
		MessageFile: messagePath,
	})
	if err == nil || !strings.Contains(err.Error(), "--to") {
		t.Fatalf("expected missing to error, got %v", err)
	}
}

func TestValidateEventRejectsDigestMismatch(t *testing.T) {
	dir := t.TempDir()
	messagePath := filepath.Join(dir, "message.md")
	if err := os.WriteFile(messagePath, []byte("body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	event, err := Record(Options{
		Kind:        "corrective_feedback",
		From:        "human",
		To:          "gsd",
		Summary:     "Digest mismatch.",
		MessageFile: messagePath,
		Now:         time.Date(2026, 5, 8, 20, 32, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	event.Message.Body = "BODY\n"
	if err := ValidateEvent(event); err == nil || !strings.Contains(err.Error(), "sha256") {
		t.Fatalf("expected digest mismatch, got %v", err)
	}
}
