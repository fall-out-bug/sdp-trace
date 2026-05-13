package witness

func populateCIEnvelopeArtifacts(record *Record, runsRoot, reportDir string) error {
	// Artifact population defines the comparison set for later envelope
	// validation. A missing or unreadable run artifact blocks the profile
	// instead of silently shrinking the evidence surface.
	runArtifacts, err := hashRunArtifacts(runsRoot)
	if err != nil {
		return err
	}
	record.RunArtifacts = runArtifacts
	if reportDir == "" {
		return nil
	}
	reportArtifacts, err := hashReportArtifacts(reportDir)
	if err != nil {
		return err
	}
	record.ReportArtifacts = reportArtifacts
	return nil
}
