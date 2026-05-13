package posture

import (
	"time"
)

func buildExportResult(selectionPath string, now time.Time, input buildInput, inputs []InputSelection, metricRows []MetricRow, movementRows []MovementRow, summary MovementSummary, refusals []RefusalRow) ExportResult {
	// buildExportResult keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	result := exportResultHeader(selectionPath, now, input, metricRows, refusals)
	result.InputSelection = inputs
	result.MetricRows = metricRows
	result.MovementRows = movementRows
	result.MovementSummary = summary
	result.RefusalRows = refusals
	return result
}

func exportResultHeader(selectionPath string, now time.Time, input buildInput, metricRows []MetricRow, refusals []RefusalRow) ExportResult {
	// exportResultHeader keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return ExportResult{

		SchemaVersion:        SchemaVersion,
		ExportProfileID:      ProfileID,
		ExportProfileVersion: ProfileVer,
		ExportID:             deterministicExportID(selectionPath, metricRows, refusals),
		Producer:             "sdp-trace",
		GeneratedAt:          now.UTC().Format(time.RFC3339),

		GroupingSetID:      input.selection.GroupingSetID,
		ActiveGroupingKeys: input.activeKeys,
		Handoff:            input.handoff,
		OutputSafety:       exportOutputSafety(),
	}
}

func exportOutputSafety() OutputSafety {

	return OutputSafety{VerifiedAbsentSensitiveClasses: SensitiveClasses()}
}
