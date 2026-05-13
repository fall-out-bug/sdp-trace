package witness

func hydrateGitHubArtifacts(record *Record, runsRoot, reportDir string) error {
	// Artifact hashes are local evidence that the OIDC-backed identity must bind
	// to; they do not by themselves establish CI witness trust.
	runArtifacts, err := hashRunArtifacts(runsRoot)
	if err != nil {
		return err
	}
	record.RunArtifacts = runArtifacts
	if reportDir != "" {
		// Report artifacts are optional because some witness callers only bind
		// the CI run output directory.
		reportArtifacts, err := hashReportArtifacts(reportDir)
		if err != nil {
			return err
		}
		record.ReportArtifacts = reportArtifacts
	}
	return nil
}
