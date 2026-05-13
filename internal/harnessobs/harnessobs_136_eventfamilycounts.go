package harnessobs

func eventFamilyCounts(events []Event) map[string]int {
	// eventFamilyCounts keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	counts := map[string]int{}
	for _, event := range events {

		counts[event.EventFamily]++
	}
	return counts
}
