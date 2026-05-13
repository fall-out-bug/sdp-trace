package query

import "path/filepath"

func loadForensicInput(runDir string, inputs packInputs) (packInputs, error) {
	// loadForensicInput keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	var forensic assessmentEnvelope
	artifact, present, err := readOptionalPackArtifact(filepath.Join(runDir, "forensic-retention.assessment-result.json"), "forensic_retention", "forensic_retention", false, &forensic)
	if err != nil && artifact.Role == "" {
		return packInputs{}, err
	}
	if present {
		inputs.forensicPresent = true
		inputs.forensicArtifact = &artifact
		inputs.forensic = forensic
		inputs.forensicErr = err
	}
	return inputs, nil
}

func loadAdapterInput(runDir string, inputs packInputs) (packInputs, error) {
	// loadAdapterInput keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	var adapter assessmentEnvelope
	artifact, present, err := readOptionalPackArtifact(filepath.Join(runDir, "adapter-capture.assessment-result.json"), "adapter_capture", "adapter_capture", false, &adapter)
	if err != nil && artifact.Role == "" {
		return packInputs{}, err
	}
	if present {
		inputs.adapterPresent = true
		inputs.adapterArtifact = &artifact
		inputs.adapter = adapter
		inputs.adapterErr = err
	}
	return inputs, nil
}
