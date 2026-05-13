package adaptercapture

func runBindingCondition(run RunEvidence) Condition {
	// Run binding ties adapter evidence back to the current run chain instead of
	// trusting checked-in adapter data by itself.
	if runIdentityMissing(run) {

		return cannotVerify("run_binding_established", "run_binding_missing", "run id or nonce is missing", "Record run id and run nonce before assessing adapter capture.")
	}
	for _, event := range run.AdapterEvents {
		if condition := adapterEventRunBindingCondition(run, event); condition.ID != "" {
			return condition
		}
	}
	return pass("run_binding_established", "run_binding_established", "adapter events are bound to run id, nonce, and chain or bundle context")
}

func runIdentityMissing(run RunEvidence) bool {
	return run.RunID == "" || run.RunNonce == ""
}

func adapterEventRunBindingCondition(run RunEvidence, event AdapterEvent) Condition {
	// adapterEventRunBindingCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if adapterEventRunIdentityMismatch(run, event) {
		return fail("run_binding_established", "run_binding_mismatch", "adapter event contradicts run id or nonce", "Use adapter events bound to the selected run.")
	}

	return adapterEventBindingModeCondition(run, event)
}

func adapterEventBindingModeCondition(run RunEvidence, event AdapterEvent) Condition {
	// adapterEventBindingModeCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	switch event.BindingMode {
	case BindingSameChain:
		return sameChainBindingCondition(run, event)
	case BindingAdapterBundle:
		return adapterBundleBindingCondition(run, event)
	default:
		return cannotVerify("run_binding_established", "binding_mode_missing", "adapter event binding mode is missing or unsupported", "Use same_chain or adapter_bundle binding.")
	}
}

func adapterEventRunIdentityMismatch(run RunEvidence, event AdapterEvent) bool {
	return event.RunID != run.RunID || event.RunNonce != run.RunNonce
}
