package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestBuildWriteReadAuditBundle(t *testing.T) {
	runDir := t.TempDir()
	manifest := trace.RunManifest{
		SchemaVersion:   trace.SchemaVersion,
		RunID:           "run-a",
		RecorderVersion: trace.RecorderVersion,
		ContractID:      trace.DefaultContract.ContractID,
		EventCount:      1,
	}
	event := trace.Event{SchemaVersion: trace.SchemaVersion, EventID: "event-a"}
	event.RunID = manifest.RunID
	event.Sequence = 0
	event.EventType = trace.EventType("run_started")
	event.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	event.PrevEventHash = trace.NullEventHash
	event.EventPayload = map[string]any{"state": "started"}
	event, err := event.WithComputedEventHash()
	if err != nil {
		t.Fatal(err)
	}
	writeJSON(t, filepath.Join(runDir, "run.json"), manifest)
	if err := os.Mkdir(filepath.Join(runDir, "events"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeJSON(t, filepath.Join(runDir, "events", "000001-event.json"), event)

	result := trace.VerifierResult{Result: trace.VerdictObserved}
	table := trace.MissingEvidenceTable{Rows: []trace.MissingEvidenceRow{}}
	audit := &trace.IntegrityAudit{Issue: "none"}
	bundle, err := BuildAuditBundle(runDir, result, table, audit, nil)
	if err != nil {
		t.Fatalf("build bundle: %v", err)
	}
	if bundle.Run.RunID != manifest.RunID || len(bundle.Events) != 1 || bundle.Integrity == nil {
		t.Fatalf("unexpected bundle: %+v", bundle)
	}

	out := filepath.Join(t.TempDir(), "bundle.json")
	if err := WriteBundle(out, bundle); err != nil {
		t.Fatalf("write bundle: %v", err)
	}
	loaded, err := Read(out)
	if err != nil {
		t.Fatalf("read bundle: %v", err)
	}
	if loaded.Run.RunID != manifest.RunID {
		t.Fatalf("loaded run id = %q", loaded.Run.RunID)
	}
	if RunManifestPath(runDir) != filepath.Join(runDir, "run.json") {
		t.Fatalf("unexpected run manifest path")
	}
}

func TestBuildAuditBundleRejectsMissingRun(t *testing.T) {
	_, err := BuildAuditBundle(t.TempDir(), trace.VerifierResult{}, trace.MissingEvidenceTable{}, nil, nil)
	if err == nil {
		t.Fatalf("expected missing run error")
	}
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
