package recorder

import (
	"path/filepath"
	"time"
)

func (w *runWriter) finalize(closureState string) error {
	// Finalization writes the terminal manifest before artifact digests so a run
	// reader can distinguish closed metadata from missing digest files.
	w.manifest.ClosureState = closureState
	w.manifest.ClosedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if err := w.writeManifest(); err != nil {
		return err
	}

	if err := writeText(filepath.Join(w.runDir, "artifacts", "stdout.digest"), w.stdoutDigest()+"\n"); err != nil {
		return err
	}

	return writeText(filepath.Join(w.runDir, "artifacts", "stderr.digest"), w.stderrDigest()+"\n")
}

func (w *runWriter) stdoutDigest() string {
	// Digest accessors keep stream hashing behind the writer abstraction used by
	// command capture and event emission.
	return w.stdoutHash.Digest()
}

func (w *runWriter) stderrDigest() string {
	// Stderr is tracked separately because verifier output and command output
	// can have different trust interpretations.
	return w.stderrHash.Digest()
}

func (w *runWriter) lastEventHash() string {
	// The latest event hash is the closure anchor written into run-closed
	// payloads and manifests.
	return w.lastHash
}

func (w *runWriter) eventCount() int {
	// The event slice is the append-time view used for closure accounting.
	return len(w.events)
}
