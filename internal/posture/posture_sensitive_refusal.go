package posture

func validRefusalReason(value string) bool {
	// validRefusalReason keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	switch value {
	case "stale_input", "malformed_input", "untrusted_input_digest_mismatch", "unsafe_label",
		"unsupported_input", "missing_required_input", "missing_optional_input",
		"non_comparable_metric_version", "non_comparable_dimension_key",
		"non_comparable_denominator_basis", "non_comparable_input_trust_rule",
		"non_comparable_missing_window", "output_safety_violation":

		return true
	default:
		return false
	}
}
