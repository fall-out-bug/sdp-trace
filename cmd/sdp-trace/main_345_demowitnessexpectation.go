package main

import (
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
