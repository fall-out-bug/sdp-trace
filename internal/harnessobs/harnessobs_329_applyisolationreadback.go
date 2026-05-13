package harnessobs

func applyIsolationReadback(result *SessionIsolationResult, ok bool) {
	// applyIsolationReadback keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !ok {

		result.State = StateCannotVerify
		result.ReasonCode = "isolation_rule_absent"
	}
}
