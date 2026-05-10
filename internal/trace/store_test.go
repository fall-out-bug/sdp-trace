package trace

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
