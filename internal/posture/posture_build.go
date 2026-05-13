package posture

import (
	"time"
)

func Build(selectionPath string, now time.Time) (ExportResult, error) {
	// Build keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	selection, err := readSelection(selectionPath)
	if err != nil {
		return ExportResult{}, err
	}
	input, err := prepareBuildSelectionInput(now, selection)
	if err != nil {
		return ExportResult{}, err
	}
	inputs, refusals, groups := ingestRepositories(input.selection, input.activeKeys, input.cutoff, input.hasCutoff)
	metricRows := buildMetrics(groups)
	movementRows, summary := buildMovements(metricRows, input.selection.CurrentWindow, input.selection.PreviousWindow)

	return buildExportResult(selectionPath, now, input, inputs, metricRows, movementRows, summary, refusals), nil
}
