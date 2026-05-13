package main

func writeFunctionMIBaselineFile(path string, report qualityReport, threshold float64) error {
	// Function baselines persist only current below-threshold function rows.
	return writeJSONFile(path, buildFunctionMIBaseline(report, threshold))
}

func writeFileMIBaselineFile(path string, report qualityReport, threshold float64) error {
	// File baselines persist only current below-threshold file rows.
	return writeJSONFile(path, buildFileMIBaseline(report, threshold))
}
