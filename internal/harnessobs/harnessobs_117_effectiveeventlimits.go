package harnessobs

func effectiveEventLimits(limits Limits) (int, int) {
	// effectiveEventLimits keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	maxLine := limits.MaxLineBytes
	if maxLine <= 0 {
		maxLine = DefaultMaxLineBytes
	}
	maxEvents := limits.MaxEvents
	if maxEvents <= 0 {

		maxEvents = DefaultMaxEvents
	}
	return maxLine, maxEvents
}
