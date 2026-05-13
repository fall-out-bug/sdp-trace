package query

import "path/filepath"

func loadPackInputs(runDir string) (packInputs, error) {
	// loadPackInputs keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	var run runArtifact
	runArtifact, err := readPackArtifact(filepath.Join(runDir, "run.json"), "run", "run", true, &run)
	if err != nil && runArtifact.Role == "" {
		return packInputs{}, err
	}
	inputs := packInputs{run: run, runArtifact: runArtifact, runErr: err}
	return loadOptionalPackInputs(runDir, inputs)
}

func loadOptionalPackInputs(runDir string, inputs packInputs) (packInputs, error) {
	// loadOptionalPackInputs keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	var err error
	inputs, err = loadForensicInput(runDir, inputs)
	if err != nil {
		return packInputs{}, err
	}
	return loadAdapterInput(runDir, inputs)
}
