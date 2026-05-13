package recorder

import (
	"os"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// runWriter owns the mutable state for a single local recording. The public
// package API receives only immutable run results after the writer finalizes.

type runWriter struct {
	runDir     string
	contract   trace.Contract
	manifest   trace.RunManifest
	sequence   int
	lastHash   string
	events     []trace.Event
	stdoutHash hashWriter
	stderrHash hashWriter
}

func newRunWriter(runDir string, contract trace.Contract, task string) (*runWriter, error) {
	// Layout creation precedes manifest source capture so a failing filesystem
	// setup cannot produce a partially initialized writer.
	if err := createRunLayoutDirs(runDir); err != nil {
		return nil, err
	}

	manifest := initialRunManifest(contract, task)
	sourceDigest, sourceState := trace.LocalSourceSnapshot(currentWorkingDir())

	// Source state is captured at writer creation because it describes the code
	// tree that produced the recorded command, not a later verifier replay.
	manifest.SourceSnapshot = sourceDigest
	manifest.SourceState = sourceState

	return &runWriter{
		runDir:   runDir,
		contract: contract,
		manifest: manifest,
	}, nil
}

func createRunLayoutDirs(runDir string) error {
	// The recorder keeps evidence, artifacts, verifier output, and exports in
	// stable sibling directories so downstream tools can locate each surface.
	for _, rel := range []string{"events", "artifacts", "verifier", "export"} {
		if err := os.MkdirAll(filepath.Join(runDir, rel), 0o755); err != nil {
			return err
		}
	}
	return nil
}
