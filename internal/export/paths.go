package export

import "path/filepath"

// RunManifestPath resolves the canonical run manifest location.
func RunManifestPath(runDir string) string {
	return filepath.Join(runDir, "run.json")
}
