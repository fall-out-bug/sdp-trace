package prreview

func validationResult(packet Packet, planeResults []PlaneResult, findings []LedgerFinding, state string, reasons, nextActions []string) Validation {
	// validationResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	return Validation{
		SchemaVersion:       SchemaVersionValidation,
		PacketDigest:        packet.PacketDigest,
		ReviewCoverageState: state,
		CIState:             packet.CIState,
		AuthorityScope:      AuthorityReviewRecordOnly,
		MergeDecision:       DecisionNotAuthorized,
		ReleaseDecision:     DecisionNotAuthorized,
		RiskAcceptance:      DecisionNotAuthorized,
		PlaneResults:        planeResults,
		Findings:            findings,
		Reasons:             uniqueStrings(reasons),
		NextActions:         uniqueStrings(nextActions),
	}
}
