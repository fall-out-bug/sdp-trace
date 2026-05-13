package harnessobs

func defaultSessionStreamCapture(profile *SessionProfile) {
	// defaultSessionStreamCapture keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if profile.StreamCapture == "" {

		profile.StreamCapture = "disabled"
	}
}
