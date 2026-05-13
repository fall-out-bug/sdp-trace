package witness

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func hashRunArtifacts(runsRoot string) ([]ArtifactDigest, error) {
	// Run discovery is delegated to demo rules so witness binding and gate/demo
	// evaluation agree on what counts as a run directory.
	runDirs, err := demo.DiscoverRunDirs(runsRoot)
	if err != nil {
		return nil, err
	}
	artifacts := make([]ArtifactDigest, 0, len(runDirs))
	for _, runDir := range runDirs {
		digest, err := hashFile(filepath.Join(runDir, "run.json"))
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, ArtifactDigest{
			Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
			SHA256: digest,
		})
	}
	return artifacts, nil
}

func hashFile(path string) (string, error) {
	// File hashing returns only SHA-256 hex; callers decide how the digest is
	// bound to a run or report artifact.
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("hash %s: %w", path, err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
