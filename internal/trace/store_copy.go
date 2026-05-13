package trace

import (
	"io"
	"os"
	"path/filepath"
)

// CopyArtifactFile copies a verifier/export artifact into a run directory.
func CopyArtifactFile(src, dst string) error {
	// CopyArtifactFile preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()
	return copyArtifactReader(input, dst)
}
func copyArtifactReader(input io.Reader, dst string) error {
	// copyArtifactReader preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		// Parent creation is a local filesystem effect, not proof authority.
		return err
	}
	output, err := os.Create(dst)
	if err != nil {
		// Refuse to claim artifact retention when the destination cannot open.
		return err
	}
	defer output.Close()
	if _, err := io.Copy(output, input); err != nil {
		// Partial copies are surfaced so callers do not reference weak artifacts.
		return err
	}
	// Sync before returning so later proof material can safely reference dst.
	return output.Sync()
}
