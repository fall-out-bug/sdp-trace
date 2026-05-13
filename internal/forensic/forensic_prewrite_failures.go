package forensic

func prewriteEventFailures(event EventRetention) []conditionFailure {
	// prewriteEventFailures keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	return []conditionFailure{
		{

			matched:   prewriteEventHasSecretLike(event),
			condition: fail("redaction_prewrite_applied", "secret_like_value_persisted", "secret-like value is marked as persisted in retained metadata", "Apply pre-write redaction and retain only digests or safe references."),
		},
		{

			matched:   prewriteMissingRedactionDigests(event),
			condition: cannotVerify("redaction_prewrite_applied", "redaction_digest_missing", "pre-write redaction digests are missing", "Record pre-redaction and redacted payload digests."),
		},
		{

			matched:   prewriteRuleRefsMissing(event),
			condition: cannotVerify("redaction_prewrite_applied", "redaction_rule_refs_missing", "redaction rule references are missing", "Record the redaction rule ids applied before persistence."),
		},
	}
}
