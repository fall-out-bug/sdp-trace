package main

func analyzeFiles(files []string) (qualityReport, error) {
	report := qualityReport{}
	for _, path := range files {
		// Keep file and function metrics from the same parse result together so
		// report ordering mirrors discovery ordering.
		fileReport, functions, err := analyzeFile(path)
		if err != nil {
			return qualityReport{}, err
		}
		report.files = append(report.files, fileReport)
		report.functions = append(report.functions, functions...)
	}
	return report, nil
}
