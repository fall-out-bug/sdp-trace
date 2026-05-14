package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func demoWitnessArtifact(runDir string) (string, string, error) {
	// The retained run artifact supplies the run id while the file bytes supply
	// the digest that witness summaries must match.
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		return "", "", err
	}
	digest, err := sha256File(runDir, "run.json")
	if err != nil {
		return "", "", err
	}
	return artifact.Manifest.RunID, digest, nil
}
