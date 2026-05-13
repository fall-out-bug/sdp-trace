package witness

func newCustomerPKIRecord(runsRoot, reportDir string) (Record, error) {
	record := baseRecord(KindCustomerPKI)
	// Artifact digests are captured before external policy validation. The
	// external signer authorizes this concrete payload, not an abstract run.
	// Report artifacts are optional supporting material and do not replace run
	// artifact binding.
	// RequestedTrustScope is set to external before validation so failures still
	// show which trust upgrade was attempted.
	// ProviderKind stays customer-pki because this profile is not tied to a CI
	// vendor.
	runArtifacts, err := hashRunArtifacts(runsRoot)
	if err != nil {
		return Record{}, err
	}
	record.RunArtifacts = runArtifacts
	if reportDir != "" {
		reportArtifacts, err := hashReportArtifacts(reportDir)
		if err != nil {
			return Record{}, err
		}
		record.ReportArtifacts = reportArtifacts
	}
	record.ProfileID = "customer-pki-v1"
	record.ProfileVersion = "1.0"
	record.ProviderKind = KindCustomerPKI
	record.RequestedTrustScope = TrustScopeExternal
	return record, nil
}
