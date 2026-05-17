package posture

import "fmt"

func prepareBuildSelectionInput(selection SelectionManifest) (buildInput, error) {
	// prepareBuildSelectionInput keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	activeKeys, err := validateBuildSelection(selection)
	if err != nil {
		return buildInput{}, err
	}
	cutoff, hasCutoff, err := parseFreshnessBoundary(selection.FreshnessBoundary)
	if err != nil {
		return buildInput{}, err
	}

	handoff, err := validatedHandoff(selection.Handoff)
	if err != nil {
		return buildInput{}, err
	}
	return buildInput{selection: selection, activeKeys: activeKeys, cutoff: cutoff, hasCutoff: hasCutoff, handoff: handoff}, nil
}

func validateBuildSelection(selection SelectionManifest) ([]string, error) {
	// validateBuildSelection keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if err := validateSelection(selection); err != nil {
		return nil, err
	}

	activeKeys := groupingKeys(selection.GroupingSetID)
	if len(activeKeys) == 0 {
		return nil, fmt.Errorf("unsupported grouping set")
	}
	return activeKeys, nil
}

func validatedHandoff(handoff map[string]string) (map[string]string, error) {
	// validatedHandoff keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if handoff == nil {

		handoff = map[string]string{}
	}
	if !safeHandoff(handoff) {
		return nil, fmt.Errorf("unsafe handoff")
	}
	return handoff, nil
}
