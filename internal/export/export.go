package export

import (
	"encoding/json"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"os"
	"path/filepath"
)

// AuditBundle is a small reproducible package exported by the command layer.
type AuditBundle struct {
	Run       trace.RunManifest          `json:"run"`
	Events    []trace.Event              `json:"events"`
	Result    trace.VerifierResult       `json:"result"`
	Missing   trace.MissingEvidenceTable `json:"missing_evidence"`
	Integrity *trace.IntegrityAudit      `json:"integrity_audit,omitempty"`
}

// BuildAuditBundle composes run-level artifacts into an exportable structure.
func BuildAuditBundle(runDir string, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit, events []trace.Event) (AuditBundle, error) {
	return buildAuditBundle(runDir, result, table, audit, events)
}

// Write writes a deterministic JSON bundle file.
func WriteBundle(path string, bundle AuditBundle) error {
	return writeBundle(path, bundle)
}

// Read reads a prebuilt bundle.
func Read(path string) (AuditBundle, error) {
	return readBundle(path)
}

// RunManifestPath resolves the canonical run manifest location.
func RunManifestPath(runDir string) string {
	return filepath.Join(runDir, "run.json")
}
func buildAuditBundle(runDir string, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit, events []trace.Event) (AuditBundle, error) {

	runArtifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		return AuditBundle{}, err
	}
	events = auditBundleEvents(events, runArtifact.Events)

	return newAuditBundle(runArtifact.Manifest, result, table, audit, events), nil
}

func auditBundleEvents(events, artifactEvents []trace.Event) []trace.Event {

	if len(events) > 0 {
		return events
	}
	return artifactEvents
}

func newAuditBundle(run trace.RunManifest, result trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit, events []trace.Event) AuditBundle {

	return AuditBundle{
		Run:       run,
		Events:    events,
		Result:    result,
		Missing:   table,
		Integrity: audit,
	}
}

func writeBundle(path string, bundle AuditBundle) error {

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Read reads a prebuilt bundle.
func readBundle(path string) (AuditBundle, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return AuditBundle{}, err
	}
	var bundle AuditBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return AuditBundle{}, err
	}
	return bundle, nil
}
