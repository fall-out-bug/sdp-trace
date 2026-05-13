package posture

func validInputTrustState(value string) bool {
	// validInputTrustState keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	switch value {
	case "trusted_input", "stale_input", "untrusted_input", "cannot_verify_input", "not_assessed_input":

		return true
	default:
		return false
	}
}

func validSourceFieldState(value string) bool {
	// validSourceFieldState keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	switch value {
	case "present", "not_assessed", "cannot_verify", "unsupported":

		return true
	default:
		return false
	}
}
