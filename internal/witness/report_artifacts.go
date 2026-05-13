package witness

import (
	"os"
	"path/filepath"
	"sort"
)

func hashReportArtifacts(reportDir string) ([]ArtifactDigest, error) {
	// Report artifacts are read from one directory level and then sorted by
	// retained path for deterministic witness output.
	entries, err := os.ReadDir(reportDir)
	if err != nil {
		return nil, err
	}
	return hashReportArtifactEntries(reportDir, entries)
}

func hashReportArtifactEntries(reportDir string, entries []os.DirEntry) ([]ArtifactDigest, error) {
	// The generated ci-witness.json is excluded so rerunning witness generation
	// does not self-reference the prior output.
	artifacts := make([]ArtifactDigest, 0)
	for _, entry := range entries {
		if skipReportArtifactEntry(entry) {
			continue
		}
		artifact, err := hashReportArtifact(reportDir, entry.Name())
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	sort.Slice(artifacts, func(i, j int) bool {
		return artifacts[i].Path < artifacts[j].Path
	})
	return artifacts, nil
}

func skipReportArtifactEntry(entry os.DirEntry) bool {
	// Witness output excludes directories and the prior witness file to avoid
	// self-referential report evidence.
	return entry.IsDir() || entry.Name() == "ci-witness.json"
}

func hashReportArtifact(reportDir, name string) (ArtifactDigest, error) {
	// Paths stored in witness output are report-relative and slash-normalized for
	// cross-platform comparison.
	digest, err := hashFile(filepath.Join(reportDir, name))
	if err != nil {
		return ArtifactDigest{}, err
	}
	return ArtifactDigest{Path: filepath.ToSlash(name), SHA256: digest}, nil
}
