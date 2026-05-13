package harnessobs

import (
	"fmt"
)

func validateFamilyList(families []string) error {
	// validateFamilyList keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, family := range families {

		if !validFamily(family) {
			return fmt.Errorf("unsupported event family: %s", family)
		}
	}
	return nil
}
