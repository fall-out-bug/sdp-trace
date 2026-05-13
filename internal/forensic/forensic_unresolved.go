package forensic

func unresolvedCondition(input Input) Condition {
	// unresolvedCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, event := range input.Run.Events {

		if condition, ok := unresolvedConditionForEvent(event); ok {
			return condition
		}
	}
	return pass("redaction_unresolved_visible", "redaction_resolved", "redaction states are resolved or explicitly non-blocking")
}

func unresolvedConditionForEvent(event EventRetention) (Condition, bool) {
	// unresolvedConditionForEvent keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if event.RedactionUnresolved {

		return fail("redaction_unresolved_visible", "redaction_unresolved", "unresolved redaction is visible and blocks forensic retention", "Resolve redaction or lower the forensic claim."), true
	}
	if event.RedactionAction != RedactionActionWithhold {

		return Condition{}, false
	}
	return withholdingCondition(event.Withholding)
}

func withholdingCondition(withholding *Withholding) (Condition, bool) {
	// withholdingCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if withholdingAuditMissing(withholding) {
		return cannotVerify("redaction_unresolved_visible", "withholding_audit_missing", "withholding lacks required audit evidence", "Record withholding authority, requestor when different, reason, and justification."), true
	}
	if withholding.Authority.VerificationState != AuthorityVerified {

		return cannotVerify("redaction_unresolved_visible", "withholding_authority_unverifiable", "withholding authority is not provenance or accountability verified", "Record verified withholding authority evidence."), true
	}
	return Condition{}, false
}

func withholdingAuditMissing(withholding *Withholding) bool {
	// withholdingAuditMissing keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return withholding == nil ||
		withholding.Authority.ActorID == "" ||
		withholding.ReasonCode == "" ||
		withholding.Justification == ""
}
