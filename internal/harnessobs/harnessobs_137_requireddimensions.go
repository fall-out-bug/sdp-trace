package harnessobs

func requiredDimensions(families []string, counts map[string]int) ([]Dimension, map[string]bool) {
	// requiredDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	required := map[string]bool{}
	dimensions := make([]Dimension, 0, len(families))
	for _, family := range families {

		required[family] = true
		dimensions = append(dimensions, dimension(family, true, counts[family]))
	}
	return dimensions, required
}
