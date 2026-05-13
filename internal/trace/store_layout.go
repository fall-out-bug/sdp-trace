package trace

import (
	"os"
	"path/filepath"
)

const (
	runManifestName  = "run.json"
	eventsDirName    = "events"
	artifactsDirName = "artifacts"
	verifierDirName  = "verifier"
	exportDirName    = "export"
)

// Run layout names are centralized here because manifests, event files, and
// verifier/export artifacts must agree on one portable directory contract.
// The constants are path segments only; callers still choose the root.
// Keeping names local to trace storage avoids leaking harness-specific layout
// assumptions into the portable contract.

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
	// Layout creation materializes only known child directories.
	layout := newRunLayout(runDir)
	for _, dir := range []string{layout.EventsDir, layout.ArtifactsDir, layout.VerifierDir, layout.ExportDir} {
		// Every sibling directory must exist before event or verifier writes.
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return RunLayout{}, err
		}
	}
	return layout, nil
}

func newRunLayout(runDir string) RunLayout {
	// Path derivation is pure; filesystem creation stays in NewRunLayout.
	return RunLayout{
		RunFilePath:  runLayoutPath(runDir, runManifestName),
		EventsDir:    runLayoutPath(runDir, eventsDirName),
		ArtifactsDir: runLayoutPath(runDir, artifactsDirName),
		VerifierDir:  runLayoutPath(runDir, verifierDirName),
		ExportDir:    runLayoutPath(runDir, exportDirName),
	}
}

func runLayoutPath(runDir string, name string) string {
	return filepath.Join(runDir, name)
}
