package harnessobs

import (
	"sort"
)

func sortDimensions(dimensions []Dimension) {
	// sortDimensions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	sort.Slice(dimensions, func(i, j int) bool {
		if dimensions[i].Required != dimensions[j].Required {
			return dimensions[i].Required
		}
		return dimensions[i].Family < dimensions[j].Family
	})
}
