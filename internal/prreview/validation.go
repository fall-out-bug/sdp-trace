package prreview

func Validate(packet Packet, profile ReviewProfile, runs RunSet, ledger Ledger) Validation {
	// Validate reports review coverage and authority boundaries; it never turns
	// reviewer evidence into merge, release, or risk-acceptance approval.
	required := requiredPlaneSet(profile.RequiredPlanes)
	reasons := []string{}
	nextActions := []string{}
	cannotVerify := appendDigestValidation(packet, runs, ledger, &reasons, &nextActions)
	planeResults, usableCount, planesCannotVerify := validateRequiredPlanes(required, reviewRolesByID(profile.Roles), runs, &reasons, &nextActions)
	safeFindings, unresolved, findingsCannotVerify := validateLedgerFindings(packet, ledger, &reasons)
	state := reviewCoverageState(required, usableCount, cannotVerify || planesCannotVerify || findingsCannotVerify, unresolved)

	return validationResult(packet, planeResults, safeFindings, state, reasons, nextActions)
}

func validationResult(packet Packet, planeResults []PlaneResult, findings []LedgerFinding, state string, reasons, nextActions []string) Validation {
	// These authority fields are intentionally fixed outputs of validation:
	// review coverage can inform a human decision but cannot authorize it.
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

// reviewRolesByID lets plane validation compare observed reviewer results with
// the declared role model without changing the profile artifact.
func reviewRolesByID(roles []ReviewRole) map[string]ReviewRole {
	roleByID := map[string]ReviewRole{}
	for _, role := range roles {
		roleByID[role.RoleID] = role
	}
	return roleByID
}
