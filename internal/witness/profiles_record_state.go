package witness

func applyProfileState(record *Record, status, scope, reason string) {
	// Keep human status, trust scope, and machine-readable reason code aligned
	// whenever a profile decision is applied.
	record.Status = status
	record.TrustScope = scope
	record.EstablishedTrustScope = scope
	record.Reason = reason
	record.ReasonCodes = []string{reason}
}
