package demo

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func evaluateRequiredRuns(rows []RunRow, contract trace.Contract) []RequiredRunResult {
	// evaluateRequiredRuns keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	results := make([]RequiredRunResult, 0, len(contract.RequiredRuns))
	rowsByWrapper := firstRowByWrapper(rows)
	for _, required := range contract.RequiredRuns {

		profile := required.Profile
		if profile == "" {
			profile = GateModeObservation
		}
		result := requiredRunResultTemplate(required, profile)
		if row, ok := rowsByWrapper[required.WrapperName]; ok {

			result = matchRequiredRun(row, required, result)
		}
		result = applyProtectedFutureConstraint(result, required.ID)
		results = append(results, result)
	}
	return results
}

func requiredRunResultTemplate(required trace.RequiredRun, profile string) RequiredRunResult {
	// requiredRunResultTemplate keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return RequiredRunResult{
		ID:          required.ID,
		WrapperName: required.WrapperName,
		Profile:     profile,
		State:       GateMissingTelemetry,
		Reasons: []string{
			fmt.Sprintf("required run %s with wrapper %s is missing", required.ID, required.WrapperName),
		},
	}
}
