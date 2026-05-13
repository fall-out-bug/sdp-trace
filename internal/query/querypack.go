package query

import "strings"

func ForensicsBasicPack(runDir string) (QueryPackResult, error) {
	return forensicsBasicPack(loadPackInputs(runDir))
}

func forensicsBasicPack(inputs packInputs, err error) (QueryPackResult, error) {
	// forensicsBasicPack keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if err != nil {
		return QueryPackResult{}, err
	}
	return buildForensicsBasicPack(inputs), nil
}

func buildForensicsBasicPack(inputs packInputs) QueryPackResult {
	// buildForensicsBasicPack keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	builder := newPackBuilder(inputs)
	if inputs.runErr != nil {
		builder.addMalformedRequiredInputRows()
		return builder.result()
	}
	builder.addTimelineRows()
	builder.addRedactionRows()
	builder.addCaptureRows()
	builder.addGapRows()
	builder.addUnverifiedClaimRows()
	builder.addSummaryRows()
	return builder.result()
}

func ExplainForensicsPack(result QueryPackResult) string {
	// ExplainForensicsPack keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	var lines []string
	for _, queryName := range queryOrder {
		rows := sortedQueryRows(result.QueryRows[queryName])
		for _, row := range rows {
			lines = append(lines, explainQueryRow(queryName, row))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}
