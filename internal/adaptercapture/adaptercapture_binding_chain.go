package adaptercapture

func sameChainBindingCondition(run RunEvidence, event AdapterEvent) Condition {
	// sameChainBindingCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if eventAfterRunClosure(run, event.Sequence) {

		return cannotVerify("run_binding_established", "late_adapter_event", "adapter event appears after run closure", "Do not use late adapter events to satisfy capture-depth assessment.")
	}
	if sameChainDigestMissing(event) {

		return cannotVerify("run_binding_established", "same_chain_digest_missing", "same-chain adapter event lacks hash linkage", "Record prev_event_hash and event_hash.")
	}
	return Condition{}
}

func eventAfterRunClosure(run RunEvidence, sequence int) bool {
	return run.RunClosedSequence > 0 && sequence > run.RunClosedSequence
}

func sameChainDigestMissing(event AdapterEvent) bool {
	return event.EventHash == "" || event.PrevEventHash == ""
}

func adapterBundleBindingCondition(run RunEvidence, event AdapterEvent) Condition {
	// adapterBundleBindingCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if adapterBundleUnbound(run.AdapterBundle, event) {

		return cannotVerify("run_binding_established", "adapter_bundle_unbound", "adapter bundle is not bound to the selected run", "Bind the adapter bundle head digest into the run artifact.")
	}
	if run.RunClosedSequence > 0 && run.AdapterBundle.ReferencedSequence > run.RunClosedSequence {

		return cannotVerify("run_binding_established", "late_adapter_bundle", "adapter bundle was first referenced after run closure", "Reference adapter bundles before run closure.")
	}
	return Condition{}
}

func adapterBundleUnbound(bundle *AdapterBundle, event AdapterEvent) bool {
	// adapterBundleUnbound preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	return bundle == nil ||
		event.AdapterBundleHeadDigest == "" ||
		event.AdapterBundleHeadDigest != bundle.HeadDigest ||
		event.AdapterBundleID != bundle.BundleID
}
