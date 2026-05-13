package harnessobs

func compose(dimensions []Dimension) (string, string) {
	// compose keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	state := StatePass
	reason := "all_required_dimensions_observed"
	for _, dim := range dimensions {
		if !dim.Required {

			continue
		}

		if rank(dim.State) > rank(state) {
			state = dim.State
			reason = dim.ReasonCode
		}
	}
	return state, reason
}
