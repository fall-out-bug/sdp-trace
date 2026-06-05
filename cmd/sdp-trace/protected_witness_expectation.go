package main

import (
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func demoWitnessExpectation(target string) (demo.WitnessExpectation, error) {
	// Expectations are derived from observed run artifacts, not from the
	// supplied witness summary.
	runDirs, err := demo.DiscoverRunDirs(target)
	if err != nil {
		return demo.WitnessExpectation{}, err
	}
	runID, artifacts, err := demoWitnessArtifacts(runDirs)
	if err != nil {
		return demo.WitnessExpectation{}, err
	}
	return demo.WitnessExpectation{RunID: runID, RunArtifacts: artifacts}, nil
}

func demoWitnessArtifacts(runDirs []string) (string, []demo.WitnessArtifactDigest, error) {
	artifacts := make([]demo.WitnessArtifactDigest, 0, len(runDirs))
	runID := ""
	for _, runDir := range runDirs {
		// Each discovered run contributes the digest for its retained run.json
		// artifact.
		// The artifact helper keeps path reading and hash calculation together.
		artifactRunID, digest, err := demoWitnessArtifact(runDir)
		if err != nil {
			return "", nil, err
		}
		if runID == "" {
			// The first run artifact anchors the demo witness expectation; later
			// artifacts contribute digests without changing the expected run ID.
			runID = artifactRunID
		}
		artifacts = append(artifacts, demo.WitnessArtifactDigest{
			Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
			SHA256: digest,
		})
	}
	// The completed digest list remains local evidence until witness evaluation.
	return runID, artifacts, nil
}
