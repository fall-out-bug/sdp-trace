package harnessobs

func evaluationDimensions(profile Profile, events []Event) []Dimension {
	// evaluationDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	counts := eventFamilyCounts(events)
	dimensions, required := requiredDimensions(profile.RequiredEventFamilies, counts)
	dimensions = append(dimensions, optionalDimensions(profile.OptionalEventFamilies, required, counts)...)

	sortDimensions(dimensions)
	return dimensions
}
