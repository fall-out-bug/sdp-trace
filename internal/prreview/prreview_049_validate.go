package prreview

func Validate(packet Packet, profile ReviewProfile, runs RunSet, ledger Ledger) Validation {
	// Validate keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	required := requiredPlaneSet(profile.RequiredPlanes)
	reasons := []string{}
	nextActions := []string{}
	cannotVerify := appendDigestValidation(packet, runs, ledger, &reasons, &nextActions)
	planeResults, usableCount, planesCannotVerify := validateRequiredPlanes(required, reviewRolesByID(profile.Roles), runs, &reasons, &nextActions)
	safeFindings, unresolved, findingsCannotVerify := validateLedgerFindings(packet, ledger, &reasons)
	state := reviewCoverageState(required, usableCount, cannotVerify || planesCannotVerify || findingsCannotVerify, unresolved)

	return validationResult(packet, planeResults, safeFindings, state, reasons, nextActions)
}
