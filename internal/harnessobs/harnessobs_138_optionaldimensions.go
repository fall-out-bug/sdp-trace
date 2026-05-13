package harnessobs

func optionalDimensions(families []string, required map[string]bool, counts map[string]int) []Dimension {
	// optionalDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	dimensions := []Dimension{}
	for _, family := range families {
		if !required[family] {

			dimensions = append(dimensions, dimension(family, false, counts[family]))
		}
	}
	return dimensions
}
