package trace

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// EventWriterConfig contains the layout and run metadata for an append-only run artifact.
type EventWriterConfig struct {
	RunDir string
}

// RunLayout represents the stable on-disk arrangement for one run.
type RunLayout struct {
	RunFilePath  string
	EventsDir    string
	ArtifactsDir string
	VerifierDir  string
	ExportDir    string
}

// NewRunLayout creates all child directories and returns paths.
func NewRunLayout(runDir string) (RunLayout, error) {
	layout := newRunLayout(runDir)
	for _, dir := range []string{layout.EventsDir, layout.ArtifactsDir, layout.VerifierDir, layout.ExportDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return RunLayout{}, err
		}
	}
	return layout, nil
}

func newRunLayout(runDir string) RunLayout {
	return RunLayout{
		RunFilePath:  filepath.Join(runDir, "run.json"),
		EventsDir:    filepath.Join(runDir, "events"),
		ArtifactsDir: filepath.Join(runDir, "artifacts"),
		VerifierDir:  filepath.Join(runDir, "verifier"),
		ExportDir:    filepath.Join(runDir, "export"),
	}
}

// EventFileName returns a stable sequence-based event path inside events/.
func EventFileName(sequence int, eventType EventType) string {
	return fmt.Sprintf("%06d-%s.json", sequence, eventType)
}

// WriteEvent stores a deterministic event file under events/.
func (layout RunLayout) WriteEvent(event Event) error {
	event = event.EnsureDefaults()
	filename := EventFileName(event.Sequence, event.EventType)
	path := filepath.Join(layout.EventsDir, filename)
	encoded, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, encoded, 0o644)
}

// WriteRun writes run manifest under run.json.
func (layout RunLayout) WriteRun(manifest RunManifest) error {
	encoded, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(layout.RunFilePath, encoded, 0o644)
}

// RunArtifact is the loaded form for a run directory.
type RunArtifact struct {
	Manifest RunManifest
	Events   []Event
	Layout   RunLayout
}

// OpenRunArtifact loads the manifest and events from disk.
func OpenRunArtifact(runDir string) (RunArtifact, error) {
	layout := newRunLayout(runDir)
	manifestData, err := os.ReadFile(layout.RunFilePath)
	if err != nil {
		return RunArtifact{}, err
	}
	var manifest RunManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return RunArtifact{}, err
	}

	entries, err := os.ReadDir(layout.EventsDir)
	if err != nil {
		return RunArtifact{}, err
	}
	events, err := readRunEvents(layout.EventsDir, entries)
	if err != nil {
		return RunArtifact{}, err
	}

	return RunArtifact{
		Manifest: manifest,
		Events:   events,
		Layout:   layout,
	}, nil
}

// AppendRunEvent appends one local event to an existing run artifact and updates the run manifest chain head.
func AppendRunEvent(runDir string, eventType EventType, payload map[string]any, observedBy string) (Event, error) {
	artifact, err := OpenRunArtifact(runDir)
	if err != nil {
		return Event{}, err
	}
	prevHash := NullEventHash
	if len(artifact.Events) > 0 {
		prevHash = artifact.Events[len(artifact.Events)-1].EventHash
	}
	event := Event{
		SchemaVersion: SchemaVersion,
		RunID:         artifact.Manifest.RunID,
		EventID:       SHA256Hex(fmt.Sprintf("%s:%s:%d:%s", artifact.Manifest.RunID, eventType, len(artifact.Events), time.Now().UTC().Format(time.RFC3339Nano))),
		Sequence:      len(artifact.Events),
		EventType:     eventType,
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		PrevEventHash: prevHash,
		HashAlgorithm: HashAlgSHA256,
		Canonicalization: Canonicalization{
			Algorithm: CanonicalSchemaAlgo,
			Version:   CanonicalAlgoVersion,
		},
		EventPayload: payload,
		ObservedBy:   observedBy,
	}
	computed, err := event.WithComputedEventHash()
	if err != nil {
		return Event{}, err
	}
	event = computed
	if err := artifact.Layout.WriteEvent(event); err != nil {
		return Event{}, err
	}
	artifact.Manifest.EventCount = event.Sequence + 1
	artifact.Manifest.EventChainHead = event.EventHash
	artifact.Manifest.FinalChainHead = event.EventHash
	if err := artifact.Layout.WriteRun(artifact.Manifest); err != nil {
		return Event{}, err
	}
	return event, nil
}

// ValidateRunDirectory checks that run.json and event files are parseable.
func ValidateRunDirectory(path string, requireChain bool) error {
	runArtifact, err := OpenRunArtifact(path)
	if err != nil {
		return err
	}
	return validateRunArtifact(runArtifact, requireChain)
}

func validateRunArtifact(artifact RunArtifact, requireChain bool) error {
	if err := artifact.Manifest.Validate(); err != nil {
		return err
	}
	return validateRunDirectoryState(artifact.Manifest, artifact.Events, requireChain)
}

func validateEventChainIfRequested(events []Event, requireChain bool) error {
	if !requireChain {
		return nil
	}
	return ValidateEventChain(events)
}

func validateManifestEventCount(manifestCount int, eventCount int) error {
	if manifestCount != 0 && manifestCount != eventCount {
		return fmt.Errorf("event_count mismatch: run.json=%d files=%d", manifestCount, eventCount)
	}
	return nil
}

func validateManifestEventChainHead(manifestHead string, events []Event) error {
	if manifestHead == "" || len(events) == 0 {
		return nil
	}
	if events[len(events)-1].EventHash != manifestHead {
		return fmt.Errorf("run manifest event_chain_head does not match last event hash")
	}
	return nil
}

func validateRunDirectoryState(manifest RunManifest, events []Event, requireChain bool) error {
	if err := validateEventChainIfRequested(events, requireChain); err != nil {
		return err
	}
	if err := validateManifestEventCount(manifest.EventCount, len(events)); err != nil {
		return err
	}
	return validateManifestEventChainHead(manifest.EventChainHead, events)
}

// readRunEvents loads and sorts every *.json file in events/.
func readRunEvents(eventsDir string, entries []fs.DirEntry) ([]Event, error) {
	jsonFiles := eventJSONFiles(eventsDir, entries)
	events := make([]Event, 0, len(jsonFiles))
	for _, path := range jsonFiles {
		event, err := readRunEvent(path)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// ValidateEventChain checks that hashes and prev-event links are consistent.
func ValidateEventChain(events []Event) error {
	prevEventHash := NullEventHash
	for i, event := range events {
		if err := validateChainEvent(i, event, prevEventHash); err != nil {
			return err
		}
		prevEventHash = event.EventHash
	}
	return nil
}

func eventJSONFiles(eventsDir string, entries []fs.DirEntry) []string {
	jsonFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		if isEventJSON(entry) {
			jsonFiles = append(jsonFiles, filepath.Join(eventsDir, entry.Name()))
		}
	}
	sort.Strings(jsonFiles)
	return jsonFiles
}

func isEventJSON(entry fs.DirEntry) bool {
	return !entry.IsDir() && filepath.Ext(entry.Name()) == ".json"
}

func readRunEvent(path string) (Event, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return Event{}, err
	}
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return Event{}, err
	}
	if err := event.Validate(); err != nil {
		return Event{}, fmt.Errorf("invalid event %q: %w", path, err)
	}
	return event, nil
}

func validateChainEvent(index int, event Event, prevEventHash string) error {
	computed, err := event.WithComputedEventHash()
	if err != nil {
		return fmt.Errorf("event %d (%s) hash generation failed: %w", index, event.EventID, err)
	}
	if err := validateEventHash(index, event, computed.EventHash); err != nil {
		return err
	}
	if err := event.VerifyPayloadDigest(); err != nil {
		return fmt.Errorf("event %d (%s) payload_digest invalid: %w", index, event.EventID, err)
	}
	return validateEventPosition(index, event, prevEventHash)
}

func validateEventHash(index int, event Event, computedHash string) error {
	if event.EventHash != computedHash {
		return fmt.Errorf("event %d (%s) event_hash mismatch: expected %s got %s", index, event.EventID, computedHash, event.EventHash)
	}
	return nil
}

func validateEventPosition(index int, event Event, prevEventHash string) error {
	if event.Sequence != index {
		return fmt.Errorf("event %d has non-zero-based sequence %d", index, event.Sequence)
	}
	if event.PrevEventHash != prevEventHash {
		return fmt.Errorf("event %d (%s) prev_event_hash expected %s", index, event.EventID, prevEventHash)
	}
	return nil
}

// CopyArtifactFile copies a verifier/export artifact into a run directory.
func CopyArtifactFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		return err
	}
	return output.Sync()
}
