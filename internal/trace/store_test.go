package trace

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func mustNewTestEvent(t *testing.T, runID string, sequence int, eventType EventType, prevHash string) Event {
	t.Helper()
	event := Event{
		SchemaVersion: SchemaVersion,
		RunID:         runID,
		EventID:       string(eventType),
		Sequence:      sequence,
		EventType:     eventType,
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		PrevEventHash: prevHash,
		EventPayload:  map[string]any{"state": eventType},
	}
	computed, err := event.WithComputedEventHash()
	if err != nil {
		t.Fatalf("compute event hash: %v", err)
	}
	return computed
}

func TestRunLayoutWritesAndValidatesEventChain(t *testing.T) {
	runDir := t.TempDir()
	layout, err := NewRunLayout(runDir)
	if err != nil {
		t.Fatal(err)
	}
	event := Event{
		SchemaVersion: SchemaVersion,
		RunID:         "run-a",
		EventID:       "event-a",
		Sequence:      0,
		EventType:     EventType("run_started"),
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		PrevEventHash: NullEventHash,
		EventPayload:  map[string]any{"state": "started"},
	}
	event, err = event.WithComputedEventHash()
	if err != nil {
		t.Fatal(err)
	}
	manifest := RunManifest{
		SchemaVersion:   SchemaVersion,
		RunID:           "run-a",
		RecorderVersion: RecorderVersion,
		ContractID:      DefaultContract.ContractID,
		EventCount:      1,
		EventChainHead:  event.EventHash,
	}
	if err := layout.WriteRun(manifest); err != nil {
		t.Fatalf("write run: %v", err)
	}
	if err := layout.WriteEvent(event); err != nil {
		t.Fatalf("write event: %v", err)
	}
	if err := ValidateRunDirectory(runDir, true); err != nil {
		t.Fatalf("validate run dir: %v", err)
	}
	if _, err := OpenRunArtifact(runDir); err != nil {
		t.Fatalf("open run artifact: %v", err)
	}
}

func TestValidateRunDirectoryChecksEventCountAndChainHead(t *testing.T) {
	runDir := t.TempDir()
	layout, err := NewRunLayout(runDir)
	if err != nil {
		t.Fatal(err)
	}
	validEvent := mustNewTestEvent(t, "run-a", 0, EventRunStarted, NullEventHash)
	mismatchedEvent := mustNewTestEvent(t, "run-a", 0, EventRunClosed, NullEventHash)
	manifest := RunManifest{
		SchemaVersion:   SchemaVersion,
		RunID:           "run-a",
		RecorderVersion: RecorderVersion,
		ContractID:      DefaultContract.ContractID,
		EventCount:      1,
		EventChainHead:  validEvent.EventHash,
	}
	if err := layout.WriteRun(manifest); err != nil {
		t.Fatalf("write run: %v", err)
	}
	if err := layout.WriteEvent(validEvent); err != nil {
		t.Fatalf("write event: %v", err)
	}
	manifest.EventCount = 2
	if err := layout.WriteRun(manifest); err != nil {
		t.Fatalf("write malformed run: %v", err)
	}
	if err := ValidateRunDirectory(runDir, true); err == nil {
		t.Fatal("expected event_count mismatch")
	}

	manifest.EventCount = 1
	manifest.EventChainHead = mismatchedEvent.EventHash
	if err := layout.WriteRun(manifest); err != nil {
		t.Fatalf("write malformed chain head run: %v", err)
	}
	if err := ValidateRunDirectory(runDir, true); err == nil {
		t.Fatal("expected chain head mismatch")
	}
}

func TestValidateRunDirectorySkipsEventChainCheckWhenDisabled(t *testing.T) {
	runDir := t.TempDir()
	layout, err := NewRunLayout(runDir)
	if err != nil {
		t.Fatal(err)
	}
	firstEvent := mustNewTestEvent(t, "run-a", 0, EventRunStarted, NullEventHash)
	secondEvent := mustNewTestEvent(t, "run-a", 1, EventRunClosed, "sha256:bad")
	manifest := RunManifest{
		SchemaVersion:   SchemaVersion,
		RunID:           "run-a",
		RecorderVersion: RecorderVersion,
		ContractID:      DefaultContract.ContractID,
		EventCount:      2,
		EventChainHead:  secondEvent.EventHash,
	}
	if err := layout.WriteRun(manifest); err != nil {
		t.Fatalf("write run: %v", err)
	}
	if err := layout.WriteEvent(firstEvent); err != nil {
		t.Fatalf("write first event: %v", err)
	}
	if err := layout.WriteEvent(secondEvent); err != nil {
		t.Fatalf("write second event: %v", err)
	}
	if err := ValidateRunDirectory(runDir, false); err != nil {
		t.Fatalf("expected invalid chain not to be validated: %v", err)
	}
}

func TestAppendRunEventExtendsChainAndUpdatesManifest(t *testing.T) {
	runDir := t.TempDir()
	layout, err := NewRunLayout(runDir)
	if err != nil {
		t.Fatal(err)
	}
	first := mustNewTestEvent(t, "run-a", 0, EventRunStarted, NullEventHash)
	manifest := RunManifest{
		SchemaVersion:   SchemaVersion,
		RunID:           "run-a",
		RecorderVersion: RecorderVersion,
		ContractID:      DefaultContract.ContractID,
		EventCount:      1,
		EventChainHead:  first.EventHash,
		FinalChainHead:  first.EventHash,
	}
	if err := layout.WriteRun(manifest); err != nil {
		t.Fatalf("write run: %v", err)
	}
	if err := layout.WriteEvent(first); err != nil {
		t.Fatalf("write first event: %v", err)
	}

	appended, err := AppendRunEvent(runDir, EventCommandFinished, map[string]any{"exit_code": 0}, "test-observer")
	if err != nil {
		t.Fatalf("AppendRunEvent() error = %v", err)
	}
	if appended.Sequence != 1 {
		t.Fatalf("sequence = %d, want 1", appended.Sequence)
	}
	if appended.PrevEventHash != first.EventHash {
		t.Fatalf("prev hash = %s, want %s", appended.PrevEventHash, first.EventHash)
	}
	if appended.ObservedBy != "test-observer" {
		t.Fatalf("observed_by = %s", appended.ObservedBy)
	}
	if appended.EventHash == "" || appended.PayloadDigest == "" {
		t.Fatalf("append did not compute hashes: %+v", appended)
	}

	artifact, err := OpenRunArtifact(runDir)
	if err != nil {
		t.Fatalf("OpenRunArtifact() error = %v", err)
	}
	if artifact.Manifest.EventCount != 2 {
		t.Fatalf("event_count = %d, want 2", artifact.Manifest.EventCount)
	}
	if artifact.Manifest.EventChainHead != appended.EventHash || artifact.Manifest.FinalChainHead != appended.EventHash {
		t.Fatalf("manifest chain heads not updated: %+v appended=%s", artifact.Manifest, appended.EventHash)
	}
	if err := ValidateRunDirectory(runDir, true); err != nil {
		t.Fatalf("ValidateRunDirectory() error = %v", err)
	}
}

func TestValidateEventChainRejectsBrokenPreviousHash(t *testing.T) {
	first := Event{
		SchemaVersion: SchemaVersion,
		RunID:         "run-a",
		EventID:       "first",
		Sequence:      0,
		EventType:     EventType("run_started"),
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		PrevEventHash: NullEventHash,
		EventPayload:  map[string]any{"state": "started"},
	}
	first, err := first.WithComputedEventHash()
	if err != nil {
		t.Fatal(err)
	}
	second := Event{
		SchemaVersion: SchemaVersion,
		RunID:         "run-a",
		EventID:       "second",
		Sequence:      1,
		EventType:     EventType("run_closed"),
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		PrevEventHash: "sha256:bad",
		EventPayload:  map[string]any{"state": "closed"},
	}
	second, err = second.WithComputedEventHash()
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateEventChain([]Event{first, second}); err == nil {
		t.Fatalf("expected chain validation failure")
	}
}

func TestCopyArtifactFileAndResolveContractPath(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "artifact.txt")
	if err := os.WriteFile(src, []byte("artifact"), 0o644); err != nil {
		t.Fatal(err)
	}
	layout, err := NewRunLayout(filepath.Join(root, "run-a"))
	if err != nil {
		t.Fatal(err)
	}
	if err := layout.WriteRun(RunManifest{SchemaVersion: SchemaVersion, RunID: "run-a", RecorderVersion: RecorderVersion, ContractID: DefaultContract.ContractID}); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(layout.ArtifactsDir, "artifact.txt")
	if err := CopyArtifactFile(src, dst); err != nil {
		t.Fatalf("copy artifact: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("stat copied artifact: %v", err)
	}
	if ResolveContractPath("", "") != "" {
		t.Fatalf("empty contract path should remain empty")
	}
}
