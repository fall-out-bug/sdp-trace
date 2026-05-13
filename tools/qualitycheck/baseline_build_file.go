package main

func buildFileMIBaseline(report qualityReport, threshold float64) fileMIBaseline {
	// File baselines mirror the function ratchet: only existing below-threshold
	// files are allowed to remain.
	// New below-threshold files are therefore rejected unless the baseline is
	// deliberately regenerated and reviewed.
	records := make([]fileMIBaselineRecord, 0, len(report.files))
	for _, file := range report.files {
		// Store only below-threshold files so the baseline remains a debt ledger,
		// not a snapshot of every measured source file.
		if file.maintainabilityIndex >= threshold {
			continue
		}
		records = append(records, fileMIBaselineRecord{
			Key:                  fileKey(file),
			File:                 file.file,
			MaintainabilityIndex: roundMetric(file.maintainabilityIndex),
		})
	}
	return fileMIBaseline{
		// Keep threshold alongside rows so future invocations can detect a
		// baseline written for a different gate.
		// The schema identifies file-level keys separately from function keys.
		SchemaVersion: fileMIBaselineSchema,
		Threshold:     threshold,
		Files:         records,
	}
}
