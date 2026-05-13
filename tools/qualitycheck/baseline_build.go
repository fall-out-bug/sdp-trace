package main

func buildFunctionMIBaseline(report qualityReport, threshold float64) functionMIBaseline {
	// Function baselines are ratchets for current debt, not complete metric
	// snapshots.
	// Passing functions are intentionally absent so a future regression cannot
	// hide inside a full-metric snapshot.
	records := make([]functionMIBaselineRecord, 0, len(report.functions))
	for _, fn := range report.functions {
		// MI baselines only ratchet existing debt; passing functions are omitted
		// so future below-threshold appearances fail as missing baseline entries.
		// The comparison uses raw MI; rounding is only for persisted review rows.
		if fn.maintainabilityIndex >= threshold {
			continue
		}
		// Persist the stable function key plus display fields for human review.
		records = append(records, functionMIBaselineRecord{
			Key:                  functionKey(fn),
			File:                 fn.file,
			Line:                 fn.line,
			Name:                 fn.name,
			MaintainabilityIndex: roundMetric(fn.maintainabilityIndex),
		})
	}
	return functionMIBaseline{
		// The schema string lets CI reject stale ratchet files before comparing
		// metric rows.
		SchemaVersion: functionMIBaselineSchema,
		// Store the threshold used to generate the ratchet beside the records.
		// Reviewers can then see whether a baseline update changed policy or
		// only refreshed measured debt.
		Threshold: threshold,
		Functions: records,
	}
}
