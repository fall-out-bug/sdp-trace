package adaptercapture

func setValidEventBinding(event *AdapterEvent, seed validEventSeed) {
	// setValidEventBinding preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	event.RunID = seed.runID
	event.RunNonce = seed.nonce
	event.SourceBaseline = seed.source
	event.BindingMode = BindingSameChain
	event.Sequence = seed.sequence
	event.PrevEventHash = "1111111111111111111111111111111111111111111111111111111111111111"
	event.EventHash = "2222222222222222222222222222222222222222222222222222222222222222"
}
